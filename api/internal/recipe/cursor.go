package recipe

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

const recipeCursorVersion = 1

type RecipeCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

type recipeCursorPayload struct {
	Version   int    `json:"v"`
	CreatedAt string `json:"createdAt"`
	ID        string `json:"id"`
}

func ParseRecipeCursor(cursor *string) (*RecipeCursor, error) {
	if cursor == nil {
		return nil, nil
	}

	rawCursor := strings.TrimSpace(*cursor)
	if rawCursor == "" {
		return nil, nil
	}

	decoded, err := base64.RawURLEncoding.DecodeString(rawCursor)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	var payload recipeCursorPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, ErrInvalidCursor
	}

	if payload.Version != recipeCursorVersion {
		return nil, ErrInvalidCursor
	}

	createdAt, err := time.Parse(time.RFC3339Nano, payload.CreatedAt)
	if err != nil {
		return nil, ErrInvalidCursor
	}
	if createdAt.IsZero() {
		return nil, ErrInvalidCursor
	}

	id, err := uuid.Parse(payload.ID)
	if err != nil {
		return nil, ErrInvalidCursor
	}
	if id == uuid.Nil {
		return nil, ErrInvalidCursor
	}

	return &RecipeCursor{
		CreatedAt: createdAt,
		ID:        id,
	}, nil
}

func (c *RecipeCursor) String() string {
	if c == nil {
		return ""
	}
	if c.CreatedAt.IsZero() || c.ID == uuid.Nil {
		return ""
	}

	payload := recipeCursorPayload{
		Version:   recipeCursorVersion,
		CreatedAt: c.CreatedAt.UTC().Format(time.RFC3339Nano),
		ID:        c.ID.String(),
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(encoded)
}
