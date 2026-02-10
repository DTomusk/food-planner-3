package recipe

type IngredientUsage struct {
	Ingredient Ingredient
	Quantity   float32
	Unit       Unit
}

func ValidateIngredientUsage(usage IngredientUsage) error {
	if usage.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	return nil
}
