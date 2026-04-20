package audit

import (
	"context"
	"fmt"

	"foodplanner/internal/db"
	"foodplanner/internal/events"
)

// A type of function that maps an event to an audit entry
type AuditEntryMapper func(event events.Event) (*AuditEntry, error)

type AuditProjector struct {
	auditService *AuditService
	mappers      map[string]AuditEntryMapper
}

func NewAuditProjector(auditService *AuditService, mappers map[string]AuditEntryMapper) *AuditProjector {
	return &AuditProjector{
		auditService: auditService,
		mappers:      mappers,
	}
}

// One handle method that reduces boilerplate by centralizing the logic of looking up the right mapper and logging the resulting entry
func (p *AuditProjector) Handle(ctx context.Context, tx db.DBTX, event events.Event) error {
	if p == nil || p.auditService == nil || event == nil {
		return nil
	}

	mapper, ok := p.mappers[event.Metadata().Type]
	if !ok || mapper == nil {
		return nil
	}

	entry, err := mapper(event)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}

	return p.auditService.Log(context.WithoutCancel(ctx), tx, entry)
}

func DefaultAuditMappers() map[string]AuditEntryMapper {
	return map[string]AuditEntryMapper{
		events.RecipeCreatedEventType:     mapRecipeCreatedEvent,
		events.RecipeUpdatedEventType:     mapRecipeUpdatedEvent,
		events.UserSignedUpType:           mapUserSignedUpEvent,
		events.UserSignedInType:           mapUserSignedInEvent,
		events.UserSigninFailedType:       mapUserSigninFailedEvent,
		events.GraphQLRequestRejectedType: mapGraphQLRequestRejectedEvent,
		events.RateLimitExceededType:      mapRateLimitExceededEvent,
	}
}

func mapRecipeCreatedEvent(event events.Event) (*AuditEntry, error) {
	e, ok := event.(events.RecipeCreatedEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type for recipe created mapper: %T", event)
	}

	return NewRecipeCreatedEvent(
		e.Meta.CorrelationID,
		e.UserID,
		e.RecipeID,
		e.VersionID,
		e.IPAddress,
		e.UserAgent,
		e.Meta.OccurredAtUTC,
	)
}

func mapRecipeUpdatedEvent(event events.Event) (*AuditEntry, error) {
	e, ok := event.(events.RecipeUpdatedEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type for recipe updated mapper: %T", event)
	}

	return NewRecipeUpdatedEvent(
		e.Meta.CorrelationID,
		e.UserID,
		e.RecipeID,
		e.VersionID,
		e.IPAddress,
		e.UserAgent,
		e.Meta.OccurredAtUTC,
	)
}

func mapUserSignedUpEvent(event events.Event) (*AuditEntry, error) {
	e, ok := event.(events.UserSignedUpEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type for user signed up mapper: %T", event)
	}

	return NewUserSignupEvent(
		e.Meta.CorrelationID,
		e.UserID,
		e.Username,
		e.IPAddress,
	)
}

func mapUserSignedInEvent(event events.Event) (*AuditEntry, error) {
	e, ok := event.(events.UserSignedInEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type for user signed in mapper: %T", event)
	}

	return NewUserSigninEvent(
		e.Meta.CorrelationID,
		e.UserID,
		e.Username,
		e.IPAddress,
		e.UserAgent,
	)
}

func mapUserSigninFailedEvent(event events.Event) (*AuditEntry, error) {
	e, ok := event.(events.UserSigninFailedEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type for user signin failed mapper: %T", event)
	}

	return NewUserSigninFailureEvent(
		e.Meta.CorrelationID,
		e.UserID,
		e.Email,
		e.IPAddress,
		e.UserAgent,
		e.FailureReason,
	)
}

func mapGraphQLRequestRejectedEvent(event events.Event) (*AuditEntry, error) {
	e, ok := event.(events.GraphQLRequestRejectedEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type for graphql request rejected mapper: %T", event)
	}

	return NewGraphQLRequestRejectedEvent(
		e.Meta.CorrelationID,
		e.Meta.ActorUserID,
		e.OperationName,
		e.OperationType,
		e.IPAddress,
		e.UserAgent,
		e.Path,
		e.Reason,
		e.QueryHash,
		e.MaxComplexity,
		e.PresentedMessage,
	)
}

func mapRateLimitExceededEvent(event events.Event) (*AuditEntry, error) {
	e, ok := event.(events.RateLimitExceededEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type for rate limit exceeded mapper: %T", event)
	}

	return NewRateLimitExceededEvent(
		e.Meta.CorrelationID,
		e.Meta.ActorUserID,
		e.Subject,
		string(e.SubjectType),
		e.IPAddress,
		e.UserAgent,
		e.Method,
		e.Path,
		e.Limit,
		e.Count,
		e.WindowSeconds,
		e.RetryAfterSeconds,
		e.Meta.OccurredAtUTC,
	)
}
