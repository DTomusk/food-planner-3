package sync

import (
	"context"
	"fmt"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/logging"
	"foodplanner/internal/reference"

	"github.com/google/uuid"
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

	existingIngredients, err := s.ingredientService.GetAllIngredients(ctx, logger)
	if err != nil {
		return err
	}

	existingIDsByFileKey := make(map[string]uuid.UUID, len(existingIngredients))
	for _, existingIngredient := range existingIngredients {
		existingIDsByFileKey[existingIngredient.FileKey] = existingIngredient.ID
	}

	resolvedIDsByFileKey := make(map[string]uuid.UUID, len(fileIngredients))
	for _, fileIngredient := range fileIngredients {
		if _, exists := resolvedIDsByFileKey[fileIngredient.ID]; exists {
			return fmt.Errorf("duplicate ingredient file key %q", fileIngredient.ID)
		}

		resolvedID, exists := existingIDsByFileKey[fileIngredient.ID]
		if !exists {
			resolvedID = uuid.New()
		}
		resolvedIDsByFileKey[fileIngredient.ID] = resolvedID
	}

	domainIngredients := make([]*ingredient.Ingredient, len(fileIngredients))

	for i, fileIngredient := range fileIngredients {
		processedLevel := ingredient.Raw
		if fileIngredient.ProcessedLevel != nil {
			processedLevel = ingredient.ProcessedLevel(*fileIngredient.ProcessedLevel)
		}

		isSearchable := true
		if fileIngredient.IsSearchable != nil {
			isSearchable = *fileIngredient.IsSearchable
		}

		var taxonomyParentID *uuid.UUID
		if fileIngredient.TaxonomyParentKey != nil && *fileIngredient.TaxonomyParentKey != "" {
			resolvedParentID, exists := resolvedIDsByFileKey[*fileIngredient.TaxonomyParentKey]
			if !exists {
				return fmt.Errorf("ingredient %q references unknown taxonomy parent key %q", fileIngredient.ID, *fileIngredient.TaxonomyParentKey)
			}
			taxonomyParentID = &resolvedParentID
		}

		domainIngredients[i], err = ingredient.NewIngredient(
			fileIngredient.Name,
			fileIngredient.ID,
			fileIngredient.PreferredUnit,
			fileIngredient.Counter,
			fileIngredient.Plural,
			fileIngredient.CounterPlural,
			ingredient.AnimalProductLevel(fileIngredient.AnimalProductLevel),
			fileIngredient.ContainsGluten,
			taxonomyParentID,
			processedLevel,
			isSearchable,
		)
		if err != nil {
			return err
		}
		domainIngredients[i].ID = resolvedIDsByFileKey[fileIngredient.ID]
	}

	if err := s.ingredientService.SyncIngredientData(ctx, logger, domainIngredients); err != nil {
		return err
	}
	logger.Info("Successfully synchronized ingredient data")
	return nil
}
