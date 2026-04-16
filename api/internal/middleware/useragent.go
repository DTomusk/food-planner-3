package middleware

import (
	"context"
	"net/http"
	"strings"
)

const UserAgentKey contextKey = "user_agent"

func UserAgentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := getUserAgent(r)

		ctx := context.WithValue(r.Context(), UserAgentKey, userAgent)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserAgent(r *http.Request) string {
	return strings.TrimSpace(r.UserAgent())
}
