package resolver

import (
	"context"
	"fmt"
	"foodplanner/internal/auth"
	"foodplanner/internal/gql/graph/errors"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/middleware"
	"foodplanner/internal/recipe"
	"foodplanner/internal/user"
	"net/http"

	"github.com/google/uuid"
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
		ID:          recipeVersion.ID.String(),
		RecipeID:    recipeVersion.RecipeID.String(),
		Name:        recipeVersion.Name,
		Description: recipeVersion.Description,
		PrepMins:    int32(recipeVersion.PrepMins),
		CookMins:    int32(recipeVersion.CookMins),
		Portions:    int32(recipeVersion.Portions),
		CreatedAt:   recipeVersion.CreatedAt,
		Version:     int32(recipeVersion.Version),
		ImgSrc:      recipeVersion.ImgSrc,
	}
}

func buildRecipeListParams(
	pagination *model.PaginationInput,
	filter *model.RecipeFilterInput,
) (recipe.RecipeListParams, error) {
	first := 20
	var after *string
	if pagination != nil {
		if pagination.First != nil {
			first = int(*pagination.First)
		}
		after = pagination.After
	}

	var query *string
	var userID *uuid.UUID
	if filter != nil {
		query = filter.Query
		if filter.UserID != nil {
			parsed, err := uuid.Parse(*filter.UserID)
			if err != nil {
				return recipe.RecipeListParams{}, fmt.Errorf("invalid userID in recipe filter: %w", err)
			}
			userID = &parsed
		}
	}

	return recipe.RecipeListParams{
		Pagination: recipe.RecipePagination{
			First: first,
			After: after,
		},
		Filter: recipe.RecipeFilter{
			Query:  query,
			UserID: userID,
		},
	}, nil
}

func buildRecipeConnection(
	recipes []*recipe.RecipeWithCursor,
	endCursor *string,
) (*model.RecipeConnection, error) {
	if recipes == nil {
		return &model.RecipeConnection{
			Edges:    []*model.RecipeEdge{},
			PageInfo: &model.PageInfo{HasNextPage: false, EndCursor: nil},
		}, nil
	}

	connection := &model.RecipeConnection{
		Edges:    make([]*model.RecipeEdge, len(recipes)),
		PageInfo: &model.PageInfo{HasNextPage: endCursor != nil, EndCursor: endCursor},
	}
	for i, recipeWithCursor := range recipes {
		if recipeWithCursor == nil || recipeWithCursor.Recipe == nil {
			return nil, fmt.Errorf("recipe service returned invalid recipe page item")
		}
		connection.Edges[i] = &model.RecipeEdge{
			Cursor: recipeWithCursor.Cursor,
			Node:   mapRecipe(recipeWithCursor.Recipe),
		}
	}

	return connection, nil
}

func (r *Resolver) listRecipes(
	ctx context.Context,
	pagination *model.PaginationInput,
	filter *model.RecipeFilterInput,
) (*model.RecipeConnection, error) {
	params, err := buildRecipeListParams(pagination, filter)
	if err != nil {
		return nil, err
	}

	recipes, endCursor, err := r.RecipeService.GetRecipes(ctx, params)
	if err != nil {
		return nil, err
	}

	return buildRecipeConnection(recipes, endCursor)
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

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	return recipe.CreateRecipeRequest{
		Name:        input.Name,
		Description: description,
		Ingredients: toIngredientUsageRequests(input.IngredientUsages),
		UserID:      userID,
		PrepMins:    int(input.PrepMins),
		CookMins:    int(input.CookMins),
		Portions:    int(input.Portions),
		Source:      recipeSourceRequest,
		ImgUploadID: input.ImgUploadID,
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

func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	w, _ := ctx.Value(middleware.ResponseWriterKey).(http.ResponseWriter)
	return w
}

func GetRequest(ctx context.Context) *http.Request {
	req, _ := ctx.Value(middleware.RequestKey).(*http.Request)
	return req
}
