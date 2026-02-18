package recipe

type CreateRecipeRequest struct {
	Name        string                         `json:"name" validate:"required"`
	Ingredients []CreateIngredientUsageRequest `json:"ingredients" validate:"required,dive"`
	UserID      string                         `json:"user_id" validate:"required,uuid4"`
}

type CreateIngredientUsageRequest struct {
	IngredientID string  `json:"ingredient_id" validate:"required,uuid4"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	Unit         int     `json:"unit" validate:"required"`
}
