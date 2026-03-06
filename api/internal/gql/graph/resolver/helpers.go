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

func mapRecipes(recipes []*recipe.RecipeVersion) []*model.Recipe {
	var result []*model.Recipe

	for _, recipe := range recipes {
		result = append(result, mapRecipe(recipe))
	}

	return result
}
func mapRecipe(recipe *recipe.RecipeVersion) *model.Recipe {
	if recipe == nil {
		return nil
	}
	return nil
	// return &model.Recipe{
	// 	ID:       recipe.ID.String(),
	// 	Name:     recipe.Name,
	// 	PrepMins: int32(recipe.PrepMins),
	// 	CookMins: int32(recipe.CookMins),
	// 	Portions: int32(recipe.Portions),
	// }
}
