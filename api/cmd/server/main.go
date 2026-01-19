package main

import (
	"database/sql"
	"food-planner/internal/auth"
	"food-planner/internal/config"
	"food-planner/internal/gql/graph/generated"
	"food-planner/internal/gql/graph/resolver"
	"food-planner/internal/recipe"
	"food-planner/internal/user"
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

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	log.Printf("Starting server on port %s", cfg.ServerPort)

	recipeRepo := recipe.NewRepo()
	recipeService := recipe.NewService(db, recipeRepo)

	userRepo := user.NewUserRepo()
	userService := user.NewUserService(db, userRepo)
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationMinutes)
	authService := auth.NewAuthService(db, userService, jwtService)

	srv := handler.New(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolver.Resolver{
					AuthService:   authService,
					RecipeService: recipeService,
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
	http.Handle("/query", srv)

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
