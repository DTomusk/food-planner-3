package main

import (
	"context"
	"database/sql"
	"foodplanner/internal/config"
	"foodplanner/internal/db"
	"foodplanner/internal/upload"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	txRunner := db.NewDBTxRunner(database)

	uploadProvider, err := upload.NewR2UploadProvider(ctx, upload.R2UploadProviderConfig{
		AccountID:       cfg.R2AccountID,
		EndpointURL:     cfg.R2EndpointURL,
		BucketName:      cfg.R2BucketName,
		AccessKeyID:     cfg.R2AccessKeyID,
		SecretAccessKey: cfg.R2SecretAccessKey,
		PublicBaseURL:   cfg.R2PublicBaseURL,
		Region:          cfg.R2Region,
		PresignExpiry:   time.Duration(cfg.R2PresignExpiry) * time.Second,
	})

	if err != nil {
		log.Fatalf("Failed to create R2 upload provider: %v", err)
	}

	uploadService := upload.NewUploadServiceWithProvider(txRunner.DB(), uploadProvider, cfg.UploadMaxImageSizeBytes, upload.NewUploadRepo())

	if err := uploadService.DeleteExpiredUploads(ctx); err != nil {
		log.Fatalf("Failed to delete expired uploads: %v", err)
	}
}
