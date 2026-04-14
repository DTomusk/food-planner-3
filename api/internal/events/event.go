package events

import (
	"time"

	"github.com/google/uuid"
)

// Metadata is transport-level information shared by all events.
type Metadata struct {
	ID            uuid.UUID
	Type          string
	Version       int
	OccurredAtUTC time.Time
	CorrelationID uuid.UUID
	CausationID   *uuid.UUID
	TraceID       string
	ActorUserID   *uuid.UUID
	Source        string
}

type Event interface {
	Metadata() Metadata
}
