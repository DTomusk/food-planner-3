package recipe

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

const recipeCursorVersion = 2

type RecipeCursorMode string

const (
	RecipeCursorModeNewest    RecipeCursorMode = "newest"
	RecipeCursorModeRelevance RecipeCursorMode = "relevance"
)

type RecipeCursor struct {
	// Mode indicates what order is used in the sort and hence what cursor values to use.
	Mode RecipeCursorMode
	// FilterHash binds the cursor to a normalized filter/sort state.
	FilterHash string
	CreatedAt  time.Time
	ID         uuid.UUID
	// RelevanceScore stores the last row's score for relevance-ordered pagination.
	RelevanceScore *float64
}

type recipeCursorPayload struct {
	Version        int      `json:"v"`
	Mode           string   `json:"mode,omitempty"`
	FilterHash     string   `json:"filterHash,omitempty"`
	CreatedAt      string   `json:"createdAt"`
	ID             string   `json:"id"`
	RelevanceScore *float64 `json:"relevanceScore,omitempty"`
}

func normalizeRecipeCursorMode(mode RecipeCursorMode) RecipeCursorMode {
	if mode == "" {
		return RecipeCursorModeNewest
	}
	return mode
}

// Turns an encoded cursor string into a structured cursor object
func ParseRecipeCursor(cursor *string) (*RecipeCursor, error) {
	if cursor == nil {
		return nil, nil
	}

	rawCursor := strings.TrimSpace(*cursor)
	if rawCursor == "" {
		return nil, nil
	}

	// Cursor is base64 encoded, so needs to be decoded accordingly
	decoded, err := base64.RawURLEncoding.DecodeString(rawCursor)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	// Once decoded, we can unmarshal the json string into an object
	var payload recipeCursorPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, ErrInvalidCursor
	}

	// The version must be valid and supported
	if payload.Version != recipeCursorVersion && payload.Version != 1 {
		return nil, ErrInvalidCursor
	}

	// Since we didn't have a mode before, an empty mode is treated as the default newest mode
	mode := normalizeRecipeCursorMode(RecipeCursorMode(payload.Mode))
	if mode != RecipeCursorModeNewest && mode != RecipeCursorModeRelevance {
		return nil, ErrInvalidCursor
	}

	// CreatedAt and ID are required for all cursor modes, so we validate them regardless of mode
	createdAt, err := time.Parse(time.RFC3339Nano, payload.CreatedAt)
	if err != nil {
		return nil, ErrInvalidCursor
	}
	if createdAt.IsZero() {
		return nil, ErrInvalidCursor
	}

	id, err := uuid.Parse(payload.ID)
	if err != nil || id == uuid.Nil {
		return nil, ErrInvalidCursor
	}

	// Check relevance score if cursor is relevance based
	if mode == RecipeCursorModeRelevance {
		if payload.RelevanceScore == nil || math.IsNaN(*payload.RelevanceScore) || math.IsInf(*payload.RelevanceScore, 0) {
			return nil, ErrInvalidCursor
		}
	}

	return &RecipeCursor{
		Mode:           mode,
		FilterHash:     payload.FilterHash,
		CreatedAt:      createdAt,
		ID:             id,
		RelevanceScore: payload.RelevanceScore,
	}, nil
}

// Encodes a cursor into a base64 encoded string
func (c *RecipeCursor) String() string {
	if c == nil {
		return ""
	}
	if c.CreatedAt.IsZero() || c.ID == uuid.Nil {
		return ""
	}

	mode := normalizeRecipeCursorMode(c.Mode)
	if mode != RecipeCursorModeNewest && mode != RecipeCursorModeRelevance {
		return ""
	}

	payload := recipeCursorPayload{
		Version:    recipeCursorVersion,
		Mode:       string(mode),
		FilterHash: c.FilterHash,
		CreatedAt:  c.CreatedAt.UTC().Format(time.RFC3339Nano),
		ID:         c.ID.String(),
	}

	if mode == RecipeCursorModeRelevance {
		if c.RelevanceScore == nil || math.IsNaN(*c.RelevanceScore) || math.IsInf(*c.RelevanceScore, 0) {
			return ""
		}
		payload.RelevanceScore = c.RelevanceScore
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(encoded)
}

// Check that the receiving cursor matches the mode and hash passed in
func (c *RecipeCursor) Matches(mode RecipeCursorMode, filterHash string) bool {
	if c == nil {
		return true
	}

	normalizedCursorMode := normalizeRecipeCursorMode(c.Mode)
	normalizedMode := normalizeRecipeCursorMode(mode)
	return normalizedCursorMode == normalizedMode && c.FilterHash == filterHash
}
