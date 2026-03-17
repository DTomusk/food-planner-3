package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/auth"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/recipe"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestUserResolver_User(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		r := &Resolver{
			UserService: user.NewUserService(tx, user.NewUserRepo()),
		}
		queryResolver := &queryResolver{r}

		userModel, err := queryResolver.User(ctx, testUser.ID.String())
		require.NoError(t, err, "Failed to fetch user")
		require.NotNil(t, userModel, "Expected user model, got nil")
		require.Equal(t, testUser.ID.String(), userModel.ID)
		require.Equal(t, testUser.Email, userModel.Email)
		require.Equal(t, testUser.Username, userModel.Username)
	})
}

func TestUserResolver_Me(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		r := &Resolver{
			UserService: user.NewUserService(tx, user.NewUserRepo()),
		}
		queryResolver := &queryResolver{r}

		claims := auth.Claims{UserID: testUser.ID.String()}
		ctx = auth.ContextWithClaims(ctx, &claims)

		userModel, err := queryResolver.Me(ctx)
		require.NoError(t, err, "Failed to fetch current user")
		require.NotNil(t, userModel, "Expected current user model, got nil")
		require.Equal(t, testUser.ID.String(), userModel.ID)
		require.Equal(t, testUser.Email, userModel.Email)
		require.Equal(t, testUser.Username, userModel.Username)
	})
}

func TestUserResolver_Me_Unauthenticated(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		r := &Resolver{
			UserService: user.NewUserService(tx, user.NewUserRepo()),
		}
		queryResolver := &queryResolver{r}

		userModel, err := queryResolver.Me(context.Background())
		require.Error(t, err, "Expected unauthenticated error")
		require.Nil(t, userModel)

		gqlErr, ok := err.(*gqlerror.Error)
		require.True(t, ok, "Expected GraphQL error type")
		require.Equal(t, "user is not authenticated", gqlErr.Message)
		require.Equal(t, "UNAUTHENTICATED", gqlErr.Extensions["code"])
	})
}

func TestUserResolver_Recipes(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		otherUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed other test user")

		expectedRecipeNames := []string{"User Recipe 1", "User Recipe 2"}
		for _, recipeName := range expectedRecipeNames {
			recipeID := uuid.New()
			versionID := uuid.New()

			err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
			require.NoError(t, err)

			err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, recipeName, 20, 40, 4, 1)
			require.NoError(t, err)

			err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
			require.NoError(t, err)
		}

		otherRecipeID := uuid.New()
		otherVersionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, otherRecipeID, otherUser.ID)
		require.NoError(t, err)

		err = seeds.InsertRecipeVersion(ctx, tx, otherVersionID, otherRecipeID, "Other User Recipe", 10, 30, 2, 1)
		require.NoError(t, err)

		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, otherRecipeID, otherVersionID)
		require.NoError(t, err)

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		recipeService := recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredientService, recipe.NewIngredientUsageRepo(), nil)

		r := &Resolver{
			RecipeService: recipeService,
		}
		userResolver := &userResolver{r}

		recipeModels, err := userResolver.Recipes(ctx, &model.User{ID: testUser.ID.String()})
		require.NoError(t, err, "Failed to fetch recipes for user")
		require.Len(t, recipeModels, 2, "Expected to find exactly 2 recipes for user")

		actualNames := make(map[string]struct{}, len(recipeModels))
		for _, recipeModel := range recipeModels {
			require.Equal(t, testUser.ID.String(), recipeModel.AuthorID, "Expected recipe to belong to test user")
			require.NotNil(t, recipeModel.CurrentVersion, "Expected current version to be populated")
			actualNames[recipeModel.CurrentVersion.Name] = struct{}{}
		}

		for _, expectedName := range expectedRecipeNames {
			_, found := actualNames[expectedName]
			require.True(t, found, "Expected recipe name to be present")
		}
	})
}

func TestUserResolver_Recipes_NoRecipes(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		recipeService := recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredientService, recipe.NewIngredientUsageRepo(), nil)

		r := &Resolver{
			RecipeService: recipeService,
		}
		userResolver := &userResolver{r}

		recipeModels, err := userResolver.Recipes(ctx, &model.User{ID: testUser.ID.String()})
		require.NoError(t, err, "Failed to fetch recipes for user")
		require.Len(t, recipeModels, 0, "Expected no recipes for user")
	})
}
