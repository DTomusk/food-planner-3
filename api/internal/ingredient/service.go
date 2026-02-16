package ingredient

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"
)

type IngredientService struct {
	txRunner        db.TxRunner
	repo            *IngredientRepo
	upsertBatchSize int
}

func NewIngredientService(txRunner db.TxRunner, repo *IngredientRepo) *IngredientService {
	return &IngredientService{
		txRunner: txRunner,
		repo:     repo,
	}
}

func (s *IngredientService) Exists(ctx context.Context, ingredientID string) (bool, error) {
	return s.repo.IngredientExists(ctx, s.txRunner.DB(), ingredientID)
}

func (s *IngredientService) SyncIngredientData(ctx context.Context, ingredients []*Ingredient) error {
	// Separate into batches then run all batches in tx
	err := s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		return s.repo.UpsertIngredients(ctx, tx, ingredients)
	})
	return err
}

// TODO: remove and replace with paginated (and maybe filtered) search
func (s *IngredientService) GetAllIngredients(ctx context.Context) ([]*Ingredient, error) {
	return s.repo.GetAllIngredients(ctx, s.txRunner.DB())
}
