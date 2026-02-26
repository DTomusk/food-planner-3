package recipe

type CreateRecipeRequest struct {
	Name        string                         `json:"name" validate:"required"`
	Ingredients []CreateIngredientUsageRequest `json:"ingredients" validate:"required,dive"`
	UserID      string                         `json:"user_id" validate:"required,uuid4"`
	PrepMins    int                            `json:"prep_mins" validate:"required,gte=0"`
	CookMins    int                            `json:"cook_mins" validate:"required,gte=0"`
	Portions    int                            `json:"portions" validate:"required,gt=0"`
	Source      CreateRecipeSourceRequest      `json:"source" validate:"required,dive"`
}

type CreateIngredientUsageRequest struct {
	IngredientID string  `json:"ingredient_id" validate:"required,uuid4"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	Unit         int     `json:"unit" validate:"required"`
}

type CreateRecipeSourceRequest struct {
	Type         SourceType `json:"type" validate:"required,oneof=url book_reference original"`
	URL          string     `json:"url,omitempty"`
	BookTitle    string     `json:"book_title,omitempty"`
	BookPage     int        `json:"book_page,omitempty"`
	Instructions string     `json:"instructions,omitempty"`
}
