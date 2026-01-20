package auth

import "context"

type contextKey string

// UserContextKey is the context key used to store and retrieve the authenticated user from the request context.
// It's used as opposed to context.WithValue(ctx, "user", user) to avoid potential key collisions.
const UserClaimsContextKey contextKey = "user_claims"

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, UserClaimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*Claims)
	if !ok || claims == nil {
		return nil, ErrUnauthenticated
	}
	return claims, nil
}
