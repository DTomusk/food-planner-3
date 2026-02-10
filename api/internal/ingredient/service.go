package ingredient

import (
	"context"
	"foodplanner/internal/db"
)

type IngredientService struct {
	repo *IngredientRepo
}

func NewIngredientService(repo *IngredientRepo) *IngredientService {
	return &IngredientService{
		repo: repo,
	}
}

func (s *IngredientService) Exists(ctx context.Context, db db.DBTX, ingredientID string) (bool, error) {
	return s.repo.IngredientExists(ctx, db, ingredientID)
}
