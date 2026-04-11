package recipe

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type listedRecipeSeed struct {
	RecipeID  uuid.UUID
	VersionID uuid.UUID
	Name      string
	CreatedAt time.Time
}

func TestNewRepo_ValidatesWeights(t *testing.T) {
	_, err := NewRecipeRepo(-0.1, 1.1)
	require.Error(t, err)
	_, err = NewRecipeRepo(0.5, 0.6)
	require.Error(t, err)
	_, err = NewRecipeRepo(0.15, 0.85)
	require.NoError(t, err)
	_, err = NewRecipeRepo(0, 1)
	require.NoError(t, err)
	_, err = NewRecipeRepo(1, 0)
	require.NoError(t, err)
}

func TestGetRecipe_DoesNotErrorWhenNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		// Act
		recipeContainer, err := r.getRecipeByID(context.Background(), tx, uuid.MustParse("04061e4e-6d4c-41d1-abcf-8b214927e1ed"))

		// Assert
		require.NoError(t, err)
		require.Nil(t, recipeContainer)
	})
}

func TestGetRecipeByID_ReturnsCurrentVersionImgSrc(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		recipeID := uuid.New()
		versionID := uuid.New()
		expectedImgSrc := "https://cdn.example.com/recipe-images/test-user/test-upload.png"

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)

		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Image Recipe", 30, 60, 8, 1)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx, `UPDATE recipe_versions SET img_src = $1 WHERE id = $2`, expectedImgSrc, versionID)
		require.NoError(t, err)

		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		recipeContainer, err := r.getRecipeByID(ctx, tx, recipeID)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.NotNil(t, recipeContainer.CurrentVersion)
		require.NotNil(t, recipeContainer.CurrentVersion.ImgSrc)
		require.Equal(t, expectedImgSrc, *recipeContainer.CurrentVersion.ImgSrc)
	})
}

func TestGetRecipeByID_ReturnsCurrentVersionAnimalProductLevel(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		versionID := uuid.New()
		expectedAnimalProductLevel := 2

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)

		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Animal Level Recipe", 30, 60, 8, 1)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx, `UPDATE recipe_versions SET animal_product_level = $1 WHERE id = $2`, expectedAnimalProductLevel, versionID)
		require.NoError(t, err)

		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		recipeContainer, err := r.getRecipeByID(ctx, tx, recipeID)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.NotNil(t, recipeContainer.CurrentVersion)
		require.Equal(t, expectedAnimalProductLevel, recipeContainer.CurrentVersion.AnimalProductLevel)
	})
}

func TestGetRecipesByUserID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		otherUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed other test user")

		expectedRecipeNames := []string{"Test Recipe 0", "Test Recipe 1"}
		for _, recipeName := range expectedRecipeNames {
			recipeID := uuid.New()
			versionID := uuid.New()

			err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
			require.NoError(t, err)

			err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, recipeName, 30, 60, 8, 1)
			require.NoError(t, err)

			err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
			require.NoError(t, err)
		}

		otherRecipeID := uuid.New()
		otherVersionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, otherRecipeID, otherUser.ID)
		require.NoError(t, err)

		err = seeds.InsertRecipeVersion(ctx, tx, otherVersionID, otherRecipeID, "Other User Recipe", 30, 60, 8, 1)
		require.NoError(t, err)

		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, otherRecipeID, otherVersionID)
		require.NoError(t, err)

		// Act
		recipes, err := r.getRecipesByUserID(ctx, tx, testUser.ID)

		// Assert
		require.NoError(t, err)
		require.Len(t, recipes, 2, "Expected to find 2 recipes for test user")

		actualNames := make(map[string]struct{}, len(recipes))
		for _, recipe := range recipes {
			require.Equal(t, testUser.ID, recipe.UserID, "Expected all recipes to belong to test user")
			require.NotNil(t, recipe.CurrentVersion, "Expected recipe to include current version")
			require.Equal(t, recipe.ID, recipe.CurrentVersion.RecipeID, "Expected current version to point to recipe container")
			actualNames[recipe.CurrentVersion.Name] = struct{}{}
		}

		for _, expectedName := range expectedRecipeNames {
			_, found := actualNames[expectedName]
			require.True(t, found, "Expected recipe name to be present")
		}
	})
}

func TestGetRecipes_ReturnsActiveRecipesOrderedByCreatedAtAndIDDesc(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		sameCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
		olderCreatedAt := sameCreatedAt.Add(-1 * time.Minute)
		deletedCreatedAt := sameCreatedAt.Add(1 * time.Minute)
		deletedOn := deletedCreatedAt.Add(1 * time.Hour)

		newestHighID := uuid.MustParse("ffffffff-ffff-ffff-ffff-fffffffffff1")
		newestLowID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
		olderID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		deletedID := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")

		newestHigh := seedRecipeForListTests(t, ctx, tx, testUser.ID, newestHighID, uuid.New(), "Newest High", sameCreatedAt, nil, nil)
		newestLow := seedRecipeForListTests(t, ctx, tx, testUser.ID, newestLowID, uuid.New(), "Newest Low", sameCreatedAt, nil, nil)
		older := seedRecipeForListTests(t, ctx, tx, testUser.ID, olderID, uuid.New(), "Older", olderCreatedAt, nil, nil)
		seedRecipeForListTests(t, ctx, tx, testUser.ID, deletedID, uuid.New(), "Deleted", deletedCreatedAt, &deletedOn, nil)

		recipes, err := r.getRecipesByCreatedAt(ctx, tx, 10, nil, nil, nil)

		require.NoError(t, err)
		require.Len(t, recipes, 3)
		require.Equal(t, newestHigh.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, newestLow.RecipeID, recipes[1].Recipe.ID)
		require.Equal(t, older.RecipeID, recipes[2].Recipe.ID)
	})
}

func TestGetRecipes_AppliesCursorBoundary(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		sameCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
		olderCreatedAt := sameCreatedAt.Add(-1 * time.Minute)

		newestHigh := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("ffffffff-ffff-ffff-ffff-fffffffffff1"), uuid.New(), "Newest High", sameCreatedAt, nil, nil)
		newestLow := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("00000000-0000-0000-0000-000000000002"), uuid.New(), "Newest Low", sameCreatedAt, nil, nil)
		older := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.New(), "Older", olderCreatedAt, nil, nil)

		cursor := &RecipeCursor{
			CreatedAt: newestHigh.CreatedAt,
			ID:        newestHigh.RecipeID,
		}

		recipes, err := r.getRecipesByCreatedAt(ctx, tx, 10, cursor, nil, nil)

		require.NoError(t, err)
		require.Len(t, recipes, 2)
		require.Equal(t, newestLow.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, older.RecipeID, recipes[1].Recipe.ID)
	})
}

func TestGetRecipesByRelevance_ReturnsMatchesOrderedByScoreThenCreatedAtThenIDDesc(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		sameCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
		olderCreatedAt := sameCreatedAt.Add(-1 * time.Minute)

		exactHigh := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("ffffffff-ffff-ffff-ffff-fffffffffff1"), uuid.New(), "Chicken Soup", sameCreatedAt, nil, nil)
		exactLow := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("00000000-0000-0000-0000-000000000002"), uuid.New(), "Chicken Soup", sameCreatedAt, nil, nil)
		fuzzy := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.New(), "Chikcen Soup", olderCreatedAt, nil, nil)
		seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("22222222-2222-2222-2222-222222222222"), uuid.New(), "Beef Chili", olderCreatedAt, nil, nil)

		recipes, err := r.getRecipesByRelevance(ctx, tx, "chicken soup", 10, nil, nil, nil)

		require.NoError(t, err)
		require.Len(t, recipes, 3)
		require.Equal(t, exactHigh.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, exactLow.RecipeID, recipes[1].Recipe.ID)
		require.Equal(t, fuzzy.RecipeID, recipes[2].Recipe.ID)
		require.NotNil(t, recipes[0].RelevanceScore)
		require.NotNil(t, recipes[1].RelevanceScore)
		require.NotNil(t, recipes[2].RelevanceScore)
		require.GreaterOrEqual(t, *recipes[0].RelevanceScore, *recipes[2].RelevanceScore)
	})
}

func TestGetRecipesByRelevance_AppliesCursorBoundary(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		sameCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
		olderCreatedAt := sameCreatedAt.Add(-1 * time.Minute)

		seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("ffffffff-ffff-ffff-ffff-fffffffffff1"), uuid.New(), "Chicken Soup", sameCreatedAt, nil, nil)
		seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("00000000-0000-0000-0000-000000000002"), uuid.New(), "Chicken Soup", sameCreatedAt, nil, nil)
		fuzzy := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.New(), "Chikcen Soup", olderCreatedAt, nil, nil)

		firstPage, err := r.getRecipesByRelevance(ctx, tx, "chicken soup", 2, nil, nil, nil)
		require.NoError(t, err)
		require.Len(t, firstPage, 2)
		require.NotNil(t, firstPage[1].RelevanceScore)

		cursor := &RecipeCursor{
			CreatedAt:      firstPage[1].Recipe.CreatedAt,
			ID:             firstPage[1].Recipe.ID,
			RelevanceScore: firstPage[1].RelevanceScore,
		}

		secondPage, err := r.getRecipesByRelevance(ctx, tx, "chicken soup", 10, cursor, nil, nil)

		require.NoError(t, err)
		require.Len(t, secondPage, 1)
		require.Equal(t, fuzzy.RecipeID, secondPage[0].Recipe.ID)
	})
}

func TestGetRecipesByRelevance_ReturnsErrInvalidCursorWhenScoreMissing(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		cursor := &RecipeCursor{
			CreatedAt: time.Now().UTC(),
			ID:        uuid.New(),
		}

		recipes, err := r.getRecipesByRelevance(ctx, tx, "chicken", 10, cursor, nil, nil)

		require.ErrorIs(t, err, ErrInvalidCursor)
		require.Nil(t, recipes)
	})
}

func TestGetRecipesByRelevance_SupportsFuzzyMatching(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		createdAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
		expected := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("33333333-3333-3333-3333-333333333333"), uuid.New(), "Chicken Soup", createdAt, nil, nil)

		recipes, err := r.getRecipesByRelevance(ctx, tx, "chikcen soup", 10, nil, nil, nil)

		require.NoError(t, err)
		require.NotEmpty(t, recipes)
		require.Equal(t, expected.RecipeID, recipes[0].Recipe.ID)
		require.NotNil(t, recipes[0].RelevanceScore)
	})
}

func seedRecipeForListTests(
	t *testing.T,
	ctx context.Context,
	tx *sql.Tx,
	userID, recipeID, versionID uuid.UUID,
	name string,
	createdAt time.Time,
	deletedOn *time.Time,
	animalProductLevel *int,
) listedRecipeSeed {
	t.Helper()

	err := seeds.InsertRecipeContainer(ctx, tx, recipeID, userID)
	require.NoError(t, err)

	err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, name, 30, 60, 8, 1)
	require.NoError(t, err)

	err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
	require.NoError(t, err)

	_, err = tx.ExecContext(ctx, `UPDATE recipe_containers SET created_at = $1, deleted_on = $2 WHERE id = $3`, createdAt, deletedOn, recipeID)
	require.NoError(t, err)

	_, err = tx.ExecContext(ctx, `UPDATE recipe_versions SET created_at = $1 WHERE id = $2`, createdAt, versionID)
	require.NoError(t, err)

	return listedRecipeSeed{
		RecipeID:  recipeID,
		VersionID: versionID,
		Name:      name,
		CreatedAt: createdAt,
	}
}
