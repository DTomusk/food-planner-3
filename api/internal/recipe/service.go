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
		_, err := s.recipeRepo.createRecipeContainer(ctx, tx, recipeContainer)
		if err != nil {
			return err
		}

		_, err = s.recipeVersionRepo.createRecipeVersion(ctx, tx, recipeContainer.CurrentVersion)
		if err != nil {
			return err
		}

		err = s.recipeRepo.updateRecipeCurrentVersion(ctx, tx, recipeContainer.ID, recipeContainer.CurrentVersion.ID)
		if err != nil {
			return err
		}

		err = s.ingredientUsageRepo.insertIngredientUsages(ctx, tx, ingredientUsages, recipeContainer.CurrentVersion.ID)
		if err != nil {
			return err
		}

		if recipeContainer.CurrentVersion.Source != nil {
			err = s.recipeVersionRepo.insertRecipeSource(ctx, tx, recipeContainer.CurrentVersion.ID, recipeContainer.CurrentVersion.Source)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("Error persisting recipe", "error", err)
		return nil, err
	}

	// Consider returning values as they're inserted rather than doing a separate query
	persistedRecipe, err := s.recipeRepo.getRecipeByID(ctx, s.txRunner.DB(), recipeContainer.ID)
	if err != nil {
		logger.Error("Error retrieving persisted recipe", "error", err)
		return nil, err
	}

	return persistedRecipe, nil
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
	existingRecipe, err := s.recipeRepo.getRecipeByID(ctx, s.txRunner.DB(), uuid.MustParse(request.RecipeId))
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
	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		// Create new recipe version
		_, err := s.recipeVersionRepo.createRecipeVersion(ctx, tx, recipeVersion)
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

		if recipeVersion.Source != nil {
			err = s.recipeVersionRepo.insertRecipeSource(ctx, tx, recipeVersion.ID, recipeVersion.Source)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("Error persisting recipe", "error", err)
		return nil, err
	}

	// Return updated recipe container
	dbRecipeContainer, err := s.recipeRepo.getRecipeByID(ctx, s.txRunner.DB(), existingRecipe.ID)
	if err != nil {
		logger.Error("Error retrieving updated recipe", "error", err)
		return nil, err
	}

	return dbRecipeContainer, nil
}

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

const (
	defaultRecipePageSize = 20
	maxRecipePageSize     = 100
)

func (s *Service) GetRecipes(ctx context.Context, params RecipeListParams) ([]*RecipeWithCursor, *string, error) {
	count := params.Pagination.First

	// Default limits if invalid
	if count <= 0 {
		count = defaultRecipePageSize
	}
	if count > maxRecipePageSize {
		count = maxRecipePageSize
	}

	// Extract query from filter object
	var query *string
	if params.Filter.Query != nil && *params.Filter.Query != "" {
		trimmed := strings.TrimSpace(*params.Filter.Query)
		if trimmed != "" {
			query = &trimmed
		}
	}

	// If a query is included, then set the mode to relevance, otherwise use the default newest mode
	mode := RecipeCursorModeNewest
	if query != nil {
		mode = RecipeCursorModeRelevance
	}
	filterHash := filterHashForParams(mode, query, params.Filter.UserID)

	cursor := params.Pagination.After
	c, err := ParseRecipeCursor(cursor)
	if err != nil {
		return nil, nil, err
	}

	// Matches compares the mode and filter hash of the cursor to the current request
	// If they don't match, then we nullify the cursor and query from the start
	if !c.Matches(mode, filterHash) {
		c = nil
	}

	// Query recipes
	var rows []*RecipeListRow
	switch mode {
	case RecipeCursorModeNewest:
		rows, err = s.recipeRepo.getRecipesByCreatedAt(ctx, s.txRunner.DB(), count+1, c)
	case RecipeCursorModeRelevance:
		// query can be safely dereferenced as checked above
		rows, err = s.recipeRepo.getRecipesByRelevance(ctx, s.txRunner.DB(), *query, count+1, c)
	default:
		return nil, nil, ErrInvalidCursor
	}

	if err != nil {
		return nil, nil, err
	}

	// Once we've done the query, set the starting cursor for the next page
	var nextCursor *string
	if len(rows) > count {
		rows = rows[:count]
		last := rows[len(rows)-1]
		encoded, err := recipeContainerCursor(last, mode, filterHash)
		if err != nil {
			return nil, nil, err
		}
		nextCursor = &encoded
	}

	// Construct the response objects with cursors
	recipes := make([]*RecipeWithCursor, len(rows))
	for i, row := range rows {
		encoded, err := recipeContainerCursor(row, mode, filterHash)
		if err != nil {
			return nil, nil, err
		}
		recipes[i] = &RecipeWithCursor{
			Recipe: row.Recipe,
			Cursor: encoded,
		}
	}

	return recipes, nextCursor, nil
}

func recipeContainerCursor(
	row *RecipeListRow,
	mode RecipeCursorMode,
	filterHash string,
) (string, error) {
	if row == nil {
		return "", ErrInvalidCursor
	}

	if row.Recipe == nil {
		return "", ErrInvalidCursor
	}

	cursor := &RecipeCursor{
		Mode:       mode,
		FilterHash: filterHash,
		CreatedAt:  row.Recipe.CreatedAt,
		ID:         row.Recipe.ID,
	}

	if mode == RecipeCursorModeRelevance && row.RelevanceScore != nil {
		cursor.RelevanceScore = row.RelevanceScore
	}

	encoded := cursor.String()
	if encoded == "" {
		return "", ErrInvalidCursor
	}

	return encoded, nil
}

func (s *Service) GetRecipeByID(ctx context.Context, id uuid.UUID) (*RecipeContainer, error) {
	return s.recipeRepo.getRecipeByID(ctx, s.txRunner.DB(), id)
}

func (s *Service) GetRecipeVersionByID(ctx context.Context, id uuid.UUID) (*RecipeVersion, error) {
	return s.recipeVersionRepo.getRecipeVersionByID(ctx, s.txRunner.DB(), id)
}

func (s *Service) GetRecipeVersionByRecipeIDAndVersion(ctx context.Context, id uuid.UUID, version int) (*RecipeVersion, error) {
	return s.recipeVersionRepo.getRecipeVersionByRecipeIDAndVersion(ctx, s.txRunner.DB(), id, version)
}

func (s *Service) GetRecipesByUserID(ctx context.Context, userID uuid.UUID) ([]*RecipeContainer, error) {
	return s.recipeRepo.getRecipesByUserID(ctx, s.txRunner.DB(), userID)
}

func (s *Service) GetIngredientUsagesByRecipeVersionID(ctx context.Context, recipeVersionID uuid.UUID) ([]*IngredientUsage, error) {
	return s.ingredientUsageRepo.getIngredientUsagesForRecipeVersion(ctx, s.txRunner.DB(), recipeVersionID)
}

func (s *Service) GetRecipeSourceByRecipeVersionID(ctx context.Context, recipeVersionID uuid.UUID) (*RecipeSource, error) {
	return s.recipeVersionRepo.getRecipeSourceByRecipeVersionID(ctx, s.txRunner.DB(), recipeVersionID)
}

func (s *Service) GetRecipeVersionsByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]*RecipeVersion, error) {
	return s.recipeVersionRepo.getRecipeVersionsByRecipeID(ctx, s.txRunner.DB(), recipeID)
}
