package resolver

import (
	"context"
	"fmt"
	"foodplanner/internal/auth"
	"foodplanner/internal/gql/graph/errors"
	"foodplanner/internal/gql/graph/middleware"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/recipe"
	"foodplanner/internal/user"
)

func mapUsers(users []*user.User) []*model.User {
	var result []*model.User

	for _, user := range users {
		result = append(result, mapUser(user))
	}

	return result
}
func mapUser(user *user.User) *model.User {
	if user == nil {
		return nil
	}
	return &model.User{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}
}

func mapRecipes(recipes []*recipe.RecipeContainer) []*model.Recipe {
	var result []*model.Recipe

	for _, recipe := range recipes {
		result = append(result, mapRecipe(recipe))
	}

	return result
}
func mapRecipe(recipe *recipe.RecipeContainer) *model.Recipe {
	if recipe == nil {
		return nil
	}
	return &model.Recipe{
		ID:               recipe.ID.String(),
		CreatedAt:        recipe.CreatedAt,
		AuthorID:         recipe.UserID.String(),
		CurrentVersionID: recipe.CurrentVersionID.String(),
		CurrentVersion:   mapRecipeVersion(recipe.CurrentVersion),
	}
}

func mapRecipeVersion(recipeVersion *recipe.RecipeVersion) *model.RecipeVersion {
	if recipeVersion == nil {
		return nil
	}
	return &model.RecipeVersion{
		ID:        recipeVersion.ID.String(),
		RecipeID:  recipeVersion.RecipeID.String(),
		Name:      recipeVersion.Name,
		PrepMins:  int32(recipeVersion.PrepMins),
		CookMins:  int32(recipeVersion.CookMins),
		Portions:  int32(recipeVersion.Portions),
		CreatedAt: recipeVersion.CreatedAt,
		Version:   int32(recipeVersion.Version),
	}
}

func toCreateRecipeRequest(input *model.CreateRecipeInput, userID string) (recipe.CreateRecipeRequest, error) {
	if input == nil {
		return recipe.CreateRecipeRequest{}, fmt.Errorf("recipe details are required")
	}
	if input.RecipeSource == nil {
		return recipe.CreateRecipeRequest{}, fmt.Errorf("recipe source is required")
	}

	recipeSourceRequest := recipe.CreateRecipeSourceRequest{
		Type:         int(input.RecipeSource.Type),
		URL:          input.RecipeSource.URL,
		BookTitle:    input.RecipeSource.BookTitle,
		BookPage:     input.RecipeSource.BookPage,
		Instructions: input.RecipeSource.Instructions,
	}

	return recipe.CreateRecipeRequest{
		Name:        input.Name,
		Ingredients: toIngredientUsageRequests(input.IngredientUsages),
		UserID:      userID,
		PrepMins:    int(input.PrepMins),
		CookMins:    int(input.CookMins),
		Portions:    int(input.Portions),
		Source:      recipeSourceRequest,
	}, nil
}

func toIngredientUsageRequests(usages []*model.CreateIngredientUsageInput) []recipe.CreateIngredientUsageRequest {
	ingredientUsageRequests := make([]recipe.CreateIngredientUsageRequest, len(usages))
	for i, usage := range usages {
		ingredientUsageRequests[i] = recipe.CreateIngredientUsageRequest{
			IngredientID: usage.IngredientID,
			Quantity:     usage.Quantity,
			Unit:         int(usage.Unit),
		}
	}
	return ingredientUsageRequests
}

func RequireAuth(ctx context.Context) (*auth.Claims, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, errors.NewUnauthenticatedError("user is not authenticated")
	}
	return claims, nil
}

func GetIPAddress(ctx context.Context) (string, error) {
	ip, ok := ctx.Value(middleware.IPKey).(string)
	if !ok {
		return "", errors.NewInternalError("failed to retrieve IP address from context")
	}
	return ip, nil
}
