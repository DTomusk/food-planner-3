package recipe

import (
	"github.com/google/uuid"
)

type Recipe struct {
	ID          uuid.UUID
	Name        string
	Ingredients []*IngredientUsage
}

func NewRecipe(name string, ingredients []*IngredientUsage) (*Recipe, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	return &Recipe{
		ID:          uuid.New(),
		Name:        name,
		Ingredients: ingredients,
	}, nil
}

func (r *Recipe) String() string {
	return r.Name
}
