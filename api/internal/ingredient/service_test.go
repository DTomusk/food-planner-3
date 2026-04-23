package ingredient

import (
	"database/sql"
	"foodplanner/internal/logging"
	"foodplanner/internal/testutil"
	"foodplanner/internal/unit"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSyncIngredientData_PersistsAnimalProductLevelAcrossBatches(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		service := NewIngredientService(testutil.NewTestTxRunner(tx), NewIngredientRepo(), 2)

		ingredients := []*Ingredient{
			{
				ID:                 uuid.New(),
				Name:               "Service Vegan Ingredient",
				FileKey:            "service_vegan_ingredient",
				PreferredUnit:      unit.Quantum,
				AnimalProductLevel: Vegan,
				ProcessedLevel:     Raw,
			},
			{
				ID:                 uuid.New(),
				Name:               "Service Vegetarian Ingredient",
				FileKey:            "service_vegetarian_ingredient",
				PreferredUnit:      unit.Gram,
				AnimalProductLevel: Vegetarian,
				ProcessedLevel:     Derived,
			},
			{
				ID:                 uuid.New(),
				Name:               "Service Meat Ingredient",
				FileKey:            "service_meat_ingredient",
				PreferredUnit:      unit.Gram,
				AnimalProductLevel: Meat,
				ProcessedLevel:     Derived,
			},
		}

		err := service.SyncIngredientData(ctx, logger, ingredients)
		require.NoError(t, err)

		ids := []string{ingredients[0].ID.String(), ingredients[1].ID.String(), ingredients[2].ID.String()}
		persisted, err := service.GetIngredientsByIDs(ctx, logger, ids)
		require.NoError(t, err)
		require.Len(t, persisted, 3)

		persistedByID := make(map[uuid.UUID]*Ingredient, len(persisted))
		for _, ingredient := range persisted {
			persistedByID[ingredient.ID] = ingredient
		}

		require.Equal(t, Vegan, persistedByID[ingredients[0].ID].AnimalProductLevel)
		require.Equal(t, Vegetarian, persistedByID[ingredients[1].ID].AnimalProductLevel)
		require.Equal(t, Meat, persistedByID[ingredients[2].ID].AnimalProductLevel)
	})
}

func TestExists_ReturnsTrueForPersistedIngredient(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		service := NewIngredientService(testutil.NewTestTxRunner(tx), NewIngredientRepo(), 10)

		ingredientID := uuid.New()
		err := service.SyncIngredientData(ctx, logger, []*Ingredient{
			{
				ID:                 ingredientID,
				Name:               "Service Exists Ingredient",
				FileKey:            "service_exists_ingredient",
				PreferredUnit:      unit.Quantum,
				AnimalProductLevel: Vegan,
				ProcessedLevel:     Raw,
			},
		})
		require.NoError(t, err)

		exists, err := service.Exists(ctx, ingredientID.String())
		require.NoError(t, err)
		require.True(t, exists)

		notExists, err := service.Exists(ctx, uuid.NewString())
		require.NoError(t, err)
		require.False(t, notExists)
	})
}

func TestGetIngredientsByIDs_EmptyIDsReturnsEmptySlice(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		service := NewIngredientService(testutil.NewTestTxRunner(tx), NewIngredientRepo(), 10)

		ingredients, err := service.GetIngredientsByIDs(ctx, logger, []string{})
		require.NoError(t, err)
		require.Empty(t, ingredients)
	})
}

func TestServiceGetAllIngredients_FiltersOutNonSearchable(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		service := NewIngredientService(testutil.NewTestTxRunner(tx), NewIngredientRepo(), 10)

		searchableID := uuid.New()
		nonSearchableID := uuid.New()

		err := service.SyncIngredientData(ctx, logger, []*Ingredient{
			{
				ID:                 searchableID,
				Name:               "Service Searchable Ingredient",
				FileKey:            "service_searchable_ingredient",
				PreferredUnit:      unit.Gram,
				AnimalProductLevel: Vegan,
				ProcessedLevel:     Raw,
				IsSearchable:       true,
			},
			{
				ID:                 nonSearchableID,
				Name:               "Service Non Searchable Ingredient",
				FileKey:            "service_non_searchable_ingredient",
				PreferredUnit:      unit.UnitUnknown,
				AnimalProductLevel: Vegan,
				ProcessedLevel:     Raw,
				IsSearchable:       false,
			},
		})
		require.NoError(t, err)

		allIngredients, err := service.GetAllIngredients(ctx, logger)
		require.NoError(t, err)
		require.Len(t, allIngredients, 1)
		require.Equal(t, searchableID, allIngredients[0].ID)
		require.True(t, allIngredients[0].IsSearchable)
	})
}

func TestGetIngredientsByIDs_ReturnsNonSearchableIngredientWhenRequested(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		service := NewIngredientService(testutil.NewTestTxRunner(tx), NewIngredientRepo(), 10)

		ingredientID := uuid.New()
		err := service.SyncIngredientData(ctx, logger, []*Ingredient{
			{
				ID:                 ingredientID,
				Name:               "Service Hidden Ingredient",
				FileKey:            "service_hidden_ingredient",
				PreferredUnit:      unit.UnitUnknown,
				AnimalProductLevel: Vegan,
				ProcessedLevel:     Derived,
				IsSearchable:       false,
			},
		})
		require.NoError(t, err)

		ingredients, err := service.GetIngredientsByIDs(ctx, logger, []string{ingredientID.String()})
		require.NoError(t, err)
		require.Len(t, ingredients, 1)
		require.Equal(t, ingredientID, ingredients[0].ID)
		require.False(t, ingredients[0].IsSearchable)
		require.Equal(t, unit.UnitUnknown, ingredients[0].PreferredUnit)
	})
}

func TestGetAllIngredientsUnfiltered_IncludesNonSearchable(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		service := NewIngredientService(testutil.NewTestTxRunner(tx), NewIngredientRepo(), 10)

		searchableID := uuid.New()
		nonSearchableID := uuid.New()

		err := service.SyncIngredientData(ctx, logger, []*Ingredient{
			{
				ID:                 searchableID,
				Name:               "Service Searchable Unfiltered",
				FileKey:            "service_searchable_unfiltered",
				PreferredUnit:      unit.Gram,
				AnimalProductLevel: Vegan,
				ProcessedLevel:     Raw,
				IsSearchable:       true,
			},
			{
				ID:                 nonSearchableID,
				Name:               "Service Non Searchable Unfiltered",
				FileKey:            "service_non_searchable_unfiltered",
				PreferredUnit:      unit.UnitUnknown,
				AnimalProductLevel: Vegan,
				ProcessedLevel:     Derived,
				IsSearchable:       false,
			},
		})
		require.NoError(t, err)

		ingredients, err := service.GetAllIngredientsUnfiltered(ctx, logger)
		require.NoError(t, err)

		byID := make(map[uuid.UUID]*Ingredient, len(ingredients))
		for _, ingredient := range ingredients {
			byID[ingredient.ID] = ingredient
		}

		require.Contains(t, byID, searchableID)
		require.Contains(t, byID, nonSearchableID)
		require.False(t, byID[nonSearchableID].IsSearchable)
		require.Equal(t, unit.UnitUnknown, byID[nonSearchableID].PreferredUnit)
	})
}
