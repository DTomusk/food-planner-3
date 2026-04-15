package events

import (
	"context"
	"foodplanner/internal/db"
)

// Concrete handlers live in their respective packages but must implement this interface
type Handler interface {
	Handle(ctx context.Context, tx db.DBTX, event Event) error
}

type HandlerFunc func(ctx context.Context, tx db.DBTX, event Event) error

func (f HandlerFunc) Handle(ctx context.Context, tx db.DBTX, event Event) error {
	return f(ctx, tx, event)
}
