package recipe

import (
	"crypto/sha256"
	"fmt"
)

type RecipeFilter struct {
	Query *string
}

func filterHashForParams(mode RecipeCursorMode, query *string) string {
	q := ""
	if query != nil {
		q = *query
	}
	h := sha256.Sum256([]byte(fmt.Sprintf("mode=%s|q=%s", mode, q)))
	return fmt.Sprintf("%x", h[:8]) // 8 bytes = 16 hex chars, compact but collision-safe enough
}
