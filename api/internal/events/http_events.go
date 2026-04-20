package events

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const RateLimitExceededType = "http.rate_limit_exceeded"

type RateLimitSubjectType string

const (
	SubjectTypeUser RateLimitSubjectType = "user"
	SubjectTypeIP   RateLimitSubjectType = "ip"
)

type RateLimitExceededEvent struct {
	Meta              Metadata
	Subject           string // e.g. "user:1234" or "ip:192.168.0.1"
	SubjectType       RateLimitSubjectType
	IPAddress         string
	UserAgent         string
	Method            string
	Path              string
	Limit             int
	Count             int
	WindowSeconds     int
	RetryAfterSeconds int
}

func (e RateLimitExceededEvent) Metadata() Metadata {
	return e.Meta
}

func NewRateLimitExceededEvent(
	correlationID uuid.UUID,
	subject, ipAddress, userAgent, method, path string,
	limit, count, windowSeconds, retryAfterSeconds int,
) RateLimitExceededEvent {
	return RateLimitExceededEvent{
		Meta: Metadata{
			ID:            uuid.New(),
			Type:          RateLimitExceededType,
			Version:       1,
			OccurredAtUTC: time.Now().UTC(),
			CorrelationID: correlationID,
			ActorUserID:   actorUserIDFromSubject(subject),
			Source:        GraphQLServerSource,
		},
		Subject:           subject,
		SubjectType:       rateLimitSubjectTypeFromSubject(subject),
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		Method:            method,
		Path:              path,
		Limit:             limit,
		Count:             count,
		WindowSeconds:     windowSeconds,
		RetryAfterSeconds: retryAfterSeconds,
	}
}

func rateLimitSubjectTypeFromSubject(subject string) RateLimitSubjectType {
	if strings.HasPrefix(subject, "user:") {
		return SubjectTypeUser
	}
	return SubjectTypeIP
}

func actorUserIDFromSubject(subject string) *uuid.UUID {
	if !strings.HasPrefix(subject, "user:") {
		return nil
	}

	userID := strings.TrimPrefix(subject, "user:")
	parsed, err := uuid.Parse(userID)
	if err != nil {
		return nil
	}

	return &parsed
}
