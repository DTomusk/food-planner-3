package middleware

import (
	"foodplanner/internal/correlationid"
	"net/http"

	"github.com/google/uuid"
)

// CorrelationIDMiddleware sets a correlation ID for each request in context.
// It reads from the X-Correlation-ID request header if present, otherwise generates a new UUID.
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()

		if headerVal := r.Header.Get("X-Correlation-ID"); headerVal != "" {
			if parsed, err := uuid.Parse(headerVal); err == nil {
				id = parsed
			}
		}

		ctx := correlationid.WithContext(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
