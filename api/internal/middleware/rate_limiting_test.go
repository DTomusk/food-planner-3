package middleware

import (
	"context"
	"foodplanner/internal/auth"
	"foodplanner/internal/events"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type recordingPublisher struct {
	events []events.Event
}

func (p *recordingPublisher) Publish(ctx context.Context, event events.Event) error {
	p.events = append(p.events, event)
	return nil
}

func TestRateLimitingMiddleware_AllowsRequestUnderLimit_Anonymous(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       2,
		AuthenticatedLimit:   5,
		FailOpenOnRedisError: false,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/query", nil)
	req = req.WithContext(context.WithValue(req.Context(), IPKey, "203.0.113.10"))
	rr := httptest.NewRecorder()

	middleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Equal(t, "2", rr.Header().Get("X-RateLimit-Limit"))
	require.Equal(t, "1", rr.Header().Get("X-RateLimit-Remaining"))
}

func TestRateLimitingMiddleware_BlocksRequestOverLimit_Anonymous(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   5,
		FailOpenOnRedisError: false,
	}, nil)

	req1 := httptest.NewRequest(http.MethodGet, "/query", nil)
	req1 = req1.WithContext(context.WithValue(req1.Context(), IPKey, "203.0.113.11"))
	rr1 := httptest.NewRecorder()
	middleware(next).ServeHTTP(rr1, req1)
	require.Equal(t, http.StatusNoContent, rr1.Code)

	req2 := httptest.NewRequest(http.MethodGet, "/query", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), IPKey, "203.0.113.11"))
	rr2 := httptest.NewRecorder()
	middleware(next).ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusTooManyRequests, rr2.Code)
	require.NotEmpty(t, rr2.Header().Get("Retry-After"))
	require.Equal(t, "1", rr2.Header().Get("X-RateLimit-Limit"))
	require.Equal(t, "0", rr2.Header().Get("X-RateLimit-Remaining"))
}

func TestRateLimitingMiddleware_UsesAuthenticatedBucketWhenClaimsPresent(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   2,
		FailOpenOnRedisError: false,
	}, nil)

	claims := &auth.Claims{UserID: "user-123"}

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		ctx := context.WithValue(req.Context(), IPKey, "203.0.113.12")
		ctx = auth.ContextWithClaims(ctx, claims)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		middleware(next).ServeHTTP(rr, req)
		require.Equal(t, http.StatusNoContent, rr.Code)
	}

	req3 := httptest.NewRequest(http.MethodGet, "/query", nil)
	ctx3 := context.WithValue(req3.Context(), IPKey, "203.0.113.12")
	ctx3 = auth.ContextWithClaims(ctx3, claims)
	req3 = req3.WithContext(ctx3)
	rr3 := httptest.NewRecorder()
	middleware(next).ServeHTTP(rr3, req3)
	require.Equal(t, http.StatusTooManyRequests, rr3.Code)
}

func TestRateLimitingMiddleware_DifferentAnonymousSubjectsAreIsolated(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   2,
		FailOpenOnRedisError: false,
	}, nil)

	reqA := httptest.NewRequest(http.MethodGet, "/query", nil)
	reqA = reqA.WithContext(context.WithValue(reqA.Context(), IPKey, "203.0.113.13"))
	rrA := httptest.NewRecorder()
	middleware(next).ServeHTTP(rrA, reqA)
	require.Equal(t, http.StatusNoContent, rrA.Code)

	reqB := httptest.NewRequest(http.MethodGet, "/query", nil)
	reqB = reqB.WithContext(context.WithValue(reqB.Context(), IPKey, "203.0.113.14"))
	rrB := httptest.NewRecorder()
	middleware(next).ServeHTTP(rrB, reqB)
	require.Equal(t, http.StatusNoContent, rrB.Code)

	reqA2 := httptest.NewRequest(http.MethodGet, "/query", nil)
	reqA2 = reqA2.WithContext(context.WithValue(reqA2.Context(), IPKey, "203.0.113.13"))
	rrA2 := httptest.NewRecorder()
	middleware(next).ServeHTTP(rrA2, reqA2)
	require.Equal(t, http.StatusTooManyRequests, rrA2.Code)
}

func TestRateLimitingMiddleware_AnonymousAndAuthenticatedBucketsAreIsolatedForSameIP(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   2,
		FailOpenOnRedisError: false,
	}, nil)

	// First anonymous request from this IP should consume the anon bucket.
	anonReq := httptest.NewRequest(http.MethodGet, "/query", nil)
	anonReq = anonReq.WithContext(context.WithValue(anonReq.Context(), IPKey, "203.0.113.19"))
	anonRR := httptest.NewRecorder()
	middleware(next).ServeHTTP(anonRR, anonReq)
	require.Equal(t, http.StatusNoContent, anonRR.Code)
	require.Equal(t, "0", anonRR.Header().Get("X-RateLimit-Remaining"))

	// Authenticated request from the same IP should use user bucket, not IP bucket.
	authReq := httptest.NewRequest(http.MethodGet, "/query", nil)
	authCtx := context.WithValue(authReq.Context(), IPKey, "203.0.113.19")
	authCtx = auth.ContextWithClaims(authCtx, &auth.Claims{UserID: "user-isolated"})
	authReq = authReq.WithContext(authCtx)
	authRR := httptest.NewRecorder()
	middleware(next).ServeHTTP(authRR, authReq)
	require.Equal(t, http.StatusNoContent, authRR.Code)
	require.Equal(t, "1", authRR.Header().Get("X-RateLimit-Remaining"))

	// Anonymous bucket should still be exhausted for this IP.
	anonReq2 := httptest.NewRequest(http.MethodGet, "/query", nil)
	anonReq2 = anonReq2.WithContext(context.WithValue(anonReq2.Context(), IPKey, "203.0.113.19"))
	anonRR2 := httptest.NewRecorder()
	middleware(next).ServeHTTP(anonRR2, anonReq2)
	require.Equal(t, http.StatusTooManyRequests, anonRR2.Code)
}

func TestRateLimitingMiddleware_OptionsRequestsBypassLimiter(t *testing.T) {
	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(nil, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   1,
		FailOpenOnRedisError: false,
	}, nil)

	req := httptest.NewRequest(http.MethodOptions, "/query", nil)
	rr := httptest.NewRecorder()

	middleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestRateLimitingMiddleware_FailOpenOnRedisError_AllowsRequest(t *testing.T) {
	s := miniredis.RunT(t)
	addr := s.Addr()
	s.Close()

	client := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() { _ = client.Close() })

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   1,
		FailOpenOnRedisError: true,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/query", nil)
	ctx, cancel := context.WithTimeout(req.Context(), 200*time.Millisecond)
	defer cancel()
	ctx = context.WithValue(ctx, IPKey, "203.0.113.15")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	middleware(next).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestRateLimitingMiddleware_FailClosedOnRedisError_ReturnsServiceUnavailable(t *testing.T) {
	s := miniredis.RunT(t)
	addr := s.Addr()
	s.Close()

	client := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() { _ = client.Close() })

	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   1,
		FailOpenOnRedisError: false,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/query", nil)
	ctx, cancel := context.WithTimeout(req.Context(), 200*time.Millisecond)
	defer cancel()
	ctx = context.WithValue(ctx, IPKey, "203.0.113.16")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	middleware(next).ServeHTTP(rr, req)

	require.False(t, handlerCalled)
	require.Equal(t, http.StatusServiceUnavailable, rr.Code)
}

func TestRateLimitingMiddleware_UsesFallbackLimitsWhenConfiguredLimitsInvalid(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       0,
		AuthenticatedLimit:   0,
		FailOpenOnRedisError: false,
	}, nil)

	anonReq := httptest.NewRequest(http.MethodGet, "/query", nil)
	anonReq = anonReq.WithContext(context.WithValue(anonReq.Context(), IPKey, "203.0.113.17"))
	anonRR := httptest.NewRecorder()
	middleware(next).ServeHTTP(anonRR, anonReq)
	require.Equal(t, http.StatusNoContent, anonRR.Code)
	require.Equal(t, "60", anonRR.Header().Get("X-RateLimit-Limit"))

	authReq := httptest.NewRequest(http.MethodGet, "/query", nil)
	authCtx := context.WithValue(authReq.Context(), IPKey, "203.0.113.18")
	authCtx = auth.ContextWithClaims(authCtx, &auth.Claims{UserID: "user-fallback"})
	authReq = authReq.WithContext(authCtx)
	authRR := httptest.NewRecorder()
	middleware(next).ServeHTTP(authRR, authReq)
	require.Equal(t, http.StatusNoContent, authRR.Code)
	require.Equal(t, "180", authRR.Header().Get("X-RateLimit-Limit"))
}

func TestRateLimitingMiddleware_PublishesOnlyOnFirstBreachAndIncrementsMetricPerBreach(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	publisher := &recordingPublisher{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   2,
		FailOpenOnRedisError: false,
	}, publisher)

	before := events.SnapshotMetrics().RateLimiterBlockedTotal

	for i := 0; i < 4; i++ {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req = req.WithContext(context.WithValue(req.Context(), IPKey, "203.0.113.30"))
		rr := httptest.NewRecorder()
		middleware(next).ServeHTTP(rr, req)

		if i == 0 {
			require.Equal(t, http.StatusNoContent, rr.Code)
		} else {
			require.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}

	after := events.SnapshotMetrics().RateLimiterBlockedTotal
	require.Equal(t, uint64(3), after-before, "expected blocked metric to increment once per blocked request")
	require.Len(t, publisher.events, 1, "expected event publish only on first breach in window")

	rateLimitEvent, ok := publisher.events[0].(events.RateLimitExceededEvent)
	require.True(t, ok)
	require.Equal(t, "ip:203.0.113.30", rateLimitEvent.Subject)
	require.Equal(t, 1, rateLimitEvent.Limit)
	require.Equal(t, 2, rateLimitEvent.Count, "first breach should occur at count=2 when limit=1")
}

func TestRateLimitingMiddleware_PublishesFirstBreachEventPerSubject(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	publisher := &recordingPublisher{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	middleware := NewRateLimitingMiddleware(client, RateLimitingConfig{
		Window:               time.Minute,
		AnonymousLimit:       1,
		AuthenticatedLimit:   2,
		FailOpenOnRedisError: false,
	}, publisher)

	// Subject A: first request allowed, second request is first breach.
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req = req.WithContext(context.WithValue(req.Context(), IPKey, "203.0.113.40"))
		rr := httptest.NewRecorder()
		middleware(next).ServeHTTP(rr, req)
	}

	// Subject B: first request allowed, second request is first breach.
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req = req.WithContext(context.WithValue(req.Context(), IPKey, "203.0.113.41"))
		rr := httptest.NewRecorder()
		middleware(next).ServeHTTP(rr, req)
	}

	require.Len(t, publisher.events, 2, "expected one first-breach event per subject")

	eventA, ok := publisher.events[0].(events.RateLimitExceededEvent)
	require.True(t, ok)
	require.Equal(t, "ip:203.0.113.40", eventA.Subject)

	eventB, ok := publisher.events[1].(events.RateLimitExceededEvent)
	require.True(t, ok)
	require.Equal(t, "ip:203.0.113.41", eventB.Subject)
}
