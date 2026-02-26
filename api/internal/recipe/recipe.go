package recipe

import (
	"github.com/google/uuid"
)

type Recipe struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Ingredients []*IngredientUsage
}

func NewRecipe(name string, userID uuid.UUID, ingredients []*IngredientUsage) (*Recipe, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	return &Recipe{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Ingredients: ingredients,
	}, nil
}

func (r *Recipe) String() string {
	return r.Name
}
