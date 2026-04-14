package events

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMarshalEvent_NilEvent_ReturnsError(t *testing.T) {
	_, err := MarshalEvent(nil)
	require.ErrorIs(t, err, ErrNilEvent, "expected ErrNilEvent, got %v", err)
}

func TestMarshalEvent_EmptyType_ReturnsError(t *testing.T) {
	e := UserSignedUpEvent{}
	_, err := MarshalEvent(e)
	require.ErrorIs(t, err, ErrEmptyEventType, "expected ErrEmptyEventType, got %v", err)
}

func TestEnvelopeRoundTrip(t *testing.T) {
	correlationID := uuid.New()
	userID := uuid.New()
	e := NewUserSignedUpEvent(correlationID, userID, "demo", "demo@example.com", "127.0.0.1")

	env, err := MarshalEvent(e)
	require.NoError(t, err, "marshal event failed: %v", err)

	data, err := MarshalEnvelope(env)
	require.NoError(t, err, "marshal envelope failed: %v", err)

	decodedEnv, err := UnmarshalEnvelope(data)
	require.NoError(t, err, "unmarshal envelope failed: %v", err)

	registry := NewRegistry()
	registry.Register(UserSignedUpType, func() Event { return &UserSignedUpEvent{} })

	decodedEvent, err := UnmarshalEvent(decodedEnv, registry)
	require.NoError(t, err, "unmarshal event failed: %v", err)

	signup, ok := decodedEvent.(*UserSignedUpEvent)
	require.True(t, ok, "expected *UserSignedUpEvent, got %T", decodedEvent)
	require.Equal(t, userID, signup.UserID, "expected userID %s, got %s", userID, signup.UserID)
	require.Equal(t, correlationID, signup.Meta.CorrelationID, "expected correlationID %s, got %s", correlationID, signup.Meta.CorrelationID)
	require.Equal(t, "demo@example.com", signup.Email, "unexpected email: %s", signup.Email)
}

func TestUnmarshalEvent_UnknownType_ReturnsError(t *testing.T) {
	env := Envelope{Type: "does.not.exist", Payload: []byte(`{"x":1}`)}
	registry := NewRegistry()

	_, err := UnmarshalEvent(env, registry)
	require.ErrorIs(t, err, ErrUnknownEventType, "expected ErrUnknownEventType, got %v", err)
}

func TestUnmarshalEnvelope_EmptyData_ReturnsError(t *testing.T) {
	_, err := UnmarshalEnvelope(nil)
	require.ErrorIs(t, err, ErrEmptyEnvelope, "expected ErrEmptyEnvelope, got %v", err)
}
