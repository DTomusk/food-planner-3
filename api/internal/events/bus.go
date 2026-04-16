package events

import "context"

// Publisher is the narrow interface that domain services depend on for publishing events
// Satisfied by InMemoryEventBus and RedisPublisher.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// EventBus extends Publisher with subscription and lifecycle management.
// Currently not used in production as we have redis publisher and redis worker
type EventBus interface {
	Publisher
	Subscribe(eventType string, handler Handler) (func(), error)
	Close(ctx context.Context) error
}
