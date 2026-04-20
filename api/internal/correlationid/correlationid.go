package correlationid

import (
	"context"

	"github.com/google/uuid"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

// WithContext returns a new context with the given correlation ID stored in it.
func WithContext(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKey, id)
}

// FromContext returns the correlation ID stored in the context.
// If none is present, a new UUID is generated as a fallback.
func FromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(contextKey).(uuid.UUID); ok {
		return id
	}
	return uuid.New()
}
