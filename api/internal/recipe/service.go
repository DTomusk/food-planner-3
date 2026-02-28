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
	txRunner          db.TxRunner
	Repo              *Repo
	IngredientService *ingredient.IngredientService
}

func NewService(txRunner db.TxRunner, repo *Repo, ingredientService *ingredient.IngredientService) *Service {
	return &Service{
		txRunner:          txRunner,
		Repo:              repo,
		IngredientService: ingredientService,
	}
}

func (s *Service) CreateRecipe(ctx context.Context, request CreateRecipeRequest) (*Recipe, error) {
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

	recipe, err := NewRecipe(
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
	// We should use a transaction
	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		recipe, err = s.Repo.CreateRecipe(ctx, tx, recipe)
		return err
	})
	if err != nil {
		logger.Error("Error persisting recipe", "error", err)
		return nil, err
	}
	return recipe, nil
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
		ingredients, err := s.IngredientService.GetIngredientsByIDs(ctx, logger, []string{ingredientRequest.IngredientID})
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

func (s *Service) GetAllRecipes(ctx context.Context) ([]*Recipe, error) {
	return s.Repo.GetAllRecipes(ctx, s.txRunner.DB())
}

func (s *Service) GetRecipeByID(ctx context.Context, id string) (*Recipe, error) {
	return s.Repo.GetRecipeByID(ctx, s.txRunner.DB(), id)
}

func (s *Service) GetIngredientUsagesByRecipeID(ctx context.Context, recipeID string) ([]*IngredientUsage, error) {
	return s.Repo.GetIngredientUsagesForRecipe(ctx, s.txRunner.DB(), recipeID)
}

func (s *Service) GetRecipeSourceByRecipeID(ctx context.Context, recipeID string) (*RecipeSource, error) {
	return s.Repo.GetRecipeSourceByRecipeID(ctx, s.txRunner.DB(), recipeID)
}

func (s *Service) DeleteRecipe(ctx context.Context, recipeID, userID string) (*Recipe, error) {
	// Get the recipe to ensure it exists and belongs to the user
	recipe, err := s.GetRecipeByID(ctx, recipeID)
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, ErrRecipeNotFound
	}
	if recipe.UserID.String() != userID {
		return nil, ErrUnauthorized
	}
	dbRecipe, err := s.Repo.DeleteRecipe(ctx, s.txRunner.DB(), recipeID)
	if err != nil {
		return nil, err
	}
	return dbRecipe, nil
}
