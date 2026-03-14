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

// ingredient_usage.go
func newIngredientUsages(
	requests []CreateIngredientUsageRequest,
	ingredientByID map[string]*ingredient.Ingredient,
) ([]*IngredientUsage, error) {
	usages := make([]*IngredientUsage, len(requests))
	for i, req := range requests {
		ing := ingredientByID[req.IngredientID]
		if ing == nil {
			return nil, ErrIngredientNotFound
		}
		if req.Unit != int(ing.PreferredUnit) {
			return nil, ErrInvalidUnit
		}

		usage, err := NewIngredientUsage(req)
		if err != nil {
			return nil, err
		}
		usages[i] = usage
	}
	return usages, nil
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
