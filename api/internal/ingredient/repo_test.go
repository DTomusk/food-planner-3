package ingredient

import (
	"context"
	"database/sql"
	"foodplanner/internal/testutil"
	"foodplanner/internal/unit"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUpsertIngredients_PersistsAnimalProductLevel(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewIngredientRepo()

		ingredientA := &Ingredient{
			ID:                 uuid.New(),
			Name:               "Test Vegan Ingredient",
			FileKey:            "test_vegan_ingredient",
			PreferredUnit:      unit.Quantum,
			Counter:            nil,
			Plural:             nil,
			CounterPlural:      nil,
			AnimalProductLevel: Vegan,
			ContainsGluten:     false,
		}
		ingredientB := &Ingredient{
			ID:                 uuid.New(),
			Name:               "Test Meat Ingredient",
			FileKey:            "test_meat_ingredient",
			PreferredUnit:      unit.Gram,
			Counter:            nil,
			Plural:             nil,
			CounterPlural:      nil,
			AnimalProductLevel: Meat,
			ContainsGluten:     true,
		}

		err := repo.UpsertIngredients(ctx, tx, []*Ingredient{ingredientA, ingredientB})
		require.NoError(t, err)

		persisted, err := repo.GetIngredientsByIDs(ctx, tx, []string{ingredientA.ID.String(), ingredientB.ID.String()})
		require.NoError(t, err)
		require.Len(t, persisted, 2)

		persistedByID := make(map[uuid.UUID]*Ingredient, len(persisted))
		for _, ingredient := range persisted {
			persistedByID[ingredient.ID] = ingredient
		}

		require.Equal(t, Vegan, persistedByID[ingredientA.ID].AnimalProductLevel)
		require.Equal(t, Meat, persistedByID[ingredientB.ID].AnimalProductLevel)
		require.Equal(t, ingredientA.ContainsGluten, persistedByID[ingredientA.ID].ContainsGluten)
		require.Equal(t, ingredientB.ContainsGluten, persistedByID[ingredientB.ID].ContainsGluten)
	})
}

func TestUpsertIngredients_UpdatesExistingIngredientByFileKey(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewIngredientRepo()

		original := &Ingredient{
			ID:                 uuid.New(),
			Name:               "Test Ingredient Original",
			FileKey:            "test_update_by_file_key",
			PreferredUnit:      unit.Quantum,
			Counter:            nil,
			Plural:             nil,
			CounterPlural:      nil,
			AnimalProductLevel: Vegan,
			ContainsGluten:     false,
		}

		err := repo.UpsertIngredients(ctx, tx, []*Ingredient{original})
		require.NoError(t, err)

		updated := &Ingredient{
			ID:                 uuid.New(),
			Name:               "Test Ingredient Updated",
			FileKey:            original.FileKey,
			PreferredUnit:      unit.Gram,
			Counter:            testutil.PtrString("slice"),
			Plural:             testutil.PtrString("Test Ingredient Updateds"),
			CounterPlural:      testutil.PtrString("slices"),
			AnimalProductLevel: Vegetarian,
			ContainsGluten:     true,
		}

		err = repo.UpsertIngredients(ctx, tx, []*Ingredient{updated})
		require.NoError(t, err)

		persisted, err := repo.GetIngredientsByIDs(ctx, tx, []string{original.ID.String()})
		require.NoError(t, err)
		require.Len(t, persisted, 1)

		require.Equal(t, original.ID, persisted[0].ID)
		require.Equal(t, updated.Name, persisted[0].Name)
		require.Equal(t, updated.PreferredUnit, persisted[0].PreferredUnit)
		require.Equal(t, updated.Counter, persisted[0].Counter)
		require.Equal(t, updated.Plural, persisted[0].Plural)
		require.Equal(t, updated.CounterPlural, persisted[0].CounterPlural)
		require.Equal(t, updated.AnimalProductLevel, persisted[0].AnimalProductLevel)
		require.Equal(t, updated.ContainsGluten, persisted[0].ContainsGluten)
	})
}

func TestGetIngredientsByIDsRepo_EmptyIDsReturnsEmptySlice(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewIngredientRepo()

		ingredients, err := repo.GetIngredientsByIDs(context.Background(), tx, []string{})
		require.NoError(t, err)
		require.Empty(t, ingredients)
	})
}
