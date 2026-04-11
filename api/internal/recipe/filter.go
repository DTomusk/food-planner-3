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
	ContainsGluten     *bool
}

// normalizedRecipeFilter holds a validated, scrubbed form of RecipeFilter.
// Always construct via normalizeFilter() to guarantee invariants.
type normalizedRecipeFilter struct {
	UserID             *uuid.UUID
	AnimalProductLevel *int
	ContainsGluten     *bool
}

func normalizeFilter(f RecipeFilter) normalizedRecipeFilter {
	return normalizedRecipeFilter{
		UserID:             f.UserID,
		AnimalProductLevel: normalizedAnimalProductLevelFilter(f.AnimalProductLevel),
		ContainsGluten:     f.ContainsGluten,
	}
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

func filterHash(mode RecipeCursorMode, query *string, nf normalizedRecipeFilter) string {
	q := ""
	if query != nil {
		q = *query
	}
	u := ""
	if nf.UserID != nil {
		u = nf.UserID.String()
	}
	a := ""
	if nf.AnimalProductLevel != nil {
		a = fmt.Sprintf("%d", *nf.AnimalProductLevel)
	}
	g := ""
	if nf.ContainsGluten != nil {
		g = fmt.Sprintf("%t", *nf.ContainsGluten)
	}
	h := sha256.Sum256([]byte(fmt.Sprintf("mode=%s|q=%s|u=%s|a=%s|g=%s", mode, q, u, a, g)))
	return fmt.Sprintf("%x", h[:8]) // 8 bytes = 16 hex chars, compact but collision-safe enough
}
