package auth

import "context"

type contextKey string

// UserContextKey is the context key used to store and retrieve the authenticated user from the request context.
// It's used as opposed to context.WithValue(ctx, "user", user) to avoid potential key collisions.
const UserClaimsContextKey contextKey = "user_claims"
const InvalidAuthTokenContextKey contextKey = "invalid_auth_token"

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, UserClaimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*Claims)
	return claims, ok && claims != nil
}

// Mark in context that the auth token is invalid, so that resolvers can return unauthenticated errors instead of just treating it as no user
func ContextWithInvalidAuthToken(ctx context.Context) context.Context {
	return context.WithValue(ctx, InvalidAuthTokenContextKey, true)
}

func InvalidAuthTokenFromContext(ctx context.Context) bool {
	isInvalid, ok := ctx.Value(InvalidAuthTokenContextKey).(bool)
	return ok && isInvalid
}
