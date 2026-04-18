package events

import (
	"time"

	"github.com/google/uuid"
)

const GraphQLRequestRejectedType = "graphql.request_rejected"

type GraphQLRequestRejectedEvent struct {
	Meta             Metadata
	OperationName    string
	OperationType    string
	IPAddress        string
	UserAgent        string
	Path             string
	Reason           string
	QueryHash        string
	MaxComplexity    int
	PresentedMessage string
}

func (e GraphQLRequestRejectedEvent) Metadata() Metadata {
	return e.Meta
}

func NewGraphQLRequestRejectedEvent(
	correlationID uuid.UUID,
	actorUserID *uuid.UUID,
	operationName, operationType, ipAddress, userAgent, path, reason, queryHash string,
	maxComplexity int,
	presentedMessage string,
) GraphQLRequestRejectedEvent {
	return GraphQLRequestRejectedEvent{
		Meta: Metadata{
			ID:            uuid.New(),
			Type:          GraphQLRequestRejectedType,
			Version:       1,
			OccurredAtUTC: time.Now().UTC(),
			CorrelationID: correlationID,
			ActorUserID:   actorUserID,
			Source:        GraphQLServerSource,
		},
		OperationName:    operationName,
		OperationType:    operationType,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		Path:             path,
		Reason:           reason,
		QueryHash:        queryHash,
		MaxComplexity:    maxComplexity,
		PresentedMessage: presentedMessage,
	}
}
