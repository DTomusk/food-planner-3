package graph

import (
	"testing"

	"foodplanner/internal/gql/graph/model"

	"github.com/stretchr/testify/require"
)

func ptrInt32(v int32) *int32 {
	return &v
}

func TestPageSizeForComplexity(t *testing.T) {
	t.Run("defaults when pagination is nil", func(t *testing.T) {
		require.Equal(t, defaultPageSize, pageSizeForComplexity(nil))
	})

	t.Run("defaults when first is nil", func(t *testing.T) {
		require.Equal(t, defaultPageSize, pageSizeForComplexity(&model.PaginationInput{}))
	})

	t.Run("clamps to minimum", func(t *testing.T) {
		first := int32(0)
		require.Equal(t, 1, pageSizeForComplexity(&model.PaginationInput{First: &first}))
	})

	t.Run("clamps to maximum", func(t *testing.T) {
		first := int32(1000)
		require.Equal(t, maxPageSizeForCost, pageSizeForComplexity(&model.PaginationInput{First: &first}))
	})

	t.Run("uses provided first", func(t *testing.T) {
		first := int32(37)
		require.Equal(t, 37, pageSizeForComplexity(&model.PaginationInput{First: &first}))
	})
}

func TestComplexityCalculations(t *testing.T) {
	c := NewComplexityRoot()

	t.Run("query recipes uses default page size", func(t *testing.T) {
		score := c.Query.Recipes(3, nil, nil)
		require.Equal(t, 5+defaultPageSize*3, score)
	})

	t.Run("query recipes clamps page size to max", func(t *testing.T) {
		score := c.Query.Recipes(2, &model.PaginationInput{First: ptrInt32(500)}, nil)
		require.Equal(t, 5+maxPageSizeForCost*2, score)
	})

	t.Run("user recipes uses requested page size", func(t *testing.T) {
		score := c.User.Recipes(4, &model.PaginationInput{First: ptrInt32(7)}, nil)
		require.Equal(t, 5+7*4, score)
	})

	t.Run("ingredients query adds flat cost", func(t *testing.T) {
		require.Equal(t, 19, c.Query.Ingredients(4))
	})

	t.Run("recipe versions multiplies child cost", func(t *testing.T) {
		require.Equal(t, 50, c.Recipe.Versions(3))
	})

	t.Run("ingredient usages multiplies child cost", func(t *testing.T) {
		require.Equal(t, 44, c.RecipeVersion.IngredientUsages(3))
	})

	t.Run("create recipe includes ingredient list size", func(t *testing.T) {
		input := model.CreateRecipeInput{
			IngredientUsages: []*model.CreateIngredientUsageInput{{}, {}, {}},
		}
		require.Equal(t, 51, c.Mutation.CreateRecipe(5, input))
	})

	t.Run("update recipe handles nil details", func(t *testing.T) {
		input := model.UpdateRecipeInput{Details: nil}
		require.Equal(t, 39, c.Mutation.UpdateRecipe(4, input))
	})

	t.Run("update recipe counts nested ingredient usages", func(t *testing.T) {
		input := model.UpdateRecipeInput{
			Details: &model.CreateRecipeInput{
				IngredientUsages: []*model.CreateIngredientUsageInput{{}, {}, {}, {}},
			},
		}
		require.Equal(t, 45, c.Mutation.UpdateRecipe(2, input))
	})

	t.Run("signin/signup/refresh use fixed base", func(t *testing.T) {
		require.Equal(t, 37, c.Mutation.Signin(7, model.SignInInput{}))
		require.Equal(t, 42, c.Mutation.Signup(7, model.SignUpInput{}))
		require.Equal(t, 27, c.Mutation.Refresh(7))
	})
}
