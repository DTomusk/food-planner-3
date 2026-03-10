package recipe

import (
	"time"

	"github.com/google/uuid"
)

type RecipeContainer struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	CurrentVersionID uuid.UUID
	CurrentVersion   *RecipeVersion

	CreatedAt time.Time
	DeletedAt *time.Time
}

type RecipeVersion struct {
	ID       uuid.UUID
	RecipeID uuid.UUID
	Version  int

	Name        string
	Ingredients []*IngredientUsage
	PrepMins    int
	CookMins    int
	Portions    int
	Source      *RecipeSource

	CreatedAt time.Time
}

func NewRecipe(
	name string,
	userID uuid.UUID,
	ingredients []*IngredientUsage,
	prepMins,
	cookMins,
	portions int,
	source *RecipeSource,
) (*RecipeContainer, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	if len(ingredients) == 0 {
		return nil, ErrNoIngredients
	}
	if prepMins < 0 {
		return nil, ErrInvalidPrepMins
	}
	if cookMins < 0 {
		return nil, ErrInvalidCookMins
	}
	if portions <= 0 {
		return nil, ErrInvalidPortions
	}

	recipeID := uuid.New()
	versionID := uuid.New()
	now := time.Now()

	version := &RecipeVersion{
		ID:          versionID,
		RecipeID:    recipeID,
		Version:     1,
		Name:        name,
		Ingredients: ingredients,
		PrepMins:    prepMins,
		CookMins:    cookMins,
		Portions:    portions,
		Source:      source,
		CreatedAt:   now,
	}

	recipe := &RecipeContainer{
		ID:               recipeID,
		UserID:           userID,
		CurrentVersionID: versionID,
		CurrentVersion:   version,
		CreatedAt:        now,
	}

	return recipe, nil
}

func (r *RecipeVersion) String() string {
	return r.Name
}
