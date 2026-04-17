package graph

import "foodplanner/internal/gql/graph/model"

const (
	defaultPageSize              = 20
	maxPageSizeForCost           = 100
	DefaultMaxAcceptedComplexity = 1200
)

func NewComplexityRoot() ComplexityRoot {
	complexity := ComplexityRoot{}

	complexity.Query.Recipes = func(childComplexity int, pagination *model.PaginationInput, filter *model.RecipeFilterInput) int {
		first := pageSizeForComplexity(pagination)
		return 5 + first*childComplexity
	}

	complexity.User.Recipes = func(childComplexity int, pagination *model.PaginationInput, filter *model.RecipeFilterInput) int {
		first := pageSizeForComplexity(pagination)
		return 5 + first*childComplexity
	}

	complexity.Query.Ingredients = func(childComplexity int) int {
		return 15 + childComplexity
	}

	complexity.Recipe.Versions = func(childComplexity int) int {
		return 20 + 10*childComplexity
	}

	complexity.RecipeVersion.IngredientUsages = func(childComplexity int) int {
		return 20 + 8*childComplexity
	}

	complexity.Mutation.CreateRecipe = func(childComplexity int, input model.CreateRecipeInput) int {
		ingredientCost := len(input.IngredientUsages) * 2
		return 40 + ingredientCost + childComplexity
	}

	complexity.Mutation.UpdateRecipe = func(childComplexity int, input model.UpdateRecipeInput) int {
		ingredientCost := 0
		if input.Details != nil {
			ingredientCost = len(input.Details.IngredientUsages) * 2
		}
		return 35 + ingredientCost + childComplexity
	}

	complexity.Mutation.Signin = func(childComplexity int, input model.SignInInput) int {
		return 30 + childComplexity
	}

	complexity.Mutation.Signup = func(childComplexity int, input model.SignUpInput) int {
		return 35 + childComplexity
	}

	complexity.Mutation.Refresh = func(childComplexity int) int {
		return 20 + childComplexity
	}

	return complexity
}

func pageSizeForComplexity(p *model.PaginationInput) int {
	if p == nil || p.First == nil {
		return defaultPageSize
	}

	n := int(*p.First)
	if n < 1 {
		return 1
	}
	if n > maxPageSizeForCost {
		return maxPageSizeForCost
	}
	return n
}
