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
	if len(strings.TrimSpace(request.Name)) == 0 {
		return nil, ErrEmptyName
	}
	if len(request.Ingredients) == 0 {
		return nil, ErrNoIngredients
	}

	ingredientUsages, err := s.validateAndConvertIngredientUsages(ctx, logger, request.Ingredients)
	if err != nil {
		return nil, err
	}

	recipeSource, err := s.validateAndConvertRecipeSource(&request.Source)
	if err != nil {
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

	var fetchedContainer *RecipeContainer

	// Persist recipe
	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		fetchedContainer, err = s.recipeRepo.createRecipe(ctx, tx, recipeContainer)
		if err != nil {
			return err
		}
		err = s.ingredientUsageRepo.insertIngredientUsages(ctx, tx, recipeContainer.CurrentVersion.Ingredients, recipeContainer.CurrentVersion.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("Error persisting recipe", "error", err)
		return nil, err
	}
	return fetchedContainer, nil
}

func (s *Service) validateAndConvertIngredientUsages(ctx context.Context, logger *slog.Logger, ingredientUsageRequests []CreateIngredientUsageRequest) ([]*IngredientUsage, error) {
	// Validate ingredients and check for duplicates
	// TODO: There may be cases where we allow duplicate ingredients (e.g. different sections of the same recipe)
	// But then ingredient usage table would need a new primary key as composite would no longer be unique
	seenIngredients := make(map[string]bool)
	ingredientUsages := make([]*IngredientUsage, len(ingredientUsageRequests))
	for i, ingredientRequest := range ingredientUsageRequests {
		if seenIngredients[ingredientRequest.IngredientID] {
			return nil, ErrDuplicateIngredient
		}
		seenIngredients[ingredientRequest.IngredientID] = true
		// Ensure ingredient exists and get preferred unit for validation
		// TODO: could batch this to reduce queries
		ingredients, err := s.ingredientService.GetIngredientsByIDs(ctx, logger, []string{ingredientRequest.IngredientID})
		if err != nil {
			logger.Error("Error checking ingredient existence", "ingredient_id", ingredientRequest.IngredientID, "error", err)
			return nil, err
		}
		if len(ingredients) == 0 {
			return nil, ErrIngredientNotFound
		}
		selectedIngredient := ingredients[0]
		if ingredientRequest.Unit != int(selectedIngredient.PreferredUnit) {
			return nil, ErrInvalidUnit
		}
		// TODO: get ingredient rule and validate usage at instantiation
		usage, err := NewIngredientUsage(ingredientRequest)
		if err != nil {
			return nil, err
		}
		ingredientUsages[i] = usage
	}
	return ingredientUsages, nil
}

func (s *Service) validateAndConvertRecipeSource(sourceRequest *CreateRecipeSourceRequest) (*RecipeSource, error) {
	if sourceRequest == nil {
		return nil, ErrNoSource
	}
	switch sourceRequest.Type {
	case 0:
		return nil, nil
	case 1:
		if sourceRequest.URL == nil {
			return nil, ErrMissingURL
		}
		return NewURLSource(*sourceRequest.URL)
	case 2:
		if sourceRequest.BookTitle == nil || sourceRequest.BookPage == nil {
			return nil, ErrMissingBookReference
		}
		return NewBookReferenceSource(*sourceRequest.BookTitle, *sourceRequest.BookPage)
	case 3:
		if sourceRequest.Instructions == nil {
			return nil, ErrMissingInstructions
		}
		return NewOriginalSource(*sourceRequest.Instructions)
	default:
		return nil, ErrInvalidSourceType
	}
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
