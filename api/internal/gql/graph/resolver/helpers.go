package resolver

import (
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
		Name:      recipeVersion.Name,
		PrepMins:  int32(recipeVersion.PrepMins),
		CookMins:  int32(recipeVersion.CookMins),
		Portions:  int32(recipeVersion.Portions),
		CreatedAt: recipeVersion.CreatedAt,
		Version:   int32(recipeVersion.Version),
	}
}
