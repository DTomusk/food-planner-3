package events

import (
	"context"
	"errors"
	"foodplanner/internal/logging"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisWorker struct {
	client       *redis.Client
	stream       string
	group        string
	consumerName string
	// Gets event type from json, which then goes to the correct handlers
	registry *Registry
	// Maps event types to lists of handlers so an event can go to multiple
	handlers map[string][]Handler
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

func NewRedisWorker(client *redis.Client, stream, group, consumerName string, registry *Registry) *RedisWorker {
	return &RedisWorker{
		client:       client,
		stream:       stream,
		group:        group,
		consumerName: consumerName,
		registry:     registry,
		handlers:     make(map[string][]Handler),
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
			// Block waits 5 seconds if there are no messages
			Block: 5 * time.Second,
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
				// Handle message
				// What happens if this fails?
				if err := w.handleMessage(ctx, message); err != nil {
					logger.Error("failed to process message", "messageID", message.ID, "error", err)
					continue
				}

				// Acknowledge with a context detached from cancellation so a successful
				// handler can still be acked during shutdown.
				ackCtx, cancelAck := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
				if err := w.client.XAck(ackCtx, w.stream, w.group, message.ID).Err(); err != nil {
					logger.Error("failed to ack message", "messageID", message.ID, "error", err)
				}
				cancelAck()
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
		return err
	}

	event, err := UnmarshalEvent(env, w.registry)
	if err != nil {
		return err
	}

	handlers, ok := w.handlers[event.Metadata().Type]
	if !ok {
		return nil
	}

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
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
