package middleware

import (
	"fmt"
	"foodplanner/internal/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitingConfig struct {
	// Time duration for which rate limiting is applied
	Window time.Duration
	// Number of allowed requests for unauthenticated users within time window
	AnonymousLimit int
	// Number of allowed requests for authenticated users within time window
	AuthenticatedLimit int
	// Whether to allow requests if Redis is unavailable (fail-open) or to block them (fail-closed)
	FailOpenOnRedisError bool
}

// Script for fixed window rate limiting. It increments the count for the current time bucket and sets an expiration if it's the first request
// Returns whether the request is allowed (count under limit), the current count, and the time until the bucket resets (TTL)
var fixedWindowScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("TTL", KEYS[1])
local limit = tonumber(ARGV[2])
local allowed = 1
if count > limit then
	allowed = 0
end
return {allowed, count, ttl}
`)

func NewRateLimitingMiddleware(
	client redis.UniversalClient,
	cfg RateLimitingConfig,
) func(next http.Handler) http.Handler {
	windowSeconds := int(cfg.Window.Seconds())
	if windowSeconds <= 0 {
		windowSeconds = 60 // Default to 1 minute if invalid duration is provided
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				// Skip rate limiting for preflight requests
				next.ServeHTTP(w, r)
				return
			}

			limit, subject := getLimitAndSubject(r, cfg)

			// Generate Redis key based on subject and current time bucket
			nowBucket := time.Now().UTC().Unix() / int64(windowSeconds)
			key := fmt.Sprintf("rate_limit:v1:%s:%d", subject, nowBucket)

			// Run script
			res, err := fixedWindowScript.Run(
				r.Context(),
				client,
				[]string{key},
				windowSeconds,
				limit,
			).Result()
			if err != nil {
				// If Redis error occurs, decide based on configuration whether to allow the request (fail-open) or block it (fail-closed)
				if cfg.FailOpenOnRedisError {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, "rate limiter unavailable", http.StatusServiceUnavailable)
				return
			}

			results, ok := res.([]interface{})
			if !ok || len(results) != 3 {
				// Unexpected result format from Redis script
				http.Error(w, "rate limiting error", http.StatusInternalServerError)
				return
			}

			allowed := toInt64(results[0]) == 1
			count := int(toInt64(results[1]))
			ttl := int(toInt64(results[2]))

			remaining := limit - count
			if remaining < 0 {
				remaining = 0
			}

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if ttl > 0 {
				w.Header().Set("Retry-After", strconv.Itoa(ttl))
			}

			if !allowed {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getLimitAndSubject(r *http.Request, cfg RateLimitingConfig) (int, string) {
	// If user logged in, use their ID as the subject
	if claims, ok := auth.ClaimsFromContext(r.Context()); ok && claims.UserID != "" {
		limit := cfg.AuthenticatedLimit
		if limit <= 0 {
			limit = 180 // Default authenticated limit if not set
		}
		return limit, "user:" + claims.UserID
	}
	// Otherwise, use IP address
	ip := ""
	if ipValue := r.Context().Value(IPKey); ipValue != nil {
		if ipStr, ok := ipValue.(string); ok {
			ip = ipStr
		}
	}
	if ip == "" {
		ip = "unknown"
	}
	limit := cfg.AnonymousLimit
	if limit <= 0 {
		limit = 60 // Default anonymous limit if not set
	}
	return limit, "ip:" + ip
}

func toInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	default:
		return 0
	}
}
