package recipe

import (
	"context"
	"database/sql"
	"testing"

	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetRecipeSourceByRecipeVersionID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		r := NewRecipeVersionRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		versionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)
		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Recipe With Source", 30, 60, 8, 1)
		require.NoError(t, err)
		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		url := "https://example.com/recipe"
		err = seeds.InsertRecipeSource(ctx, tx, versionID, int(URL), &url, nil, nil, nil)
		require.NoError(t, err)

		// Act
		retrievedSource, err := r.getRecipeSourceByRecipeVersionID(ctx, tx, versionID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedSource, "Expected to find recipe source")
		require.Equal(t, URL, retrievedSource.Type)
		require.NotNil(t, retrievedSource.URL)
		require.Equal(t, url, *retrievedSource.URL)
	})
}

func TestCreateRecipeVersion_PersistsAnimalProductLevel(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r := NewRecipeVersionRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)

		recipeVersion := &RecipeVersion{
			ID:                 uuid.New(),
			RecipeID:           recipeID,
			Version:            1,
			Name:               "Version With Animal Level",
			Description:        "test",
			PrepMins:           10,
			CookMins:           20,
			Portions:           2,
			AnimalProductLevel: 2,
		}

		createdVersion, err := r.createRecipeVersion(ctx, tx, recipeVersion)
		require.NoError(t, err)
		require.NotNil(t, createdVersion)
		require.Equal(t, recipeVersion.AnimalProductLevel, createdVersion.AnimalProductLevel)

		var persistedAnimalProductLevel int
		err = tx.QueryRowContext(ctx, `SELECT animal_product_level FROM recipe_versions WHERE id = $1`, recipeVersion.ID).Scan(&persistedAnimalProductLevel)
		require.NoError(t, err)
		require.Equal(t, recipeVersion.AnimalProductLevel, persistedAnimalProductLevel)
	})
}

func TestGetRecipeVersionByID_ReturnsAnimalProductLevel(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r := NewRecipeVersionRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		versionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)
		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Recipe With Animal Level", 30, 60, 8, 1)
		require.NoError(t, err)

		expectedAnimalProductLevel := 1
		_, err = tx.ExecContext(ctx, `UPDATE recipe_versions SET animal_product_level = $1 WHERE id = $2`, expectedAnimalProductLevel, versionID)
		require.NoError(t, err)

		retrievedVersion, err := r.getRecipeVersionByID(ctx, tx, versionID)
		require.NoError(t, err)
		require.NotNil(t, retrievedVersion)
		require.Equal(t, expectedAnimalProductLevel, retrievedVersion.AnimalProductLevel)
	})
}

func TestCreateRecipeVersion_PersistsContainsGluten(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r := NewRecipeVersionRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)

		recipeVersion := &RecipeVersion{
			ID:             uuid.New(),
			RecipeID:       recipeID,
			Version:        1,
			Name:           "Version With Gluten Flag",
			Description:    "test",
			PrepMins:       10,
			CookMins:       20,
			Portions:       2,
			ContainsGluten: true,
		}

		createdVersion, err := r.createRecipeVersion(ctx, tx, recipeVersion)
		require.NoError(t, err)
		require.NotNil(t, createdVersion)
		require.Equal(t, recipeVersion.ContainsGluten, createdVersion.ContainsGluten)

		var persistedContainsGluten bool
		err = tx.QueryRowContext(ctx, `SELECT contains_gluten FROM recipe_versions WHERE id = $1`, recipeVersion.ID).Scan(&persistedContainsGluten)
		require.NoError(t, err)
		require.Equal(t, recipeVersion.ContainsGluten, persistedContainsGluten)
	})
}

func TestGetRecipeVersionByID_ReturnsContainsGluten(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r := NewRecipeVersionRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		versionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)
		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Recipe With Gluten", 30, 60, 8, 1)
		require.NoError(t, err)

		expectedContainsGluten := true
		_, err = tx.ExecContext(ctx, `UPDATE recipe_versions SET contains_gluten = $1 WHERE id = $2`, expectedContainsGluten, versionID)
		require.NoError(t, err)

		retrievedVersion, err := r.getRecipeVersionByID(ctx, tx, versionID)
		require.NoError(t, err)
		require.NotNil(t, retrievedVersion)
		require.Equal(t, expectedContainsGluten, retrievedVersion.ContainsGluten)
	})
}
