package main

import (
	"database/sql"
	"foodplanner/internal/config"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/reference"
	"foodplanner/internal/sync"
	"log"
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

	ingredientService := ingredient.NewIngredientService(database, ingredient.NewIngredientRepo())
	referenceService := reference.NewReferenceService()
	syncService := sync.NewSyncService(ingredientService, referenceService)

	err = syncService.SyncIngredientData()
	if err != nil {
		log.Fatalf("Failed to sync ingredient data: %v", err)
	}
}
