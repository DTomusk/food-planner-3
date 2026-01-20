package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL          string
	ServerPort           string
	CorsAllowedOrigin    string
	JWTSecret            string
	JWTExpirationMinutes int
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

	return &Config{
		DatabaseURL:          db_url,
		ServerPort:           port,
		CorsAllowedOrigin:    corsOrigin,
		JWTSecret:            jwtSecret,
		JWTExpirationMinutes: jwtExpirationMinutes,
	}, nil
}
