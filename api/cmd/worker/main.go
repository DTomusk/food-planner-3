package main

import (
	"context"
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
	})

	// Register event type strings and their corresponding event structs
	registry := events.NewRegistry()
	registry.Register(events.UserSignedUpType, func() events.Event { return &events.UserSignedUpEvent{} })
	expectedVersions := map[string]int{
		events.UserSignedUpType: 1,
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

	// Register handlers to handle events (can register multiple handlers to an event)
	if err := worker.RegisterHandler(events.UserSignedUpType, audit.NewSignupEventHandler(auditService)); err != nil {
		log.Fatalf("Failed to register signup audit handler: %v", err)
	}

	// Run worker
	if err := worker.Run(ctx); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}
