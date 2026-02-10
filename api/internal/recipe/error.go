package recipe

import "errors"

var (
	ErrEmptyName       = errors.New("recipe name cannot be empty")
	ErrNoIngredients   = errors.New("recipe must have at least one ingredient")
	ErrInvalidQuantity = errors.New("ingredient quantity must be greater than zero")
)
