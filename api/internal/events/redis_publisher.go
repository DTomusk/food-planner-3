package events

import (
	"context"
	"foodplanner/internal/logging"

	"github.com/redis/go-redis/v9"
)

// Publisher that publishes to Redis streams
type RedisPublisher struct {
	client *redis.Client
	stream string
	maxLen int64
}

func NewRedisPublisher(client *redis.Client, stream string, maxLen int64) *RedisPublisher {
	return &RedisPublisher{
		client: client,
		stream: stream,
		maxLen: maxLen,
	}
}

func (p *RedisPublisher) Publish(ctx context.Context, event Event) error {
	logger := logging.FromContext(ctx)
	meta := event.Metadata()
	logger.Info("Publishing event to Redis stream", "stream", p.stream, "eventType", meta.Type, "eventID", meta.ID)
	env, err := MarshalEvent(event)
	if err != nil {
		incrementPublishFailure()
		logger.Error("Failed to marshal event", "stream", p.stream, "eventType", meta.Type, "eventID", meta.ID, "error", err)
		return err
	}
	data, err := MarshalEnvelope(env)
	if err != nil {
		incrementPublishFailure()
		logger.Error("Failed to marshal envelope", "stream", p.stream, "eventType", meta.Type, "eventID", meta.ID, "error", err)
		return err
	}
	// XAdd appens the serialized envelope to the redis stream
	err = p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		// Prevents stream from growing indefinitely
		MaxLen: p.maxLen,
		Approx: true,
		Values: map[string]any{
			"data": string(data),
		},
	}).Err()
	if err != nil {
		incrementPublishFailure()
		logger.Error("Failed to publish event to Redis stream", "stream", p.stream, "eventType", meta.Type, "eventID", meta.ID, "error", err)
		return err
	}
	incrementPublishSuccess()

	return nil
}
