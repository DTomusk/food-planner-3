package sync

import (
	"context"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/logging"
	"foodplanner/internal/reference"
	"log/slog"

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

	// Load ingredient trees from the reference file
	fileIngredients, err := s.dataLoader.LoadIngredientData(logger)
	if err != nil {
		logger.Error("Failed to load ingredient data", "error", err)
		return err
	}

	if len(fileIngredients) == 0 {
		logger.Warn("No ingredient data found in reference file")
		return nil
	}

	logger.Info("Retrieved reference ingredient data")

	// Get all existing ingredients (including non-searchable taxonomy parents) so
	// we can preserve stable IDs across sync runs for every file key.
	existingIngredients, err := s.ingredientService.GetAllIngredientsUnfiltered(ctx, logger)
	if err != nil {
		logger.Error("Failed to fetch existing ingredients", "error", err)
		return err
	}

	// Map file keys to existing ids
	// We can check this map to see if we need to change an id during the walk
	existingIDsByFileKey := make(map[string]uuid.UUID, len(existingIngredients))
	for _, existingIngredient := range existingIngredients {
		existingIDsByFileKey[existingIngredient.FileKey] = existingIngredient.ID
	}

	// We don't know how many ingredients there are yet because of the tree structure
	domainIngredients := make([]*ingredient.Ingredient, 0, len(fileIngredients))

	// Flatten trees into one list
	// Ensure parents come before children by appending then walking
	for _, fileIngredient := range fileIngredients {
		logger.Debug("Walking ingredient tree", "fileKey", fileIngredient.FileKey, "name", fileIngredient.Name)
		// Children always have a parent, so pass uuid instead of *uuid
		domainIngredients, err = s.walkIngredientTree(fileIngredient, domainIngredients, nil, existingIDsByFileKey, logger)
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

func (s *SyncService) walkIngredientTree(
	fileIngredient *reference.IngredientModel,
	domainIngredients []*ingredient.Ingredient,
	parentID *uuid.UUID,
	keyIDMap map[string]uuid.UUID,
	logger *slog.Logger,
) ([]*ingredient.Ingredient, error) {
	domainIngredient, err := ingredient.NewIngredient(
		fileIngredient.Name,
		fileIngredient.FileKey,
		fileIngredient.PreferredUnit,
		fileIngredient.Counter,
		fileIngredient.Plural,
		fileIngredient.CounterPlural,
		ingredient.AnimalProductLevel(fileIngredient.AnimalProductLevel),
		fileIngredient.ContainsGluten,
		parentID,
		fileIngredient.ProcessedLevel,
		fileIngredient.IsSearchable,
	)
	if err != nil {
		logger.Error("Failed to create domain ingredient", "error", err)
		return nil, err
	}

	existingID, exists := keyIDMap[fileIngredient.FileKey]
	if exists {
		domainIngredient.ID = existingID
	}

	domainIngredients = append(domainIngredients, domainIngredient)

	for _, childFileIngredient := range fileIngredient.Children {
		logger.Debug("Walking child ingredient", "fileKey", childFileIngredient.FileKey, "name", childFileIngredient.Name)
		var err error
		domainIngredients, err = s.walkIngredientTree(childFileIngredient, domainIngredients, &domainIngredient.ID, keyIDMap, logger)
		if err != nil {
			return nil, err
		}
	}

	return domainIngredients, nil
}
