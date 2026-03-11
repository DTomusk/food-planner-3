package recipe

import (
	"foodplanner/internal/ingredient"
	"foodplanner/internal/unit"

	"github.com/google/uuid"
)

type IngredientUsage struct {
	ID           uuid.UUID
	IngredientID uuid.UUID
	Quantity     float64
	Unit         unit.Unit
}

func newIngredientUsages(requests []CreateIngredientUsageRequest, availableIngredients []*ingredient.Ingredient) ([]*IngredientUsage, error) {
	ingredientMap := make(map[string]*ingredient.Ingredient, len(availableIngredients))
	for _, ing := range availableIngredients {
		ingredientMap[ing.ID.String()] = ing
	}

	seenIngredients := make(map[string]bool)
	ingredientUsages := make([]*IngredientUsage, len(requests))

	for i, ingredientRequest := range requests {
		if seenIngredients[ingredientRequest.IngredientID] {
			return nil, ErrDuplicateIngredient
		}
		seenIngredients[ingredientRequest.IngredientID] = true

		selectedIngredient := ingredientMap[ingredientRequest.IngredientID]
		if ingredientRequest.Unit != int(selectedIngredient.PreferredUnit) {
			return nil, ErrInvalidUnit
		}

		usage, err := NewIngredientUsage(ingredientRequest)
		if err != nil {
			return nil, err
		}
		ingredientUsages[i] = usage
	}
	return ingredientUsages, nil
}

func NewIngredientUsage(request CreateIngredientUsageRequest) (*IngredientUsage, error) {
	if request.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if !unit.Unit(request.Unit).IsValid() {
		return nil, ErrInvalidUnit
	}
	return &IngredientUsage{
		ID:           uuid.New(),
		IngredientID: uuid.MustParse(request.IngredientID),
		Quantity:     request.Quantity,
		Unit:         unit.Unit(request.Unit),
	}, nil
}
