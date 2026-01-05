package main

import (
	"context"
	"database/sql"
	"food-planner/internal/config"
	"food-planner/internal/recipe"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Database URL: %s", cfg.DatabaseURL)

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	repo := recipe.NewRepo()

	log.Println("Creating a new recipe...")

	id := uuid.New()

	err = repo.CreateRecipe(recipe.Recipe{
		Name: "Blah",
		ID:   id,
	}, context.Background(), db)

	if err != nil {
		log.Fatalf("Failed to create recipe: %v", err)
	}

	log.Println("Recipe created successfully")

	log.Println("Retrieving the recipe by ID...")

	retRecipe, err := repo.GetRecipeByID(id.String(), context.Background(), db)
	if err != nil {
		log.Fatalf("Failed to retrieve recipe: %v", err)
	}
	log.Printf("Retrieved Recipe: %s", retRecipe.String())
}
