package ingredient

import (
	"context"
	"foodplanner/internal/db"
)

type IngredientService struct {
	db   db.DBTX
	repo *IngredientRepo
}

func NewIngredientService(db db.DBTX, repo *IngredientRepo) *IngredientService {
	return &IngredientService{
		db:   db,
		repo: repo,
	}
}

func (s *IngredientService) Exists(ctx context.Context, ingredientID string) (bool, error) {
	return s.repo.IngredientExists(ctx, s.db, ingredientID)
}

func (s *IngredientService) SyncIngredientData(ctx context.Context, ingredients []*Ingredient) error {
	// Validate ingredients and upsert
	return nil
}
