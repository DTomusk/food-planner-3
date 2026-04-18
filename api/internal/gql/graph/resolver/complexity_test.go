package resolver

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"foodplanner/internal/events"
	"foodplanner/internal/gql/graph"
	"foodplanner/internal/gql/graph/directive"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/middleware"
	"foodplanner/internal/recipe"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

type recordingPublisher struct {
	published []events.Event
}

func (p *recordingPublisher) Publish(ctx context.Context, event events.Event) error {
	p.published = append(p.published, event)
	return nil
}

type graphQLTestResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func newComplexityTestHandler(t *testing.T, tx *sql.Tx, publisher events.Publisher) (http.Handler, *recipe.Service) {
	t.Helper()

	txRunner := testutil.NewTestTxRunner(tx)
	ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
	recipeService := newTestRecipeService(t, tx, txRunner, ingredientService)
	userService := user.NewUserService(tx, user.NewUserRepo())

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &Resolver{
					IngredientsService: ingredientService,
					RecipeService:      recipeService,
					UserService:        userService,
				},
				Directives: graph.DirectiveRoot{
					Auth: directive.AuthDirective,
				},
				Complexity: graph.NewComplexityRoot(),
			},
		),
	)

	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.SetErrorPresenter(graph.NewComplexityLimitErrorPresenter(publisher, graph.DefaultMaxAcceptedComplexity))
	srv.Use(extension.FixedComplexityLimit(graph.DefaultMaxAcceptedComplexity))

	wrapped := middleware.IPMiddleware(middleware.UserAgentMiddleware(middleware.RequestMiddleware(srv)))

	return wrapped, recipeService
}

func seedRecipeForComplexityTest(t *testing.T, tx *sql.Tx, recipeService *recipe.Service) {
	t.Helper()

	ctx := context.Background()
	testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
	require.NoError(t, err)

	testUser, err := seeds.SeedTestUser(ctx, tx)
	require.NoError(t, err)

	_, err = recipeService.CreateRecipe(ctx, recipe.CreateRecipeRequest{
		Name:        "Complexity Test Recipe",
		Description: "Used for GraphQL complexity handler tests",
		Ingredients: []recipe.CreateIngredientUsageRequest{
			{
				IngredientID: testIngredient.ID.String(),
				Quantity:     1,
				Unit:         1,
			},
		},
		UserID:    testUser.ID.String(),
		PrepMins:  10,
		CookMins:  20,
		Portions:  2,
		Source:    recipe.CreateRecipeSourceRequest{Type: 1, URL: ptrString("https://example.com/complexity-test")},
		IPAddress: "127.0.0.1",
		UserAgent: "resolver-complexity-test/1.0",
	})
	require.NoError(t, err)
}

// A deeply nested query (recipes -> versions -> ingredientUsages)
// Should be fine for a couple of results, but a large page will be over the limit
func executeGraphQLRequest(t *testing.T, gqlHandler http.Handler, first int) graphQLTestResponse {
	t.Helper()

	payload := map[string]any{
		"query": `query Recipes($first: Int!) {
		  recipes(pagination: { first: $first }) {
		    edges {
		      node {
		        versions {
		          ingredientUsages {
		            ingredient {
		              id
		            }
		          }
		        }
		      }
		    }
		  }
		}`,
		"variables": map[string]any{
			"first": first,
		},
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "resolver-complexity-test/1.0")
	rr := httptest.NewRecorder()

	gqlHandler.ServeHTTP(rr, req)

	var response graphQLTestResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	return response
}

func TestGraphQLComplexity_AllowsQueryUnderLimit(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		publisher := &recordingPublisher{}
		gqlHandler, recipeService := newComplexityTestHandler(t, tx, publisher)
		seedRecipeForComplexityTest(t, tx, recipeService)

		response := executeGraphQLRequest(t, gqlHandler, 2)

		require.Empty(t, response.Errors)
		require.NotEmpty(t, response.Data)
		require.Contains(t, string(response.Data), `"recipes"`)
		require.Contains(t, string(response.Data), `"edges"`)
		require.Empty(t, publisher.published)
	})
}

func TestGraphQLComplexity_RejectsQueryOverLimit(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		publisher := &recordingPublisher{}
		gqlHandler, recipeService := newComplexityTestHandler(t, tx, publisher)
		seedRecipeForComplexityTest(t, tx, recipeService)

		response := executeGraphQLRequest(t, gqlHandler, 100)

		require.NotEmpty(t, response.Errors)
		require.True(t, strings.Contains(response.Errors[0].Message, "complexity") || strings.Contains(response.Errors[0].Message, "limit"))
		require.Len(t, publisher.published, 1)

		rejectedEvent, ok := publisher.published[0].(events.GraphQLRequestRejectedEvent)
		require.True(t, ok)
		require.Equal(t, events.GraphQLRequestRejectedType, rejectedEvent.Metadata().Type)
		require.Equal(t, "Recipes", rejectedEvent.OperationName)
		require.Equal(t, "query", rejectedEvent.OperationType)
		require.Equal(t, "complexity_limit_exceeded", rejectedEvent.Reason)
		require.Equal(t, "/query", rejectedEvent.Path)
		require.Equal(t, graph.DefaultMaxAcceptedComplexity, rejectedEvent.MaxComplexity)
		require.NotEmpty(t, rejectedEvent.QueryHash)
	})
}
