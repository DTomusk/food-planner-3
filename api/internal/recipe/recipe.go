package recipe

import (
	"time"

	"github.com/google/uuid"
)

type RecipeContainer struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	CurrentVersionID uuid.UUID

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

func NewRecipe(name string, userID uuid.UUID, ingredients []*IngredientUsage, prepMins, cookMins, portions int, source *RecipeSource) (*RecipeContainer, *RecipeVersion, error) {
	if name == "" {
		return nil, nil, ErrEmptyName
	}
	if len(ingredients) == 0 {
		return nil, nil, ErrNoIngredients
	}
	if prepMins < 0 {
		return nil, nil, ErrInvalidPrepMins
	}
	if cookMins < 0 {
		return nil, nil, ErrInvalidCookMins
	}
	if portions <= 0 {
		return nil, nil, ErrInvalidPortions
	}

	recipeID := uuid.New()
	versionID := uuid.New()
	now := time.Now()

	recipe := &RecipeContainer{
		ID:               recipeID,
		UserID:           userID,
		CurrentVersionID: versionID,
		CreatedAt:        now,
	}

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
	return recipe, version, nil
}

func (r *RecipeVersion) String() string {
	return r.Name
}
