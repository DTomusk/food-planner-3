package events

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisPublisher_Publish_WritesEnvelopeToStream(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	publisher := NewRedisPublisher(client, "events.main", 1000)
	e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "demo@example.com", "127.0.0.1")

	err := publisher.Publish(context.Background(), e)
	require.NoError(t, err)

	entries, err := client.XRange(context.Background(), "events.main", "-", "+").Result()
	require.NoError(t, err)
	require.Len(t, entries, 1)

	raw, ok := entries[0].Values["data"].(string)
	require.True(t, ok)
	env, err := UnmarshalEnvelope([]byte(raw))
	require.NoError(t, err)
	require.Equal(t, UserSignedUpType, env.Type)
}

func TestRedisPublisher_Publish_InvalidEvent_ReturnsError(t *testing.T) {
	s := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	publisher := NewRedisPublisher(client, "events.main", 1000)

	err := publisher.Publish(context.Background(), UserSignedUpEvent{})
	require.ErrorIs(t, err, ErrEmptyEventType)
}

func TestRedisPublisher_Publish_RedisUnavailable_ReturnsError(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	t.Cleanup(func() { _ = client.Close() })

	publisher := NewRedisPublisher(client, "events.main", 1000)
	e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "demo@example.com", "127.0.0.1")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := publisher.Publish(ctx, e)
	require.Error(t, err)
}
