package events

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler_EmptyEventType_ReturnsError(t *testing.T) {
	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", NewRegistry())

	err := worker.RegisterHandler("", HandlerFunc(func(context.Context, Event) error { return nil }))
	require.ErrorIs(t, err, ErrEmptyEventType)
}

func TestRegisterHandler_NilHandler_ReturnsError(t *testing.T) {
	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", NewRegistry())

	err := worker.RegisterHandler(UserSignedUpType, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "handler cannot be nil")
}

func TestHandleMessage_DispatchesToRegisteredHandler(t *testing.T) {
	registry := NewRegistry()
	registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", registry)

	called := 0
	err := worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(_ context.Context, event Event) error {
		called++
		signup, ok := event.(UserSignedUpEvent)
		require.True(t, ok)
		require.Equal(t, "user@example.com", signup.Email)
		return nil
	}))
	require.NoError(t, err)

	e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "user@example.com", "127.0.0.1")
	env, err := MarshalEvent(e)
	require.NoError(t, err)
	data, err := MarshalEnvelope(env)
	require.NoError(t, err)

	err = worker.handleMessage(context.Background(), redis.XMessage{
		ID:     "1-0",
		Values: map[string]any{"data": string(data)},
	})
	require.NoError(t, err)
	require.Equal(t, 1, called)
}

func TestHandleMessage_NoHandlersForEventType_ReturnsNil(t *testing.T) {
	registry := NewRegistry()
	registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })
	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", registry)

	e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "user@example.com", "127.0.0.1")
	env, err := MarshalEvent(e)
	require.NoError(t, err)
	data, err := MarshalEnvelope(env)
	require.NoError(t, err)

	err = worker.handleMessage(context.Background(), redis.XMessage{
		ID:     "1-0",
		Values: map[string]any{"data": string(data)},
	})
	require.NoError(t, err)
}

func TestHandleMessage_MissingDataField_ReturnsError(t *testing.T) {
	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", NewRegistry())

	err := worker.handleMessage(context.Background(), redis.XMessage{
		ID:     "1-0",
		Values: map[string]any{"wrong": "value"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "message missing 'data' field")
}

func TestHandleMessage_HandlerError_BubblesUp(t *testing.T) {
	registry := NewRegistry()
	registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })
	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", registry)

	handlerErr := errors.New("handler failed")
	err := worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(context.Context, Event) error {
		return handlerErr
	}))
	require.NoError(t, err)

	e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "user@example.com", "127.0.0.1")
	env, err := MarshalEvent(e)
	require.NoError(t, err)
	data, err := MarshalEnvelope(env)
	require.NoError(t, err)

	err = worker.handleMessage(context.Background(), redis.XMessage{
		ID:     "1-0",
		Values: map[string]any{"data": string(data)},
	})
	require.ErrorIs(t, err, handlerErr)
}
