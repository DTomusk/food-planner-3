package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/logging"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	txRunner            db.TxRunner
	recipeRepo          *recipeRepo
	recipeVersionRepo   *recipeVersionRepo
	ingredientService   *ingredient.IngredientService
	ingredientUsageRepo *ingredientUsageRepo
	recipeRetentionDays *int
}

func NewService(
	txRunner db.TxRunner,
	repo *recipeRepo,
	recipeVersionRepo *recipeVersionRepo,
	ingredientService *ingredient.IngredientService,
	ingredientUsageRepo *ingredientUsageRepo,
	recipeRetentionDays *int,
) *Service {
	return &Service{
		txRunner:            txRunner,
		recipeRepo:          repo,
		recipeVersionRepo:   recipeVersionRepo,
		ingredientService:   ingredientService,
		ingredientUsageRepo: ingredientUsageRepo,
		recipeRetentionDays: recipeRetentionDays,
	}
}

func (s *Service) CreateRecipe(ctx context.Context, request CreateRecipeRequest) (*RecipeContainer, error) {
	logger := logging.FromContext(ctx).With("method", "CreateRecipe", "request", request)
	logger.Debug("Creating recipe")
	// Basic request validation
	// We might want to enforce unique name per user
	ingredientUsages, recipeSource, err := s.validateRecipeRequest(ctx, logger, request)
	if err != nil {
		logger.Error("Error validating recipe request", "error", err)
		return nil, err
	}

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		logger.Error("Error parsing user ID", "error", err)
		return nil, err
	}

	recipeContainer, err := NewRecipe(
		request.Name,
		userID,
		ingredientUsages,
		request.PrepMins,
		request.CookMins,
		request.Portions,
		recipeSource,
	)
	if err != nil {
		logger.Error("Error creating recipe", "error", err)
		return nil, err
	}

	// Persist recipe
	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		recipeContainer, err := s.recipeRepo.createRecipe(ctx, tx, recipeContainer)
		if err != nil {
			return err
		}
		err = s.ingredientUsageRepo.insertIngredientUsages(ctx, tx, ingredientUsages, recipeContainer.CurrentVersion.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("Error persisting recipe", "error", err)
		return nil, err
	}
	return recipeContainer, nil
}

func (s *Service) validateRecipeRequest(ctx context.Context, logger *slog.Logger, request CreateRecipeRequest) ([]*IngredientUsage, *RecipeSource, error) {
	if len(strings.TrimSpace(request.Name)) == 0 {
		return nil, nil, ErrEmptyName
	}
	if len(request.Ingredients) == 0 {
		return nil, nil, ErrNoIngredients
	}

	ingredientUsages, err := s.validateAndConvertIngredientUsages(ctx, logger, request.Ingredients)
	if err != nil {
		return nil, nil, err
	}

	recipeSource, err := newSource(&request.Source)
	if err != nil {
		return nil, nil, err
	}

	return ingredientUsages, recipeSource, nil
}

func (s *Service) UpdateRecipe(ctx context.Context, request UpdateRecipeRequest) (*RecipeContainer, error) {
	// Get recipe container and ensure user owns the recipe
	logger := logging.FromContext(ctx).With("method", "UpdateRecipe", "request", request)
	existingRecipe, err := s.recipeRepo.getRecipeByID(ctx, s.txRunner.DB(), request.RecipeId)
	if err != nil {
		return nil, err
	}
	if existingRecipe == nil {
		return nil, ErrRecipeNotFound
	}
	if existingRecipe.UserID.String() != request.Request.UserID {
		return nil, ErrUnauthorized
	}
	// Validate request
	ingredientUsages, recipeSource, err := s.validateRecipeRequest(ctx, logger, request.Request)
	if err != nil {
		return nil, err
	}
	// Instantiate entity
	recipeVersion, err := NewRecipeVersion(
		existingRecipe.ID,
		existingRecipe.CurrentVersion.Version+1,
		request.Request.Name,
		ingredientUsages,
		request.Request.PrepMins,
		request.Request.CookMins,
		request.Request.Portions,
		recipeSource,
	)
	if err != nil {
		return nil, err
	}
	// Persist
	var dbRecipeContainer *RecipeContainer
	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		// Create new recipe version
		dbVersion, err := s.recipeVersionRepo.createRecipeVersion(ctx, tx, recipeVersion)
		if err != nil {
			return err
		}
		// Insert ingredients
		err = s.ingredientUsageRepo.insertIngredientUsages(ctx, tx, ingredientUsages, recipeVersion.ID)
		if err != nil {
			return err
		}
		// Update current version id on recipe container
		err = s.recipeRepo.updateRecipeCurrentVersion(ctx, tx, existingRecipe.ID, recipeVersion.ID)
		if err != nil {
			return err
		}
		// Return updated recipe container
		dbRecipeContainer, err = s.recipeRepo.getRecipeByID(ctx, tx, existingRecipe.ID.String())
		dbRecipeContainer.CurrentVersion = dbVersion
		return nil
	})
	if err != nil {
		logger.Error("Error persisting recipe", "error", err)
		return nil, err
	}
	return dbRecipeContainer, nil
}

// service.go
func (s *Service) validateAndConvertIngredientUsages(ctx context.Context, logger *slog.Logger, ingredientUsageRequests []CreateIngredientUsageRequest) ([]*IngredientUsage, error) {
	seen := make(map[string]struct{}, len(ingredientUsageRequests))
	uniqueIDs := make([]string, 0, len(ingredientUsageRequests))
	for _, req := range ingredientUsageRequests {
		if _, exists := seen[req.IngredientID]; exists {
			return nil, ErrDuplicateIngredient
		}
		seen[req.IngredientID] = struct{}{}
		uniqueIDs = append(uniqueIDs, req.IngredientID)
	}

	dbIngredients, err := s.ingredientService.GetIngredientsByIDs(ctx, logger, uniqueIDs)
	if err != nil {
		return nil, err
	}
	if len(dbIngredients) != len(uniqueIDs) {
		return nil, ErrIngredientNotFound
	}

	ingredientByID := make(map[string]*ingredient.Ingredient, len(dbIngredients))
	for _, ing := range dbIngredients {
		ingredientByID[ing.ID.String()] = ing
	}

	return newIngredientUsages(ingredientUsageRequests, ingredientByID)
}

func (s *Service) GetAllRecipes(ctx context.Context) ([]*RecipeContainer, error) {
	return s.recipeRepo.getAllRecipes(ctx, s.txRunner.DB())
}

func (s *Service) GetRecipeByID(ctx context.Context, id string) (*RecipeContainer, error) {
	return s.recipeRepo.getRecipeByID(ctx, s.txRunner.DB(), id)
}

func (s *Service) GetRecipeVersionByID(ctx context.Context, id string) (*RecipeVersion, error) {
	return s.recipeVersionRepo.getRecipeVersionByID(ctx, s.txRunner.DB(), id)
}

func (s *Service) GetRecipesByUserID(ctx context.Context, userID string) ([]*RecipeContainer, error) {
	return s.recipeRepo.getRecipesByUserID(ctx, s.txRunner.DB(), userID)
}

func (s *Service) GetIngredientUsagesByRecipeVersionID(ctx context.Context, recipeVersionID string) ([]*IngredientUsage, error) {
	return s.ingredientUsageRepo.getIngredientUsagesForRecipeVersion(ctx, s.txRunner.DB(), recipeVersionID)
}

func (s *Service) GetRecipeSourceByRecipeVersionID(ctx context.Context, recipeVersionID string) (*RecipeSource, error) {
	return s.recipeVersionRepo.getRecipeSourceByRecipeVersionID(ctx, s.txRunner.DB(), recipeVersionID)
}

func (s *Service) GetRecipeVersionsByRecipeID(ctx context.Context, recipeID string) ([]*RecipeVersion, error) {
	return s.recipeVersionRepo.getRecipeVersionsByRecipeID(ctx, s.txRunner.DB(), recipeID)
}
