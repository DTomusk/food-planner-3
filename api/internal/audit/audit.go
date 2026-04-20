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
	ActionRecipeCreated          Action = "recipe.created"
	ActionRecipeUpdated          Action = "recipe.updated"
	ActionGraphQLRequestRejected Action = "graphql.request_rejected"
	ActionRateLimitExceeded      Action = "http.rate_limit_exceeded"
	ActionUserSignup             Action = "user.signup"
	ActionUserSignin             Action = "user.signin"

	ResourceTypeGraphQLRequest ResourceType = "graphql_request"
	ResourceTypeHTTPRequest    ResourceType = "http_request"
	ResourceTypeRecipe         ResourceType = "recipe"
	ResourceTypeUser           ResourceType = "user"

	ResultSuccess Result = "success"
	ResultFailure Result = "failure"
)
