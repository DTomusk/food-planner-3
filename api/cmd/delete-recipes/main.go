package main

import (
	"context"
	"database/sql"
	"foodplanner/internal/config"
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

	// txRunner := db.NewDBTxRunner(database)

	// ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), cfg.IngredientUpsertBatchSize)
	// recipeService := recipe.NewService(txRunner, recipe.NewRepo(), ingredientService, &cfg.RecipeRetentionDays)

	log.Println("Starting recipe deletion...")

	// deletedCount, err := recipeService.DeleteOldRecipes(ctx)
	// if err != nil {
	// 	log.Fatalf("Failed to delete old recipes: %v", err)
	// }
	// log.Printf("Successfully deleted %d old recipes", deletedCount)
}
