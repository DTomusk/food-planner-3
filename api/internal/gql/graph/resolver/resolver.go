package resolver

import (
	"foodplanner/internal/auth"
	"foodplanner/internal/recipe"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	AuthService   *auth.AuthService
	RecipeService *recipe.Service
}
