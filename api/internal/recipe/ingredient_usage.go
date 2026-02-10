package recipe

import "github.com/google/uuid"

type IngredientUsage struct {
	Ingredient Ingredient
	Quantity   float32
	Unit       Unit
}

func NewIngredientUsage(request CreateIngredientUsageRequest) (*IngredientUsage, error) {
	if request.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	return &IngredientUsage{
		Ingredient: Ingredient{ID: uuid.MustParse(request.IngredientID)},
		Quantity:   request.Quantity,
		Unit:       Unit(request.Unit),
	}, nil
}
