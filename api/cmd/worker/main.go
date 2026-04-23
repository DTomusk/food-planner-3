package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"foodplanner/internal/audit"
	"foodplanner/internal/config"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
		TLSConfig: func() *tls.Config {
			if !cfg.RedisUseTLS {
				return nil
			}
			return &tls.Config{MinVersion: tls.VersionTLS12}
		}(),
	})

	// Register event type strings and their corresponding event structs
	registry := events.NewRegistry()
	registry.Register(events.GraphQLRequestRejectedType, func() events.Event { return &events.GraphQLRequestRejectedEvent{} })
	registry.Register(events.RateLimitExceededType, func() events.Event { return &events.RateLimitExceededEvent{} })
	registry.Register(events.RecipeCreatedEventType, func() events.Event { return &events.RecipeCreatedEvent{} })
	registry.Register(events.RecipeUpdatedEventType, func() events.Event { return &events.RecipeUpdatedEvent{} })
	registry.Register(events.UserSignedUpType, func() events.Event { return &events.UserSignedUpEvent{} })
	registry.Register(events.UserSignedInType, func() events.Event { return &events.UserSignedInEvent{} })
	registry.Register(events.UserSigninFailedType, func() events.Event { return &events.UserSigninFailedEvent{} })
	expectedVersions := map[string]int{
		events.GraphQLRequestRejectedType: 1,
		events.RateLimitExceededType:      1,
		events.RecipeCreatedEventType:     1,
		events.RecipeUpdatedEventType:     1,
		events.UserSignedUpType:           1,
		events.UserSignedInType:           1,
		events.UserSigninFailedType:       1,
	}
	eventsRepo := events.NewEventsRepo()

	worker := events.NewRedisWorker(
		redisClient,
		cfg.RedisStream,
		"worker-group",
		"worker-1",
		registry,
		expectedVersions,
		eventsRepo,
		txRunner,
	)

	auditService := audit.NewAuditService(audit.NewRepo())
	auditProjector := audit.NewAuditProjector(auditService, audit.DefaultAuditMappers())

	// Register a single projector for all event types that map to audit entries.
	for _, eventType := range []string{
		events.RecipeCreatedEventType,
		events.RecipeUpdatedEventType,
		events.UserSignedUpType,
		events.UserSignedInType,
		events.UserSigninFailedType,
		events.GraphQLRequestRejectedType,
		events.RateLimitExceededType,
	} {
		if err := worker.RegisterHandler(eventType, auditProjector); err != nil {
			log.Fatalf("Failed to register audit projector for %s: %v", eventType, err)
		}
	}

	// Run worker
	if err := worker.Run(ctx); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}
