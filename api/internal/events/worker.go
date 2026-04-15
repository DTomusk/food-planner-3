package events

import (
	"context"
	"database/sql"
	"errors"
	"foodplanner/internal/db"
	"foodplanner/internal/logging"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisWorker struct {
	client       *redis.Client
	stream       string
	group        string
	consumerName string
	readBlock    time.Duration
	// Gets event type from json, which then goes to the correct handlers
	registry *Registry
	// Maps event types to lists of handlers so an event can go to multiple
	handlers map[string][]Handler
	// Optional store for idempotency dedupe markers.
	eventsRepo *EventsRepo
	txRunner   db.TxRunner
}

type eventProcessingError struct {
	eventType string
	eventID   string
	err       error
}

func (e *eventProcessingError) Error() string {
	return e.err.Error()
}

func (e *eventProcessingError) Unwrap() error {
	return e.err
}

func newEventProcessingError(eventType, eventID string, err error) error {
	if err == nil {
		return nil
	}
	return &eventProcessingError{
		eventType: eventType,
		eventID:   eventID,
		err:       err,
	}
}

func (w *RedisWorker) RegisterHandler(eventType string, handler Handler) error {
	if eventType == "" {
		return ErrEmptyEventType
	}
	if handler == nil {
		return errors.New("handler cannot be nil")
	}

	w.handlers[eventType] = append(w.handlers[eventType], handler)
	return nil
}

func NewRedisWorker(client *redis.Client, stream, group, consumerName string, registry *Registry, eventsRepo *EventsRepo, txRunner db.TxRunner) *RedisWorker {
	return &RedisWorker{
		client:       client,
		stream:       stream,
		group:        group,
		consumerName: consumerName,
		readBlock:    1 * time.Second,
		registry:     registry,
		handlers:     make(map[string][]Handler),
		eventsRepo:   eventsRepo,
		txRunner:     txRunner,
	}
}

func (w *RedisWorker) Run(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Info("Starting Redis worker", "stream", w.stream, "group", w.group, "consumer", w.consumerName)
	err := w.ensureGroup(ctx)
	if err != nil {
		logger.Error("Failed to ensure Redis group", "error", err)
		return err
	}

	// Infinite loop to read Redis stream and process messages
	for {
		select {
		// Shutdown gracefully if context done
		case <-ctx.Done():
			logger.Info("Shutting down Redis worker", "stream", w.stream, "group", w.group, "consumer", w.consumerName)
			return nil
		default:
		}

		// Reads messages from Redis stream
		streams, err := w.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    w.group,
			Consumer: w.consumerName,
			Streams:  []string{w.stream, ">"},
			// Count controls how many messages are read at once
			Count: 10,
			// Keep this short so shutdown remains responsive.
			Block: w.readBlock,
		}).Result()
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			if errors.Is(err, redis.Nil) {
				continue
			}
			logger.Error("Failed to read from Redis stream", "error", err)
			return err
		}

		// Streams contain messages, each message contains an event to process
		for _, stream := range streams {
			for _, message := range stream.Messages {
				msgCtx, cancelMsg := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
				// Handle message
				// What happens if this fails?
				if err := w.handleMessage(msgCtx, message); err != nil {
					incrementWorkerHandlerFailure()
					cancelMsg()
					var procErr *eventProcessingError
					if errors.As(err, &procErr) {
						logger.Error("failed to process message", "messageID", message.ID, "eventType", procErr.eventType, "eventID", procErr.eventID, "error", procErr.err)
					} else {
						logger.Error("failed to process message", "messageID", message.ID, "eventType", "unknown", "eventID", "unknown", "error", err)
					}
					continue
				}

				// Acknowledge with a context detached from cancellation so a successful
				// handler can still be acked during shutdown.
				ackCtx, cancelAck := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
				if err := w.client.XAck(ackCtx, w.stream, w.group, message.ID).Err(); err != nil {
					incrementWorkerAckFailure()
					logger.Error("failed to ack message", "messageID", message.ID, "error", err)
				} else {
					incrementWorkerProcessed()
				}
				cancelAck()
				cancelMsg()
			}
		}
	}
}

// Get event data from json message and pass to relevant handlers
func (w *RedisWorker) handleMessage(ctx context.Context, message redis.XMessage) error {
	var payload []byte
	switch raw := message.Values["data"].(type) {
	case string:
		payload = []byte(raw)
	case []byte:
		payload = raw
	default:
		return errors.New("message missing 'data' field or 'data' is not a string")
	}

	env, err := UnmarshalEnvelope(payload)
	if err != nil {
		return newEventProcessingError("unknown", "unknown", err)
	}

	event, err := UnmarshalEvent(env, w.registry)
	if err != nil {
		return newEventProcessingError(env.Type, "unknown", err)
	}
	eventType := event.Metadata().Type
	eventID := event.Metadata().ID.String()

	handlers, ok := w.handlers[event.Metadata().Type]
	if !ok {
		return nil
	}

	for _, handler := range handlers {
		handlerName := resolveHandlerName(handler)

		// Check idempotency key to see if event has already been processed by handler
		processed, err := w.eventsRepo.checkEventProcessed(ctx, w.txRunner.DB(), eventID, w.group, handlerName)
		if err != nil {
			return newEventProcessingError(eventType, eventID, err)
		}
		// Skip if already processed
		if processed {
			continue
		}

		// Use a transaction to ensure that processing the event and marking it as processed are atomic
		err = w.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
			if err := handler.Handle(ctx, tx, event); err != nil {
				return err
			}

			// Save the fact that the event has been processed by this handler and group
			if err := w.eventsRepo.markEventProcessed(ctx, tx, eventID, w.group, handlerName); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return newEventProcessingError(eventType, eventID, err)
		}
	}

	return nil
}

// Ensure that the group and stream exist
// Try creating
// If already exists, return no error
func (w *RedisWorker) ensureGroup(ctx context.Context) error {
	err := w.client.XGroupCreateMkStream(ctx, w.stream, w.group, "$").Err()
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "BUSYGROUP") {
		return nil
	}
	return err
}

func resolveHandlerName(handler Handler) string {
	v := reflect.ValueOf(handler)
	if !v.IsValid() {
		return "unknown-handler"
	}

	if v.Kind() == reflect.Func {
		if fn := runtime.FuncForPC(v.Pointer()); fn != nil {
			return fn.Name()
		}
	}

	return reflect.TypeOf(handler).String()
}
