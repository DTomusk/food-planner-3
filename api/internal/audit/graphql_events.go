package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func NewGraphQLRequestRejectedEvent(
	correlationID uuid.UUID,
	actorUserID *uuid.UUID,
	operationName, operationType, ipAddress, userAgent, path, reason, queryHash string,
	maxComplexity int,
	presentedMessage string,
) (*AuditEntry, error) {
	createdAt := time.Now().UTC()

	contextData, err := json.Marshal(struct {
		Source           string `json:"source"`
		Operation        string `json:"operation"`
		OperationType    string `json:"operation_type"`
		IPAddress        string `json:"ip_address"`
		UserAgent        string `json:"user_agent"`
		Path             string `json:"path"`
		Reason           string `json:"reason"`
		QueryHash        string `json:"query_hash"`
		MaxComplexity    int    `json:"max_complexity"`
		PresentedMessage string `json:"presented_message"`
	}{
		Source:           "graphql",
		Operation:        operationName,
		OperationType:    operationType,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		Path:             path,
		Reason:           reason,
		QueryHash:        queryHash,
		MaxComplexity:    maxComplexity,
		PresentedMessage: presentedMessage,
	})
	if err != nil {
		return nil, err
	}

	return &AuditEntry{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		ActorID:       actorUserID,
		ResourceType:  ResourceTypeGraphQLRequest,
		Action:        ActionGraphQLRequestRejected,
		Result:        ResultFailure,
		CreatedAt:     createdAt,
		Reason:        &reason,
		Context:       contextData,
	}, nil
}
