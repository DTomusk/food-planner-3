package middleware

import (
	"context"
	"fmt"
	"foodplanner/internal/auth"
	"foodplanner/internal/correlationid"
	"foodplanner/internal/events"
	"foodplanner/internal/logging"
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
	eventPublisher events.Publisher,
) func(next http.Handler) http.Handler {
	windowSeconds := int(cfg.Window.Seconds())
	if windowSeconds <= 0 {
		windowSeconds = 60 // Default to 1 minute if invalid duration is provided
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.FromContext(r.Context())

			if r.Method == http.MethodOptions {
				// Skip rate limiting for preflight requests
				next.ServeHTTP(w, r)
				return
			}

			limit, subject := getLimitAndSubject(r, cfg)

			// Generate Redis key based on subject and current time bucket
			nowBucket := getNowBucket(windowSeconds)
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
				// Track how many times requests were blocked by the rate limiter
				events.IncrementRateLimiterBlocked()

				// Publish an audit event only on first breach in this window (not on every blocked request)
				if eventPublisher != nil && isFirstBreachInWindow(r.Context(), client, subject, windowSeconds) {
					event := events.NewRateLimitExceededEvent(
						correlationid.FromContext(r.Context()),
						subject,
						getIPFromContext(r.Context()),
						getUserAgentFromRequestContext(r),
						r.Method,
						r.URL.Path,
						limit,
						count,
						windowSeconds,
						ttl,
					)

					err = eventPublisher.Publish(r.Context(), event)
					if err != nil {
						logger.Warn(
							"Failed to publish rate limit exceeded event",
							"subject", subject,
							"method", r.Method,
							"path", r.URL.Path,
							"err", err,
						)
					}
				}

				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getNowBucket(windowSeconds int) int64 {
	return time.Now().UTC().Unix() / int64(windowSeconds)
}

func getIPFromContext(ctx context.Context) string {
	ip, _ := ctx.Value(IPKey).(string)
	return ip
}

func getUserAgentFromRequestContext(r *http.Request) string {
	if r == nil {
		return ""
	}
	if userAgent, ok := r.Context().Value(UserAgentKey).(string); ok && userAgent != "" {
		return userAgent
	}
	return r.UserAgent()
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

// Check if first breach in window, used for auditing
func isFirstBreachInWindow(ctx context.Context, client redis.UniversalClient, subject string, windowSeconds int) bool {
	nowBucket := getNowBucket(windowSeconds)
	breachKey := fmt.Sprintf("rl:audit:first:%s:%d", subject, nowBucket)

	// sets key if it doesn't exist i.e. is first breach for subject and bucket
	result, err := client.SetNX(ctx, breachKey, "1", time.Duration(windowSeconds)*time.Second).Result()
	if err != nil {
		// return false on error to avoid flooding audit log
		return false
	}
	return result
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
