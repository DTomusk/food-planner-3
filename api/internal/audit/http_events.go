package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func NewRateLimitExceededEvent(
	correlationID uuid.UUID,
	actorUserID *uuid.UUID,
	subject string,
	subjectType string,
	ipAddress string,
	userAgent string,
	method string,
	path string,
	limit int,
	count int,
	windowSeconds int,
	retryAfterSeconds int,
	occurredAtUTC time.Time,
) (*AuditEntry, error) {
	createdAt := time.Now().UTC()
	reason := "rate_limit_exceeded"

	contextData, err := json.Marshal(struct {
		Source            string `json:"source"`
		Operation         string `json:"operation"`
		Subject           string `json:"subject"`
		SubjectType       string `json:"subject_type"`
		IPAddress         string `json:"ip_address"`
		UserAgent         string `json:"user_agent"`
		Method            string `json:"method"`
		Path              string `json:"path"`
		Limit             int    `json:"limit"`
		Count             int    `json:"count"`
		WindowSeconds     int    `json:"window_seconds"`
		RetryAfterSeconds int    `json:"retry_after_seconds"`
		OccurredAtUTC     string `json:"occurred_at_utc"`
	}{
		Source:            "http",
		Operation:         "rate_limit_exceeded",
		Subject:           subject,
		SubjectType:       subjectType,
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		Method:            method,
		Path:              path,
		Limit:             limit,
		Count:             count,
		WindowSeconds:     windowSeconds,
		RetryAfterSeconds: retryAfterSeconds,
		OccurredAtUTC:     occurredAtUTC.UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return nil, err
	}

	return &AuditEntry{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		ActorID:       actorUserID,
		ResourceType:  ResourceTypeHTTPRequest,
		Action:        ActionRateLimitExceeded,
		Result:        ResultFailure,
		CreatedAt:     createdAt,
		Reason:        &reason,
		Context:       contextData,
	}, nil
}
