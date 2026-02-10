package recipe

import (
	"context"
	"foodplanner/internal/db"
)

type Service struct {
	db   db.DBTX
	Repo *Repo
}

func NewService(db db.DBTX, repo *Repo) *Service {
	return &Service{
		db:   db,
		Repo: repo,
	}
}

func (s *Service) CreateRecipe(ctx context.Context, request CreateRecipeRequest) (*Recipe, error) {
	if len(request.Ingredients) == 0 {
		return nil, ErrNoIngredients
	}
	ingredients, err := s.CreateIngredientUsages(ctx, request.Ingredients)
	if err != nil {
		return nil, err
	}
	entity, err := NewRecipe(request.Name, ingredients)
	if err != nil {
		return nil, err
	}
	return s.Repo.CreateRecipe(ctx, s.db, entity)
}

func (s *Service) CreateIngredientUsages(ctx context.Context, requests []CreateIngredientUsageRequest) ([]IngredientUsage, error) {
	// For each ingredient usage, validate quantity and unit, and fetch the ingredient details from the database
	var usages []IngredientUsage
	for _, req := range requests {
		// validate usage
		if req.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
}

func (s *Service) GetAllRecipes(ctx context.Context) ([]*Recipe, error) {
	return s.Repo.GetAllRecipes(ctx, s.db)
}

func (s *Service) GetRecipeByID(ctx context.Context, id string) (*Recipe, error) {
	return s.Repo.GetRecipeByID(ctx, s.db, id)
}
