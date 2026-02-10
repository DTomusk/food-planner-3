package recipe

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/logging"
	"strings"
)

type Service struct {
	db                db.DBTX
	Repo              *Repo
	IngredientService *ingredient.IngredientService
}

func NewService(db db.DBTX, repo *Repo, ingredientService *ingredient.IngredientService) *Service {
	return &Service{
		db:                db,
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

	// Validate ingredients and check for duplicates
	// There may be cases where we allow duplicate ingredients (e.g. different sections of the same recipe)
	seenIngredients := make(map[string]bool)
	ingredientUsages := make([]*IngredientUsage, len(request.Ingredients))
	for i, ingredientRequest := range request.Ingredients {
		if seenIngredients[ingredientRequest.IngredientID] {
			return nil, ErrDuplicateIngredient
		}
		seenIngredients[ingredientRequest.IngredientID] = true
		// ensure ingredient exists (in future, grab ingredient validation rules)
		exists, err := s.IngredientService.Exists(ctx, s.db, ingredientRequest.IngredientID)
		if err != nil {
			logger.Error("Error checking ingredient existence", "ingredient_id", ingredientRequest.IngredientID, "error", err)
			return nil, err
		}
		if !exists {
			return nil, ErrIngredientNotFound
		}
		// TODO: get ingredient rule and validate usage at instantiation
		usage, err := NewIngredientUsage(ingredientRequest)
		if err != nil {
			return nil, err
		}
		ingredientUsages[i] = usage
	}

	// Once we've confirmed all ingredients are valid, we can create the recipe
	recipe, err := NewRecipe(request.Name, ingredientUsages)
	if err != nil {
		logger.Error("Error creating recipe", "error", err)
		return nil, err
	}

	// Persist recipe
	// We should use a transaction
	recipe, err = s.Repo.CreateRecipe(ctx, s.db, recipe)
	if err != nil {
		logger.Error("Error persisting recipe", "error", err)
		return nil, err
	}
	return recipe, nil
}

func (s *Service) GetAllRecipes(ctx context.Context) ([]*Recipe, error) {
	return s.Repo.GetAllRecipes(ctx, s.db)
}

func (s *Service) GetRecipeByID(ctx context.Context, id string) (*Recipe, error) {
	return s.Repo.GetRecipeByID(ctx, s.db, id)
}
