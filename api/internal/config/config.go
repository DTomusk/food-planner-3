package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	ServerPort        string
	CorsAllowedOrigin string
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

	return &Config{
		DatabaseURL:       db_url,
		ServerPort:        port,
		CorsAllowedOrigin: corsOrigin,
	}, nil
}
