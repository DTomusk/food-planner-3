package events

import (
	"context"

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
	env, err := MarshalEvent(event)
	if err != nil {
		return err
	}
	data, err := MarshalEnvelope(env)
	if err != nil {
		return err
	}
	// XAdd appens the serialized envelope to the redis stream
	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		// Prevents stream from growing indefinitely
		MaxLen: p.maxLen,
		Approx: true,
		Values: map[string]any{
			"data": string(data),
		},
	}).Err()
}
