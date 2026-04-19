package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"foodplanner/internal/auth"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/config"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
	"foodplanner/internal/gql/graph"
	"foodplanner/internal/gql/graph/directive"
	"foodplanner/internal/gql/graph/resolver"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/middleware"
	"foodplanner/internal/recipe"
	"foodplanner/internal/upload"
	"foodplanner/internal/user"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	redisPublisher := events.NewRedisPublisher(redisClient, cfg.RedisStream, cfg.RedisStreamMaxLen)

	ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), cfg.IngredientUpsertBatchSize)

	recipeRepo, err := recipe.NewRecipeRepo(cfg.RecipeSearchTrigramWeight, cfg.RecipeSearchFullTextWeight)
	if err != nil {
		log.Fatalf("Failed to create recipe repository: %v", err)
	}

	userService := user.NewUserService(txRunner.DB(), user.NewUserRepo())
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationMinutes)
	refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), cfg.RefreshTokenSecret, cfg.RefreshTokenExpirationDays)
	authService := auth.NewAuthService(txRunner.DB(), userService, jwtService, refreshTokenService, redisPublisher)

	uploadProvider, err := upload.NewR2UploadProvider(ctx, upload.R2UploadProviderConfig{
		AccountID:       cfg.R2AccountID,
		EndpointURL:     cfg.R2EndpointURL,
		BucketName:      cfg.R2BucketName,
		AccessKeyID:     cfg.R2AccessKeyID,
		SecretAccessKey: cfg.R2SecretAccessKey,
		PublicBaseURL:   cfg.R2PublicBaseURL,
		Region:          cfg.R2Region,
		PresignExpiry:   time.Duration(cfg.R2PresignExpiry) * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create R2 upload provider: %v", err)
	}
	uploadService := upload.NewUploadServiceWithProvider(txRunner.DB(), uploadProvider, cfg.UploadMaxImageSizeBytes, upload.NewUploadRepo())

	recipeService := recipe.NewService(
		txRunner,
		recipeRepo,
		recipe.NewRecipeVersionRepo(),
		ingredientService,
		recipe.NewIngredientUsageRepo(),
		nil,
		uploadService,
		redisPublisher,
	)

	complexity := graph.NewComplexityRoot()

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &resolver.Resolver{
					AuthService:        authService,
					RecipeService:      recipeService,
					IngredientsService: ingredientService,
					UploadService:      uploadService,
					UserService:        userService,
				},
				Directives: graph.DirectiveRoot{
					Auth: directive.AuthDirective,
				},
				Complexity: complexity,
			},
		),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	// Set custom error presenter to check for complexity limit errors
	// Used for auditing
	srv.SetErrorPresenter(graph.NewComplexityLimitErrorPresenter(redisPublisher, graph.DefaultMaxAcceptedComplexity))

	srv.Use(extension.Introspection{})
	// Enforce complexity limit to prevent expensive queries
	srv.Use(extension.FixedComplexityLimit(graph.DefaultMaxAcceptedComplexity))
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	authMiddleware := auth.Middleware(jwtService)
	ipMiddleware := middleware.IPMiddleware
	userAgentMiddleware := middleware.UserAgentMiddleware
	rateLimitingMiddleware := middleware.NewRateLimitingMiddleware(redisClient, middleware.RateLimitingConfig{
		Window:               cfg.RateLimitingWindow,
		AnonymousLimit:       cfg.RateLimitingAnonymousLimit,
		AuthenticatedLimit:   cfg.RateLimitingAuthenticatedLimit,
		FailOpenOnRedisError: cfg.RateLimitingFailOpenOnRedisError,
	})
	responseWriterMiddleware := middleware.ResponseWriterMiddleware
	requestMiddleware := middleware.RequestMiddleware

	http.Handle("/query",
		ipMiddleware(
			userAgentMiddleware(
				authMiddleware(
					rateLimitingMiddleware(
						responseWriterMiddleware(
							requestMiddleware(srv),
						),
					),
				),
			),
		),
	)

	// Check API health (process is running)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
		})
	})
	// Check API dependencies, currently database and redis
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		readyCtx, cancelReady := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancelReady()

		response := map[string]any{
			"status":     "ready",
			"database":   "ok",
			"redis":      "ok",
			"checked_at": time.Now().UTC(),
		}

		if err := database.PingContext(readyCtx); err != nil {
			response["status"] = "not_ready"
			response["database"] = "error"
			response["database_error"] = err.Error()
		}
		if err := redisClient.Ping(readyCtx).Err(); err != nil {
			response["status"] = "not_ready"
			response["redis"] = "error"
			response["redis_error"] = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		if response["status"] != "ready" {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		_ = json.NewEncoder(w).Encode(response)
	})
	http.HandleFunc("/metrics/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		snapshot := events.SnapshotMetrics()
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			log.Printf("failed to encode metrics snapshot: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

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

	httpHandler := c.Handler(http.DefaultServeMux)
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: httpHandler,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP shutdown failed: %v", err)
	}

	log.Println("HTTP server stopped cleanly")
}
