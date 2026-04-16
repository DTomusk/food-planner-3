package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserSignedUpType     = "user.signed_up"
	UserSignedInType     = "user.signed_in"
	UserSigninFailedType = "user.signin_failed"
)

type UserSignedUpEvent struct {
	Meta      Metadata
	UserID    uuid.UUID
	Username  string
	Email     string
	IPAddress string
}

func (e UserSignedUpEvent) Metadata() Metadata {
	return e.Meta
}

func NewUserSignedUpEvent(correlationID, userID uuid.UUID, username, email, ipAddress string) UserSignedUpEvent {
	actorID := userID

	return UserSignedUpEvent{
		Meta: Metadata{
			ID:            uuid.New(),
			Type:          UserSignedUpType,
			Version:       1,
			OccurredAtUTC: time.Now().UTC(),
			CorrelationID: correlationID,
			ActorUserID:   &actorID,
			Source:        "auth.service",
		},
		UserID:    userID,
		Username:  username,
		Email:     email,
		IPAddress: ipAddress,
	}
}

type UserSignedInEvent struct {
	Meta      Metadata
	UserID    uuid.UUID
	Username  string
	Email     string
	IPAddress string
	UserAgent string
}

func (e UserSignedInEvent) Metadata() Metadata {
	return e.Meta
}

func NewUserSignedInEvent(correlationID, userID uuid.UUID, username, email, ipAddress, userAgent string) UserSignedInEvent {
	actorID := userID

	return UserSignedInEvent{
		Meta: Metadata{
			ID:            uuid.New(),
			Type:          UserSignedInType,
			Version:       1,
			OccurredAtUTC: time.Now().UTC(),
			CorrelationID: correlationID,
			ActorUserID:   &actorID,
			Source:        "auth.service",
		},
		UserID:    userID,
		Username:  username,
		Email:     email,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}

type UserSigninFailedEvent struct {
	Meta          Metadata
	UserID        *uuid.UUID
	Email         string
	IPAddress     string
	UserAgent     string
	FailureReason string
}

func (e UserSigninFailedEvent) Metadata() Metadata {
	return e.Meta
}

func NewUserSigninFailedEvent(correlationID uuid.UUID, userID *uuid.UUID, email, ipAddress, userAgent, failureReason string) UserSigninFailedEvent {
	return UserSigninFailedEvent{
		Meta: Metadata{
			ID:            uuid.New(),
			Type:          UserSigninFailedType,
			Version:       1,
			OccurredAtUTC: time.Now().UTC(),
			CorrelationID: correlationID,
			Source:        "auth.service",
		},
		UserID:        userID,
		Email:         email,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		FailureReason: failureReason,
	}
}
