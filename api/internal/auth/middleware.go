package auth

import (
	"context"
	"net/http"
	"strings"
)

func Middleware(jwtService *JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				// If a bearer token is provided but malformed, mark it as invalid
				ctx = ContextWithInvalidAuthToken(ctx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			claims, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				// Simply don't add claims, allow graphql to deal with it
				// But if there was an attempt to provide a token and it was invalid, mark it as such so resolvers can return unauthenticated errors
				ctx = ContextWithInvalidAuthToken(ctx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx = context.WithValue(ctx, UserClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
