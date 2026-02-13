package main

import (
	"context"
	"database/sql"
	"foodplanner/internal/config"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/reference"
	"foodplanner/internal/sync"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	ingredientService := ingredient.NewIngredientService(database, ingredient.NewIngredientRepo())
	referenceService := reference.NewReferenceService()
	syncService := sync.NewSyncService(ingredientService, referenceService)

	if err := syncService.SyncIngredientData(ctx); err != nil {
		log.Fatalf("Failed to sync ingredient data: %v", err)
	}
}
