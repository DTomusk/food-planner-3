package sync

import (
	"context"
	"foodplanner/internal/ingredient"
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
	// logger := logging.FromContext(ctx)
	// logger.Info("Starting ingredient data synchronization")

	// fileIngredients, err := s.dataLoader.LoadIngredientData(logger)
	// if err != nil {
	// 	logger.Error("Failed to load ingredient data", "error", err)
	// 	return err
	// }

	// if len(fileIngredients) == 0 {
	// 	logger.Warn("No ingredient data found in reference file")
	// 	return nil
	// }

	// logger.Info("Retrieved reference ingredient data")

	// existingIngredients, err := s.ingredientService.GetAllIngredients(ctx, logger)
	// if err != nil {
	// 	logger.Error("Failed to fetch existing ingredients", "error", err)
	// 	return err
	// }

	// // Map file key -> ingredient entity ID for ingredients that are already persisted
	// existingIDsByFileKey := make(map[string]uuid.UUID, len(existingIngredients))
	// for _, existingIngredient := range existingIngredients {
	// 	existingIDsByFileKey[existingIngredient.FileKey] = existingIngredient.ID
	// }

	// // Extend the map so that all file keys are present, with new UUIDs generated for those that don't exist yet
	// resolvedIDsByFileKey := make(map[string]uuid.UUID, len(fileIngredients))
	// for _, fileIngredient := range fileIngredients {
	// 	if _, exists := resolvedIDsByFileKey[fileIngredient.FileKey]; exists {
	// 		return fmt.Errorf("duplicate ingredient file key %q", fileIngredient.FileKey)
	// 	}

	// 	resolvedID, exists := existingIDsByFileKey[fileIngredient.FileKey]
	// 	if !exists {
	// 		resolvedID = uuid.New()
	// 	}
	// 	resolvedIDsByFileKey[fileIngredient.FileKey] = resolvedID
	// }

	// domainIngredients := make([]*ingredient.Ingredient, len(fileIngredients))

	// for i, fileIngredient := range fileIngredients {
	// 	// Save the optional taxonomy parent ID as a reference to the resolved ID of the parent ingredient
	// 	var taxonomyParentID *uuid.UUID
	// 	if fileIngredient.TaxonomyParentKey != nil && *fileIngredient.TaxonomyParentKey != "" {
	// 		resolvedParentID, exists := resolvedIDsByFileKey[*fileIngredient.TaxonomyParentKey]
	// 		// This is an error in the reference data, taxonomy parents have to be defined
	// 		if !exists {
	// 			return fmt.Errorf("ingredient %q references unknown taxonomy parent key %q", fileIngredient.FileKey, *fileIngredient.TaxonomyParentKey)
	// 		}
	// 		taxonomyParentID = &resolvedParentID
	// 	}

	// 	domainIngredients[i], err = ingredient.NewIngredient(
	// 		fileIngredient.Name,
	// 		fileIngredient.FileKey,
	// 		fileIngredient.PreferredUnit,
	// 		fileIngredient.Counter,
	// 		fileIngredient.Plural,
	// 		fileIngredient.CounterPlural,
	// 		ingredient.AnimalProductLevel(fileIngredient.AnimalProductLevel),
	// 		fileIngredient.ContainsGluten,
	// 		taxonomyParentID,
	// 		fileIngredient.ProcessedLevel,
	// 		fileIngredient.IsSearchable,
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Crucially, replace the randomly generated id with the actual id of the ingredient
	// 	domainIngredients[i].ID = resolvedIDsByFileKey[fileIngredient.FileKey]
	// }

	// if err := s.ingredientService.SyncIngredientData(ctx, logger, domainIngredients); err != nil {
	// 	return err
	// }
	// logger.Info("Successfully synchronized ingredient data")
	return nil
}
