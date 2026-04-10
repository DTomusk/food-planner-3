package recipe

import (
	"time"

	"github.com/google/uuid"
)

type RecipeContainer struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	CurrentVersionID uuid.UUID
	CurrentVersion   *RecipeVersion

	CreatedAt time.Time
	DeletedAt *time.Time
}

type RecipeVersion struct {
	ID       uuid.UUID
	RecipeID uuid.UUID
	Version  int

	Name        string
	Description string
	Ingredients []*IngredientUsage
	PrepMins    int
	CookMins    int
	Portions    int
	Source      *RecipeSource
	ImgSrc      *string

	CreatedAt time.Time
}

func NewRecipe(
	name string,
	description string,
	userID uuid.UUID,
	ingredients []*IngredientUsage,
	prepMins,
	cookMins,
	portions int,
	source *RecipeSource,
	imgSrc *string,
) (*RecipeContainer, error) {
	recipeID := uuid.New()
	now := time.Now()

	version, err := NewRecipeVersion(recipeID, 1, name, description, ingredients, prepMins, cookMins, portions, source, imgSrc)
	if err != nil {
		return nil, err
	}

	recipe := &RecipeContainer{
		ID:               recipeID,
		UserID:           userID,
		CurrentVersionID: version.ID,
		CurrentVersion:   version,
		CreatedAt:        now,
	}

	return recipe, nil
}

func NewRecipeVersion(
	recipeID uuid.UUID,
	version int,
	name string,
	description string,
	ingredients []*IngredientUsage,
	prepMins,
	cookMins,
	portions int,
	source *RecipeSource,
	imgSrc *string,
) (*RecipeVersion, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	if len(name) > 100 {
		return nil, ErrNameTooLong
	}

	if len(description) > 200 {
		return nil, ErrInvalidDescription
	}

	if len(ingredients) == 0 {
		return nil, ErrNoIngredients
	}

	if prepMins < 0 {
		return nil, ErrInvalidPrepMins
	}
	if cookMins < 0 {
		return nil, ErrInvalidCookMins
	}
	if portions <= 0 {
		return nil, ErrInvalidPortions
	}

	now := time.Now()
	return &RecipeVersion{
		ID:          uuid.New(),
		RecipeID:    recipeID,
		Version:     version,
		Name:        name,
		Description: description,
		Ingredients: ingredients,
		PrepMins:    prepMins,
		CookMins:    cookMins,
		Portions:    portions,
		Source:      source,
		CreatedAt:   now,
		ImgSrc:      imgSrc,
	}, nil
}

func (r *RecipeVersion) String() string {
	return r.Name
}
