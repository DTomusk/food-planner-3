package recipe

type CreateRecipeRequest struct {
	Name        string                         `json:"name" validate:"required"`
	Description string                         `json:"description,omitempty"`
	Ingredients []CreateIngredientUsageRequest `json:"ingredients" validate:"required,dive"`
	UserID      string                         `json:"user_id" validate:"required,uuid4"`
	PrepMins    int                            `json:"prep_mins" validate:"required,gte=0"`
	CookMins    int                            `json:"cook_mins" validate:"required,gte=0"`
	Portions    int                            `json:"portions" validate:"required,gt=0"`
	Source      CreateRecipeSourceRequest      `json:"source" validate:"required,dive"`
	ImgUploadID *string                        `json:"img_upload_id,omitempty"`
	IPAddress   string                         `json:"ip_address" validate:"required,ip"`
	UserAgent   string                         `json:"user_agent" validate:"required"`
	Publish     bool                           `json:"publish"`
}

type CreateIngredientUsageRequest struct {
	IngredientID string  `json:"ingredient_id" validate:"required,uuid4"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	Unit         int     `json:"unit" validate:"required"`
}

type CreateRecipeSourceRequest struct {
	Type         int     `json:"type" validate:"required"`
	URL          *string `json:"url,omitempty"`
	BookTitle    *string `json:"book_title,omitempty"`
	BookPage     *int32  `json:"book_page,omitempty"`
	Instructions *string `json:"instructions,omitempty"`
}

type UpdateRecipeRequest struct {
	RecipeId    string              `json:"recipe_id" validate:"required,uuid4"`
	Request     CreateRecipeRequest `json:"request" validate:"required,dive"`
	RemoveImage *bool               `json:"remove_image,omitempty"`
}

type RecipeListParams struct {
	Pagination    RecipePagination
	Filter        RecipeFilter
	IncludeDrafts bool
}

type RecipePagination struct {
	First int
	After *string
}

type RecipeWithCursor struct {
	Recipe *RecipeContainer
	Cursor string
}

type RecipeListRow struct {
	Recipe         *RecipeContainer
	RelevanceScore *float64
}
