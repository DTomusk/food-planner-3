package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditEntry struct {
	ID uuid.UUID
	// For grouping related changes together i.e. caused by the same action
	CorrelationID uuid.UUID
	// Who (user id) did the action, null if system action
	ActorID      *uuid.UUID
	ResourceType ResourceType
	// ID of affected resource, null if not applicable
	ResourceID *uuid.UUID
	Action     Action
	CreatedAt  time.Time
	Result     Result
	// Snapshot of affected state before and after action
	OldState json.RawMessage
	NewState json.RawMessage
	// Optional justification for action
	Reason  *string
	Context json.RawMessage
}

type Action string
type ResourceType string
type Result string

const (
	ActionUserSignup Action = "user.signup"
	ActionUserSignin Action = "user.signin"

	ResourceTypeUser ResourceType = "user"

	ResultSuccess Result = "success"
	ResultFailure Result = "failure"
)
