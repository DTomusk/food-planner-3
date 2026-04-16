package events

import (
	"context"
	"database/sql"
	"errors"
	"foodplanner/internal/db"
	"foodplanner/internal/testutil"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func newWorkerForTx(tx *sql.Tx, client *redis.Client, registry *Registry) *RedisWorker {
	repo := NewEventsRepo()
	txRunner := testutil.NewTestTxRunner(tx)
	expectedVersions := map[string]int{
		UserSignedUpType: 1,
	}
	return NewRedisWorker(client, "events.main", "worker-group", "worker-1", registry, expectedVersions, repo, txRunner)
}

func TestRegisterHandler_EmptyEventType_ReturnsError(t *testing.T) {
	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", NewRegistry(), nil, nil, nil)

	err := worker.RegisterHandler("", HandlerFunc(func(context.Context, db.DBTX, Event) error { return nil }))
	require.ErrorIs(t, err, ErrEmptyEventType)
}

func TestRegisterHandler_NilHandler_ReturnsError(t *testing.T) {
	worker := NewRedisWorker(nil, "events.main", "worker-group", "worker-1", NewRegistry(), nil, nil, nil)

	err := worker.RegisterHandler(UserSignedUpType, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "handler cannot be nil")
}

func TestRun_AcksUnsupportedVersionMessageWithoutProcessing(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := miniredis.RunT(t)
		client := redis.NewClient(&redis.Options{Addr: s.Addr()})
		t.Cleanup(func() { _ = client.Close() })

		registry := NewRegistry()
		registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

		worker := newWorkerForTx(tx, client, registry)
		require.NoError(t, worker.ensureGroup(context.Background()))

		e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "user@example.com", "127.0.0.1")
		e.Meta.Version = 2
		env, err := MarshalEvent(e)
		require.NoError(t, err)
		data, err := MarshalEnvelope(env)
		require.NoError(t, err)

		_, err = client.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "events.main",
			Values: map[string]any{"data": string(data)},
		}).Result()
		require.NoError(t, err)

		runCtx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		called := 0
		require.NoError(t, worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(context.Context, db.DBTX, Event) error {
			called++
			return nil
		})))

		runDone := make(chan error, 1)
		go func() {
			runDone <- worker.Run(runCtx)
		}()

		deadline := time.Now().Add(3 * time.Second)
		for {
			pending, err := client.XPending(context.Background(), "events.main", "worker-group").Result()
			require.NoError(t, err)
			if pending.Count == 0 {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("unsupported-version message was not acknowledged")
			}
			time.Sleep(20 * time.Millisecond)
		}

		cancel()
		select {
		case err := <-runDone:
			require.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("worker did not stop")
		}

		require.Equal(t, 0, called)
	})
}

func TestHandleMessage_DispatchesToRegisteredHandler(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		registry := NewRegistry()
		registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

		worker := newWorkerForTx(tx, nil, registry)

		called := 0
		err := worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(_ context.Context, tx db.DBTX, event Event) error {
			require.NotNil(t, tx)
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
	})
}

func TestHandleMessage_NoHandlersForEventType_ReturnsNil(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		registry := NewRegistry()
		registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })
		worker := newWorkerForTx(tx, nil, registry)

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
	})
}

func TestHandleMessage_MissingDataField_ReturnsError(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		worker := newWorkerForTx(tx, nil, NewRegistry())

		err := worker.handleMessage(context.Background(), redis.XMessage{
			ID:     "1-0",
			Values: map[string]any{"wrong": "value"},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "message missing 'data' field")
	})
}

func TestHandleMessage_ByteData_DispatchesToHandler(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		registry := NewRegistry()
		registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

		worker := newWorkerForTx(tx, nil, registry)

		called := 0
		err := worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(_ context.Context, tx db.DBTX, event Event) error {
			require.NotNil(t, tx)
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
			Values: map[string]any{"data": data},
		})
		require.NoError(t, err)
		require.Equal(t, 1, called)
	})
}

func TestHandleMessage_HandlerError_BubblesUp(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		registry := NewRegistry()
		registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })
		worker := newWorkerForTx(tx, nil, registry)

		handlerErr := errors.New("handler failed")
		err := worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(context.Context, db.DBTX, Event) error {
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
	})
}

func TestEnsureGroup_CreateAndBusyGroup_ReturnsNil(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := miniredis.RunT(t)
		client := redis.NewClient(&redis.Options{Addr: s.Addr()})
		t.Cleanup(func() { _ = client.Close() })

		worker := newWorkerForTx(tx, client, NewRegistry())

		err := worker.ensureGroup(context.Background())
		require.NoError(t, err)

		err = worker.ensureGroup(context.Background())
		require.NoError(t, err)
	})
}

func TestEnsureGroup_RedisError_BubblesUp(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
		t.Cleanup(func() { _ = client.Close() })

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		worker := newWorkerForTx(tx, client, NewRegistry())
		err := worker.ensureGroup(ctx)
		require.Error(t, err)
	})
}

func TestRun_AcksMessageOnSuccessfulHandling(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := miniredis.RunT(t)
		client := redis.NewClient(&redis.Options{Addr: s.Addr()})
		t.Cleanup(func() { _ = client.Close() })

		registry := NewRegistry()
		registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

		worker := newWorkerForTx(tx, client, registry)
		require.NoError(t, worker.ensureGroup(context.Background()))

		e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "user@example.com", "127.0.0.1")
		env, err := MarshalEvent(e)
		require.NoError(t, err)
		data, err := MarshalEnvelope(env)
		require.NoError(t, err)

		_, err = client.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "events.main",
			Values: map[string]any{"data": string(data)},
		}).Result()
		require.NoError(t, err)

		runCtx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		handled := make(chan struct{}, 1)
		require.NoError(t, worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(context.Context, db.DBTX, Event) error {
			handled <- struct{}{}
			cancel()
			return nil
		})))

		runDone := make(chan error, 1)
		go func() {
			runDone <- worker.Run(runCtx)
		}()

		select {
		case <-handled:
		case <-time.After(3 * time.Second):
			t.Fatal("handler was not invoked")
		}

		select {
		case err := <-runDone:
			require.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("worker did not stop")
		}

		pending, err := client.XPending(context.Background(), "events.main", "worker-group").Result()
		require.NoError(t, err)
		require.EqualValues(t, 0, pending.Count)
	})
}

func TestRun_DoesNotAckMessageOnHandlerFailure(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := miniredis.RunT(t)
		client := redis.NewClient(&redis.Options{Addr: s.Addr()})
		t.Cleanup(func() { _ = client.Close() })

		registry := NewRegistry()
		registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

		worker := newWorkerForTx(tx, client, registry)
		require.NoError(t, worker.ensureGroup(context.Background()))

		e := NewUserSignedUpEvent(uuid.New(), uuid.New(), "demo", "user@example.com", "127.0.0.1")
		env, err := MarshalEvent(e)
		require.NoError(t, err)
		data, err := MarshalEnvelope(env)
		require.NoError(t, err)

		_, err = client.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "events.main",
			Values: map[string]any{"data": string(data)},
		}).Result()
		require.NoError(t, err)

		runCtx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		handled := make(chan struct{}, 1)
		require.NoError(t, worker.RegisterHandler(UserSignedUpType, HandlerFunc(func(context.Context, db.DBTX, Event) error {
			handled <- struct{}{}
			return errors.New("handler failed")
		})))

		runDone := make(chan error, 1)
		go func() {
			runDone <- worker.Run(runCtx)
		}()

		select {
		case <-handled:
			cancel()
		case <-time.After(3 * time.Second):
			t.Fatal("handler was not invoked")
		}

		select {
		case err := <-runDone:
			require.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("worker did not stop")
		}

		pending, err := client.XPending(context.Background(), "events.main", "worker-group").Result()
		require.NoError(t, err)
		require.EqualValues(t, 1, pending.Count)
	})
}
