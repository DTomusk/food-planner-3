package sync

import (
	"context"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/logging"
	"foodplanner/internal/reference"
)

type SyncService struct {
	ingredientService *ingredient.IngredientService
	dataLoader        *reference.Loader
}

func NewSyncService(ingredientService *ingredient.IngredientService, dataLoader *reference.Loader) *SyncService {
	return &SyncService{
		ingredientService: ingredientService,
		dataLoader:        dataLoader,
	}
}

func (s *SyncService) SyncIngredientData(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Info("Starting ingredient data synchronization")

	fileIngredients, err := s.dataLoader.LoadIngredientData(logger)
	if err != nil {
		return err
	}

	if len(fileIngredients) == 0 {
		logger.Warn("No ingredient data found in reference file")
		return nil
	}

	logger.Info("Retrieved reference ingredient data")

	domainIngredients := make([]*ingredient.Ingredient, len(fileIngredients))

	for i, fileIngredient := range fileIngredients {
		domainIngredients[i], err = ingredient.NewIngredient(
			fileIngredient.Name,
			fileIngredient.ID,
			fileIngredient.PreferredUnit,
			fileIngredient.Counter,
		)
		if err != nil {
			return err
		}
	}

	if err := s.ingredientService.SyncIngredientData(ctx, logger, domainIngredients); err != nil {
		return err
	}
	logger.Info("Successfully synchronized ingredient data")
	return nil
}
