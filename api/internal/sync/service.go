package sync

import (
	"context"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/reference"
)

type SyncService struct {
	ingredientService *ingredient.IngredientService
	referenceService  *reference.ReferenceService
}

func NewSyncService(ingredientService *ingredient.IngredientService, referenceService *reference.ReferenceService) *SyncService {
	return &SyncService{
		ingredientService: ingredientService,
		referenceService:  referenceService,
	}
}

func (s *SyncService) SyncIngredientData() error {
	fileIngredients, err := s.referenceService.LoadIngredientData()
	if err != nil {
		return err
	}

	domainIngredients := make([]*ingredient.Ingredient, len(fileIngredients))

	for i, fileIngredient := range fileIngredients {
		domainIngredients[i], err = ingredient.NewIngredient(fileIngredient.Name, fileIngredient.FileKey, fileIngredient.PreferredUnit)
		if err != nil {
			return err
		}
	}

	if err := s.ingredientService.SyncIngredientData(context.Background(), domainIngredients); err != nil {
		return err
	}
	return nil
}
