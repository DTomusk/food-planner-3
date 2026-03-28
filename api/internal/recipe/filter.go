package recipe

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
)

type RecipeFilter struct {
	Query  *string
	UserID *uuid.UUID
}

func filterHashForParams(mode RecipeCursorMode, query *string, userID *uuid.UUID) string {
	q := ""
	if query != nil {
		q = *query
	}
	u := ""
	if userID != nil {
		u = userID.String()
	}
	h := sha256.Sum256([]byte(fmt.Sprintf("mode=%s|q=%s|u=%s", mode, q, u)))
	return fmt.Sprintf("%x", h[:8]) // 8 bytes = 16 hex chars, compact but collision-safe enough
}
