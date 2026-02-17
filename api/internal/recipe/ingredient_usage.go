package recipe

import (
	"foodplanner/internal/unit"

	"github.com/google/uuid"
)

type IngredientUsage struct {
	ID           uuid.UUID
	IngredientID uuid.UUID
	Quantity     float32
	Unit         unit.Unit
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
