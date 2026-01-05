package recipe

import (
	"errors"

	"github.com/google/uuid"
)

var ErrEmptyName = errors.New("recipe name cannot be empty")

type Recipe struct {
	ID   uuid.UUID
	Name string
}

func NewRecipe(name string) (Recipe, error) {
	if name == "" {
		return Recipe{}, ErrEmptyName
	}
	return Recipe{
		ID:   uuid.New(),
		Name: name,
	}, nil
}

func (r *Recipe) String() string {
	return r.Name
}
