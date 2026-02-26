package ingredient

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"
	"log/slog"
)

type IngredientService struct {
	txRunner        db.TxRunner
	repo            *IngredientRepo
	upsertBatchSize int
}

func NewIngredientService(txRunner db.TxRunner, repo *IngredientRepo, upsertBatchSize int) *IngredientService {
	return &IngredientService{
		txRunner:        txRunner,
		repo:            repo,
		upsertBatchSize: upsertBatchSize,
	}
}

func (s *IngredientService) Exists(ctx context.Context, ingredientID string) (bool, error) {
	return s.repo.IngredientExists(ctx, s.txRunner.DB(), ingredientID)
}

func (s *IngredientService) SyncIngredientData(ctx context.Context, logger *slog.Logger, ingredients []*Ingredient) error {
	logger.Info("Starting ingredient data sync", "numIngredients", len(ingredients), "batchSize", s.upsertBatchSize)
	err := s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		numberOfBatches := (len(ingredients) + s.upsertBatchSize - 1) / s.upsertBatchSize
		for i := 0; i < numberOfBatches; i++ {
			logger.Info("Upserting ingredient batch", "batchNumber", i+1, "totalBatches", numberOfBatches)
			start := i * s.upsertBatchSize
			end := start + s.upsertBatchSize
			if end > len(ingredients) {
				end = len(ingredients)
			}
			batch := ingredients[start:end]
			if err := s.repo.UpsertIngredients(ctx, tx, batch); err != nil {
				logger.Error("Failed to upsert ingredient batch", "batchNumber", i+1, "totalBatches", numberOfBatches, "error", err)
				return err
			}
		}
		return nil
	})
	return err
}

// TODO: remove and replace with paginated (and maybe filtered) search
func (s *IngredientService) GetAllIngredients(ctx context.Context, logger *slog.Logger) ([]*Ingredient, error) {
	ingredients, err := s.repo.GetAllIngredients(ctx, s.txRunner.DB())
	if err != nil {
		logger.Error("Failed to fetch ingredients", "error", err)
		return nil, err
	}
	return ingredients, nil
}

func (s *IngredientService) GetIngredientsByIDs(ctx context.Context, logger *slog.Logger, ids []string) ([]*Ingredient, error) {
	ingredients, err := s.repo.GetIngredientsByIDs(ctx, s.txRunner.DB(), ids)
	if err != nil {
		logger.Error("Failed to fetch ingredients by IDs", "error", err)
		return nil, err
	}
	return ingredients, nil
}
