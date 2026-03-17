package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL                string
	ServerPort                 string
	CorsAllowedOrigin          string
	JWTSecret                  string
	JWTExpirationMinutes       int
	IngredientDataFilePath     string
	IngredientUpsertBatchSize  int
	RecipeRetentionDays        int
	RefreshTokenSecret         string
	RefreshTokenExpirationDays int
}

func Load() (*Config, error) {
	if os.Getenv("ENV") != "docker" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		return nil, fmt.Errorf("DB_URL not set in environment")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return nil, fmt.Errorf("SERVER_PORT not set in environment")
	}

	corsOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if corsOrigin == "" {
		return nil, fmt.Errorf("CORS_ALLOWED_ORIGIN not set in environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set in environment")
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION_MINUTES")
	if jwtExpirationStr == "" {
		return nil, fmt.Errorf("JWT_EXPIRATION_MINUTES not set in environment")
	}

	jwtExpirationMinutes, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_MINUTES: %v", err)
	}

	ingredientDataFilePath := os.Getenv("INGREDIENT_DATA_FILE_PATH")
	if ingredientDataFilePath == "" {
		return nil, fmt.Errorf("INGREDIENT_DATA_FILE_PATH not set in environment")
	}

	ingredientUpsertBatchSizeStr := os.Getenv("INGREDIENT_UPSERT_BATCH_SIZE")
	if ingredientUpsertBatchSizeStr == "" {
		return nil, fmt.Errorf("INGREDIENT_UPSERT_BATCH_SIZE not set in environment")
	}
	ingredientUpsertBatchSize, err := strconv.Atoi(ingredientUpsertBatchSizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid INGREDIENT_UPSERT_BATCH_SIZE: %v", err)
	}

	recipeRetentionDaysStr := os.Getenv("RECIPE_RETENTION_DAYS")
	if recipeRetentionDaysStr == "" {
		return nil, fmt.Errorf("RECIPE_RETENTION_DAYS not set in environment")
	}
	recipeRetentionDays, err := strconv.Atoi(recipeRetentionDaysStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RECIPE_RETENTION_DAYS: %v", err)
	}

	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		return nil, fmt.Errorf("REFRESH_TOKEN_SECRET not set in environment")
	}

	refreshTokenExpirationDaysStr := os.Getenv("REFRESH_TOKEN_EXPIRATION_DAYS")
	if refreshTokenExpirationDaysStr == "" {
		return nil, fmt.Errorf("REFRESH_TOKEN_EXPIRATION_DAYS not set in environment")
	}
	refreshTokenExpirationDays, err := strconv.Atoi(refreshTokenExpirationDaysStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRATION_DAYS: %v", err)
	}

	return &Config{
		DatabaseURL:                db_url,
		ServerPort:                 port,
		CorsAllowedOrigin:          corsOrigin,
		JWTSecret:                  jwtSecret,
		JWTExpirationMinutes:       jwtExpirationMinutes,
		IngredientDataFilePath:     ingredientDataFilePath,
		IngredientUpsertBatchSize:  ingredientUpsertBatchSize,
		RecipeRetentionDays:        recipeRetentionDays,
		RefreshTokenSecret:         refreshTokenSecret,
		RefreshTokenExpirationDays: refreshTokenExpirationDays,
	}, nil
}
