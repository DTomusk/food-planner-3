package recipe

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Ingredients []*IngredientUsage
	PrepMins    int
	CookMins    int
	Portions    int
	Source      *RecipeSource
	DeletedOn   *time.Time
}

func NewRecipe(name string, userID uuid.UUID, ingredients []*IngredientUsage, prepMins, cookMins, portions int, source *RecipeSource) (*Recipe, error) {
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
	return &Recipe{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Ingredients: ingredients,
		PrepMins:    prepMins,
		CookMins:    cookMins,
		Portions:    portions,
		Source:      source,
	}, nil
}

func (r *Recipe) String() string {
	return r.Name
}
