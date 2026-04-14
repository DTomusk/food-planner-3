package events

import "context"

// Publisher is the narrow interface that domain services depend on.
// Satisfied by InMemoryEventBus and RedisPublisher.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// EventBus extends Publisher with subscription and lifecycle management.
// Only the composition root (main.go) and the worker process need this.
type EventBus interface {
	Publisher
	Subscribe(eventType string, handler Handler) (func(), error)
	Close(ctx context.Context) error
}
