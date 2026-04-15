package events

import (
	"encoding/json"
	"reflect"
)

// The serialized form of an event stored in redis (or other provider)
// this is how the event is transported
// the registry then contains the method for getting the correct struct that this unmarshals to
type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func MarshalEvent(e Event) (Envelope, error) {
	if isNilEvent(e) {
		return Envelope{}, ErrNilEvent
	}

	eventType := e.Metadata().Type
	if eventType == "" {
		return Envelope{}, ErrEmptyEventType
	}

	payload, err := json.Marshal(e)
	if err != nil {
		return Envelope{}, err
	}

	return Envelope{
		Type:    eventType,
		Payload: payload,
	}, nil
}

func MarshalEnvelope(env Envelope) ([]byte, error) {
	if env.Type == "" {
		return nil, ErrEmptyEnvelopeType
	}
	if len(env.Payload) == 0 {
		return nil, ErrEmptyEnvelopePayload
	}

	return json.Marshal(env)
}

func UnmarshalEnvelope(data []byte) (Envelope, error) {
	if len(data) == 0 {
		return Envelope{}, ErrEmptyEnvelope
	}

	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return Envelope{}, err
	}
	if env.Type == "" {
		return Envelope{}, ErrEmptyEnvelopeType
	}
	if len(env.Payload) == 0 {
		return Envelope{}, ErrEmptyEnvelopePayload
	}

	return env, nil
}

func UnmarshalEvent(env Envelope, registry *Registry) (Event, error) {
	if registry == nil {
		return nil, ErrNilRegistry
	}
	if env.Type == "" {
		return nil, ErrEmptyEnvelopeType
	}
	if len(env.Payload) == 0 {
		return nil, ErrEmptyEnvelopePayload
	}

	event, ok := registry.New(env.Type)
	if !ok {
		return nil, ErrUnknownEventType
	}

	if err := json.Unmarshal(env.Payload, event); err != nil {
		return nil, err
	}
	if isNilEvent(event) {
		return nil, ErrNilRegistryEvent
	}

	// Normalize pointer-backed decoded events to value events when the value
	// also satisfies Event. This keeps downstream handlers consistent.
	v := reflect.ValueOf(event)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		if normalized, ok := v.Elem().Interface().(Event); ok {
			return normalized, nil
		}
	}

	return event, nil
}

func isNilEvent(e Event) bool {
	if e == nil {
		return true
	}
	v := reflect.ValueOf(e)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
