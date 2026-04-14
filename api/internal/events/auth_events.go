package events

import (
	"time"

	"github.com/google/uuid"
)

const UserSignedUpType = "user.signed_up"

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
