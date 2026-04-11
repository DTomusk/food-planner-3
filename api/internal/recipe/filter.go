package recipe

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
)

type RecipeFilter struct {
	Query              *string
	UserID             *uuid.UUID
	AnimalProductLevel *int
}

func normalizedAnimalProductLevelFilter(animalProductLevel *int) *int {
	if animalProductLevel == nil {
		return nil
	}
	if *animalProductLevel == 0 || *animalProductLevel == 1 {
		return animalProductLevel
	}
	return nil
}

func filterHashForParams(mode RecipeCursorMode, query *string, userID *uuid.UUID, animalProductLevel *int) string {
	q := ""
	if query != nil {
		q = *query
	}
	u := ""
	if userID != nil {
		u = userID.String()
	}
	a := ""
	if animalProductLevel != nil {
		a = fmt.Sprintf("%d", *animalProductLevel)
	}
	h := sha256.Sum256([]byte(fmt.Sprintf("mode=%s|q=%s|u=%s|a=%s", mode, q, u, a)))
	return fmt.Sprintf("%x", h[:8]) // 8 bytes = 16 hex chars, compact but collision-safe enough
}
