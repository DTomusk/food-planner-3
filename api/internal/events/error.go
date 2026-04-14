package events

import "errors"

var (
	ErrNilEvent             = errors.New("event is nil")
	ErrEmptyEventType       = errors.New("event metadata type is empty")
	ErrEmptyEnvelope        = errors.New("envelope data is empty")
	ErrEmptyEnvelopeType    = errors.New("envelope type is empty")
	ErrEmptyEnvelopePayload = errors.New("envelope payload is empty")
	ErrNilRegistry          = errors.New("event registry is nil")
	ErrUnknownEventType     = errors.New("unknown event type")
	ErrNilRegistryEvent     = errors.New("event registry factory returned nil event")
)
