package events

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

var ErrEventBusClosed = errors.New("event bus is closed")

type InMemoryEventBus struct {
	mu       sync.RWMutex
	subs     map[string]map[uint64]Handler
	nextSub  uint64
	queue    chan Event
	workers  sync.WaitGroup
	closeBus sync.Once
	closed   bool
}

func NewInMemoryEventBus(workerCount, queueSize int) *InMemoryEventBus {
	if workerCount <= 0 {
		workerCount = 1
	}
	if queueSize <= 0 {
		queueSize = 128
	}

	b := &InMemoryEventBus{
		subs:  make(map[string]map[uint64]Handler),
		queue: make(chan Event, queueSize),
	}

	for i := 0; i < workerCount; i++ {
		b.workers.Add(1)
		go b.runWorker()
	}

	return b
}

func (b *InMemoryEventBus) Publish(ctx context.Context, event Event) error {
	if event == nil {
		return errors.New("event is nil")
	}

	meta := event.Metadata()
	if meta.Type == "" {
		return errors.New("event metadata type is required")
	}

	b.mu.RLock()
	closed := b.closed
	b.mu.RUnlock()
	if closed {
		return ErrEventBusClosed
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.queue <- event:
		return nil
	}
}

func (b *InMemoryEventBus) Subscribe(eventType string, handler Handler) (func(), error) {
	if eventType == "" {
		return nil, errors.New("eventType is required")
	}
	if handler == nil {
		return nil, errors.New("handler is required")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, ErrEventBusClosed
	}

	b.nextSub++
	subID := b.nextSub
	if _, ok := b.subs[eventType]; !ok {
		b.subs[eventType] = make(map[uint64]Handler)
	}
	b.subs[eventType][subID] = handler

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if handlers, ok := b.subs[eventType]; ok {
			delete(handlers, subID)
			if len(handlers) == 0 {
				delete(b.subs, eventType)
			}
		}
	}, nil
}

func (b *InMemoryEventBus) Close(ctx context.Context) error {
	var closeErr error

	b.closeBus.Do(func() {
		b.mu.Lock()
		b.closed = true
		close(b.queue)
		b.mu.Unlock()
	})

	done := make(chan struct{})
	go func() {
		b.workers.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		closeErr = ctx.Err()
	case <-done:
	}

	return closeErr
}

func (b *InMemoryEventBus) runWorker() {
	defer b.workers.Done()

	for event := range b.queue {
		meta := event.Metadata()
		handlers := b.snapshot(meta.Type)
		for _, handler := range handlers {
			handlerCtx := context.WithoutCancel(context.Background())
			func() {
				defer func() {
					if r := recover(); r != nil {
						slog.Error("Event handler panicked", "eventType", meta.Type, "eventID", meta.ID, "panic", fmt.Sprint(r))
					}
				}()
				if err := handler.Handle(handlerCtx, event); err != nil {
					slog.Warn("Event handler failed", "eventType", meta.Type, "eventID", meta.ID, "err", err)
				}
			}()
		}
	}
}

func (b *InMemoryEventBus) snapshot(eventType string) []Handler {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlersByID := b.subs[eventType]
	handlers := make([]Handler, 0, len(handlersByID))
	for _, h := range handlersByID {
		handlers = append(handlers, h)
	}
	return handlers
}
