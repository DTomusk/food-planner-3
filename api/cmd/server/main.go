package main

import (
	"database/sql"
	"foodplanner/internal/auth"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/config"
	"foodplanner/internal/db"
	"foodplanner/internal/gql/graph"
	"foodplanner/internal/gql/graph/directive"
	"foodplanner/internal/gql/graph/resolver"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/middleware"
	"foodplanner/internal/recipe"
	"foodplanner/internal/user"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	log.Printf("Starting server on port %s", cfg.ServerPort)

	txRunner := db.NewDBTxRunner(database)

	ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), cfg.IngredientUpsertBatchSize)

	recipeService := recipe.NewService(
		txRunner,
		recipe.NewRecipeRepo(),
		recipe.NewRecipeVersionRepo(),
		ingredientService,
		recipe.NewIngredientUsageRepo(),
		nil,
	)

	userService := user.NewUserService(txRunner.DB(), user.NewUserRepo())
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationMinutes)
	refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), cfg.RefreshTokenSecret, cfg.RefreshTokenExpirationDays)
	authService := auth.NewAuthService(txRunner.DB(), userService, jwtService, refreshTokenService)

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &resolver.Resolver{
					AuthService:        authService,
					RecipeService:      recipeService,
					IngredientsService: ingredientService,
					UserService:        userService,
				},
				Directives: graph.DirectiveRoot{
					Auth: directive.AuthDirective,
				},
			},
		),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	authMiddleware := auth.Middleware(jwtService)
	ipMiddleware := middleware.IPMiddleware
	responseWriterMiddleware := middleware.ResponseWriterMiddleware
	requestMiddleware := middleware.RequestMiddleware

	http.Handle("/query", ipMiddleware(authMiddleware(responseWriterMiddleware(requestMiddleware(srv)))))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			cfg.CorsAllowedOrigin,
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	})

	handler := c.Handler(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, handler))
}
