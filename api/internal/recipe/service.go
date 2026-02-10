package recipe

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/ingredient"
)

type Service struct {
	db                db.DBTX
	Repo              *Repo
	IngredientService *ingredient.IngredientService
}

func NewService(db db.DBTX, repo *Repo, ingredientService *ingredient.IngredientService) *Service {
	return &Service{
		db:                db,
		Repo:              repo,
		IngredientService: ingredientService,
	}
}

func (s *Service) CreateRecipe(ctx context.Context, request CreateRecipeRequest) (*Recipe, error) {
	if len(request.Ingredients) == 0 {
		return nil, ErrNoIngredients
	}
	seenIngredients := make(map[string]bool)
	ingredientUsages := make([]*IngredientUsage, len(request.Ingredients))
	for i, ingredient := range request.Ingredients {
		if seenIngredients[ingredient.IngredientID] {
			return nil, ErrDuplicateIngredient
		}
		seenIngredients[ingredient.IngredientID] = true
		// ensure ingredient exists (in future, grab ingredient validation rules)
		exists, err := s.IngredientService.Exists(ctx, s.db, ingredient.IngredientID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrIngredientNotFound
		}
		usage, err := NewIngredientUsage(ingredient)
		if err != nil {
			return nil, err
		}
		ingredientUsages[i] = usage
	}

	recipe, err := NewRecipe(request.Name, ingredientUsages)
	if err != nil {
		return nil, err
	}
	return s.Repo.CreateRecipe(ctx, s.db, recipe)
}

func (s *Service) GetAllRecipes(ctx context.Context) ([]*Recipe, error) {
	return s.Repo.GetAllRecipes(ctx, s.db)
}

func (s *Service) GetRecipeByID(ctx context.Context, id string) (*Recipe, error) {
	return s.Repo.GetRecipeByID(ctx, s.db, id)
}
