package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

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
	RecipeSearchTrigramWeight  float64
	RecipeSearchFullTextWeight float64
	// R2
	UploadMaxImageSizeBytes int64
	R2AccountID             string
	R2EndpointURL           string
	R2BucketName            string
	R2AccessKeyID           string
	R2SecretAccessKey       string
	R2PublicBaseURL         string
	R2Region                string
	R2PresignExpiry         int
	// Redis
	RedisAddress      string
	RedisPassword     string
	RedisDB           int
	RedisStream       string
	RedisStreamMaxLen int64
	// Rate limiting
	RateLimitingWindow               time.Duration
	RateLimitingAnonymousLimit       int
	RateLimitingAuthenticatedLimit   int
	RateLimitingFailOpenOnRedisError bool
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

	recipeSearchTrigramWeightStr := os.Getenv("RECIPE_SEARCH_TRIGRAM_WEIGHT")
	if recipeSearchTrigramWeightStr == "" {
		return nil, fmt.Errorf("RECIPE_SEARCH_TRIGRAM_WEIGHT not set in environment")
	}
	recipeSearchTrigramWeight, err := strconv.ParseFloat(recipeSearchTrigramWeightStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid RECIPE_SEARCH_TRIGRAM_WEIGHT: %v", err)
	}

	recipeSearchFullTextWeightStr := os.Getenv("RECIPE_SEARCH_FULL_TEXT_WEIGHT")
	if recipeSearchFullTextWeightStr == "" {
		return nil, fmt.Errorf("RECIPE_SEARCH_FULL_TEXT_WEIGHT not set in environment")
	}
	recipeSearchFullTextWeight, err := strconv.ParseFloat(recipeSearchFullTextWeightStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid RECIPE_SEARCH_FULL_TEXT_WEIGHT: %v", err)
	}

	uploadMaxImageSizeBytesStr := os.Getenv("UPLOAD_MAX_IMAGE_SIZE_BYTES")
	if uploadMaxImageSizeBytesStr == "" {
		return nil, fmt.Errorf("UPLOAD_MAX_IMAGE_SIZE_BYTES not set in environment")
	}
	uploadMaxImageSizeBytes, err := strconv.ParseInt(uploadMaxImageSizeBytesStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid UPLOAD_MAX_IMAGE_SIZE_BYTES: %v", err)
	}

	r2AccountID := os.Getenv("R2_ACCOUNT_ID")
	if r2AccountID == "" {
		return nil, fmt.Errorf("R2_ACCOUNT_ID not set in environment")
	}

	r2EndpointURL := os.Getenv("R2_ENDPOINT_URL")
	if r2EndpointURL == "" {
		return nil, fmt.Errorf("R2_ENDPOINT_URL not set in environment")
	}

	r2BucketName := os.Getenv("R2_BUCKET_NAME")
	if r2BucketName == "" {
		return nil, fmt.Errorf("R2_BUCKET_NAME not set in environment")
	}

	r2AccessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	if r2AccessKeyID == "" {
		return nil, fmt.Errorf("R2_ACCESS_KEY_ID not set in environment")
	}

	r2SecretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	if r2SecretAccessKey == "" {
		return nil, fmt.Errorf("R2_SECRET_ACCESS_KEY not set in environment")
	}

	r2PublicBaseURL := os.Getenv("R2_PUBLIC_BASE_URL")
	if r2PublicBaseURL == "" {
		return nil, fmt.Errorf("R2_PUBLIC_BASE_URL not set in environment")
	}

	r2Region := os.Getenv("R2_REGION")
	if r2Region == "" {
		return nil, fmt.Errorf("R2_REGION not set in environment")
	}

	r2PresignExpiryStr := os.Getenv("R2_PRESIGN_EXPIRY")
	if r2PresignExpiryStr == "" {
		return nil, fmt.Errorf("R2_PRESIGN_EXPIRY not set in environment")
	}
	r2PresignExpiry, err := strconv.Atoi(r2PresignExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid R2_PRESIGN_EXPIRY: %v", err)
	}

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		return nil, fmt.Errorf("REDIS_ADDRESS not set in environment")
	}

	// REDIS_PASSWORD is optional — Redis may run without auth
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// REDIS_DB is optional — defaults to 0
	redisDB := 0
	if redisDBStr := os.Getenv("REDIS_DB"); redisDBStr != "" {
		redisDB, err = strconv.Atoi(redisDBStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB: %v", err)
		}
	}

	redisStream := os.Getenv("REDIS_STREAM")
	if redisStream == "" {
		return nil, fmt.Errorf("REDIS_STREAM not set in environment")
	}

	redisStreamMaxLenStr := os.Getenv("REDIS_STREAM_MAX_LEN")
	if redisStreamMaxLenStr == "" {
		return nil, fmt.Errorf("REDIS_STREAM_MAX_LEN not set in environment")
	}
	redisStreamMaxLen, err := strconv.ParseInt(redisStreamMaxLenStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_STREAM_MAX_LEN: %v", err)
	}

	rateLimitingWindowStr := os.Getenv("RATE_LIMITING_WINDOW")
	if rateLimitingWindowStr == "" {
		return nil, fmt.Errorf("RATE_LIMITING_WINDOW not set in environment")
	}
	rateLimitingWindow, err := time.ParseDuration(rateLimitingWindowStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMITING_WINDOW: %v", err)
	}
	rateLimitingAnonymousLimitStr := os.Getenv("RATE_LIMITING_ANONYMOUS_LIMIT")
	if rateLimitingAnonymousLimitStr == "" {
		return nil, fmt.Errorf("RATE_LIMITING_ANONYMOUS_LIMIT not set in environment")
	}
	rateLimitingAnonymousLimit, err := strconv.Atoi(rateLimitingAnonymousLimitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMITING_ANONYMOUS_LIMIT: %v", err)
	}
	rateLimitingAuthenticatedLimitStr := os.Getenv("RATE_LIMITING_AUTHENTICATED_LIMIT")
	if rateLimitingAuthenticatedLimitStr == "" {
		return nil, fmt.Errorf("RATE_LIMITING_AUTHENTICATED_LIMIT not set in environment")
	}
	rateLimitingAuthenticatedLimit, err := strconv.Atoi(rateLimitingAuthenticatedLimitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMITING_AUTHENTICATED_LIMIT: %v", err)
	}
	rateLimitingFailOpenOnRedisErrorStr := os.Getenv("RATE_LIMITING_FAIL_OPEN_ON_REDIS_ERROR")
	if rateLimitingFailOpenOnRedisErrorStr == "" {
		return nil, fmt.Errorf("RATE_LIMITING_FAIL_OPEN_ON_REDIS_ERROR not set in environment")
	}
	rateLimitingFailOpenOnRedisError, err := strconv.ParseBool(rateLimitingFailOpenOnRedisErrorStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMITING_FAIL_OPEN_ON_REDIS_ERROR: %v", err)
	}

	return &Config{
		DatabaseURL:                      db_url,
		ServerPort:                       port,
		CorsAllowedOrigin:                corsOrigin,
		JWTSecret:                        jwtSecret,
		JWTExpirationMinutes:             jwtExpirationMinutes,
		IngredientDataFilePath:           ingredientDataFilePath,
		IngredientUpsertBatchSize:        ingredientUpsertBatchSize,
		RecipeRetentionDays:              recipeRetentionDays,
		RefreshTokenSecret:               refreshTokenSecret,
		RefreshTokenExpirationDays:       refreshTokenExpirationDays,
		RecipeSearchTrigramWeight:        recipeSearchTrigramWeight,
		RecipeSearchFullTextWeight:       recipeSearchFullTextWeight,
		UploadMaxImageSizeBytes:          uploadMaxImageSizeBytes,
		R2AccountID:                      r2AccountID,
		R2EndpointURL:                    r2EndpointURL,
		R2BucketName:                     r2BucketName,
		R2AccessKeyID:                    r2AccessKeyID,
		R2SecretAccessKey:                r2SecretAccessKey,
		R2PublicBaseURL:                  r2PublicBaseURL,
		R2Region:                         r2Region,
		R2PresignExpiry:                  r2PresignExpiry,
		RedisAddress:                     redisAddress,
		RedisPassword:                    redisPassword,
		RedisDB:                          redisDB,
		RedisStream:                      redisStream,
		RedisStreamMaxLen:                redisStreamMaxLen,
		RateLimitingWindow:               rateLimitingWindow,
		RateLimitingAnonymousLimit:       rateLimitingAnonymousLimit,
		RateLimitingAuthenticatedLimit:   rateLimitingAuthenticatedLimit,
		RateLimitingFailOpenOnRedisError: rateLimitingFailOpenOnRedisError,
	}, nil
}
