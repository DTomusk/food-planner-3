package recipe

import "errors"

var (
	ErrEmptyName           = errors.New("recipe name cannot be empty")
	ErrNoIngredients       = errors.New("recipe must have at least one ingredient")
	ErrInvalidQuantity     = errors.New("ingredient quantity must be greater than zero")
	ErrIngredientNotFound  = errors.New("ingredient not found")
	ErrDuplicateIngredient = errors.New("duplicate ingredient in recipe")
	ErrInvalidUnit         = errors.New("invalid unit for ingredient usage")
	ErrInvalidPrepMins     = errors.New("prep minutes cannot be negative")
	ErrInvalidCookMins     = errors.New("cook minutes cannot be negative")
	ErrInvalidPortions     = errors.New("portions must be greater than zero")
	ErrEmptyURL            = errors.New("URL cannot be empty")
	ErrEmptyBookTitle      = errors.New("book title cannot be empty")
	ErrInvalidBookPage     = errors.New("book page must be greater than zero")
	ErrEmptyInstructions   = errors.New("instructions cannot be empty")
	ErrNoSource            = errors.New("recipe source is required")
	ErrInvalidSourceType   = errors.New("invalid recipe source type")
)
