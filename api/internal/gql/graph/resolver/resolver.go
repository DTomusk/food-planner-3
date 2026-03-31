package resolver

import (
	"foodplanner/internal/auth"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/recipe"
	"foodplanner/internal/upload"
	"foodplanner/internal/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	AuthService        *auth.AuthService
	IngredientsService *ingredient.IngredientService
	RecipeService      *recipe.Service
	UploadService      *upload.UploadService
	UserService        *user.UserService
}
