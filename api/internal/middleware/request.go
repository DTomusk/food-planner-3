package middleware

import (
	"context"
	"net/http"
)

const RequestKey contextKey = "request"

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestKey, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
