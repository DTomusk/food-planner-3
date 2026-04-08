package upload

import (
	"context"
	"fmt"
	"foodplanner/internal/db"
	"foodplanner/internal/logging"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TODO: move to config
const (
	defaultMaxImageSizeBytes             = 5 * 1024 * 1024
	defaultDeleteExpiredUploadsBatchSize = 500
)

type UploadService struct {
	db                db.DBTX
	allowedImageTypes map[string]struct{}
	maxImageSizeBytes int64
	provider          UploadProvider
	repo              *uploadRepo
}

func NewUploadServiceWithProvider(db db.DBTX, provider UploadProvider, maxImageSizeBytes int64, repo *uploadRepo) *UploadService {
	if maxImageSizeBytes <= 0 {
		maxImageSizeBytes = defaultMaxImageSizeBytes
	}
	return &UploadService{
		db: db,
		allowedImageTypes: map[string]struct{}{
			"image/jpeg": {},
			"image/png":  {},
		},
		maxImageSizeBytes: maxImageSizeBytes,
		provider:          provider,
		repo:              repo,
	}
}

func (s *UploadService) CreateImageUploadURL(ctx context.Context, req CreateImageUploadURLRequest) (*CreateImageUploadURLResponse, error) {
	if err := s.validateCreateImageUploadURLRequest(req); err != nil {
		return nil, err
	}

	if s.provider == nil {
		return nil, ErrProviderNotConfigured
	}

	uploadID := uuid.New()
	objectKey := buildObjectKey(req.OwnerUserID, uploadID, req.FileName, req.Purpose)

	providerResponse, err := s.provider.CreateSignedUploadURL(ctx, CreateSignedUploadURLRequest{
		ObjectKey:     objectKey,
		FileType:      req.FileType,
		FileSizeBytes: req.FileSizeBytes,
	})
	if err != nil {
		return nil, err
	}
	if providerResponse.ExpiresAt.IsZero() {
		return nil, ErrMissingUploadExpiry
	}

	upload := &Upload{
		ID:            uploadID,
		OwnerUserID:   req.OwnerUserID,
		ObjectKey:     objectKey,
		FileName:      req.FileName,
		FileType:      req.FileType,
		FileSizeBytes: req.FileSizeBytes,
		Purpose:       req.Purpose,
		ExpiresAt:     providerResponse.ExpiresAt,
	}

	// persist upload metadata
	if err := s.repo.saveUploadMetadata(ctx, s.db, upload); err != nil {
		return nil, err
	}

	return &CreateImageUploadURLResponse{
		UploadID:  uploadID,
		ObjectKey: objectKey,
		UploadURL: providerResponse.UploadURL,
		FileURL:   providerResponse.FileURL,
	}, nil
}

func (s *UploadService) ValidateAndGetFileURL(ctx context.Context, req ValidateAndGetFileURLRequest) (*string, error) {
	if err := s.validateValidateAndGetFileURLRequest(req); err != nil {
		return nil, err
	}

	if s.provider == nil {
		return nil, ErrProviderNotConfigured
	}

	uploadRecord, err := s.repo.getUploadByID(ctx, s.db, req.UploadID)
	if err != nil {
		return nil, err
	}
	if uploadRecord == nil {
		return nil, ErrUploadNotFound
	}
	if uploadRecord.OwnerUserID != req.OwnerUserID {
		return nil, ErrUploadOwnershipMismatch
	}
	if uploadRecord.Purpose != req.Purpose {
		return nil, ErrUploadPurposeMismatch
	}
	if uploadRecord.Used {
		return nil, ErrUploadAlreadyUsed
	}
	if !uploadRecord.ExpiresAt.After(time.Now().UTC()) {
		return nil, ErrUploadExpired
	}

	fileURL := s.provider.FileURLForObjectKey(uploadRecord.ObjectKey)
	return &fileURL, nil
}

func (s *UploadService) MarkUploadAsUsed(ctx context.Context, tx db.DBTX, req ClaimUploadRequest) error {
	if req.UploadID == uuid.Nil {
		return ErrInvalidUploadID
	}

	if req.OwnerUserID == uuid.Nil {
		return ErrInvalidOwnerUserID
	}

	if req.EntityID == uuid.Nil {
		return ErrInvalidEntityID
	}

	if strings.TrimSpace(req.EntityType) == "" {
		return ErrInvalidEntityType
	}

	uploadRecord, err := s.repo.getUploadByID(ctx, tx, req.UploadID)
	if err != nil {
		return err
	}
	if uploadRecord == nil {
		return ErrUploadNotFound
	}
	if uploadRecord.OwnerUserID != req.OwnerUserID {
		return ErrUploadOwnershipMismatch
	}
	if uploadRecord.Used {
		return ErrUploadAlreadyUsed
	}

	return s.repo.markUploadAsUsed(ctx, tx, req.UploadID, req.EntityID, req.EntityType)
}

func (s *UploadService) validateValidateAndGetFileURLRequest(req ValidateAndGetFileURLRequest) error {
	if req.UploadID == uuid.Nil {
		return ErrInvalidUploadID
	}

	if req.OwnerUserID == uuid.Nil {
		return ErrInvalidOwnerUserID
	}

	if !req.Purpose.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidPurpose, req.Purpose)
	}

	return nil
}

func (s *UploadService) validateCreateImageUploadURLRequest(req CreateImageUploadURLRequest) error {
	if req.OwnerUserID == uuid.Nil {
		return ErrInvalidOwnerUserID
	}

	if !req.Purpose.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidPurpose, req.Purpose)
	}

	if strings.TrimSpace(req.FileName) == "" {
		return ErrInvalidFileName
	}

	if strings.TrimSpace(req.FileType) == "" {
		return ErrInvalidFileType
	}

	if _, ok := s.allowedImageTypes[req.FileType]; !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedFileType, req.FileType)
	}

	if req.FileSizeBytes <= 0 {
		return ErrInvalidFileSize
	}

	if req.FileSizeBytes > s.maxImageSizeBytes {
		return ErrFileTooLarge
	}

	return nil
}

func buildObjectKey(userID, uploadID uuid.UUID, fileName string, purpose UploadPurpose) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(fileName)))
	if ext == "" {
		ext = ".bin"
	}

	return fmt.Sprintf("%s/%s/%s%s", purpose, userID.String(), uploadID.String(), ext)
}

func (s *UploadService) DeleteExpiredUploads(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Info("Starting deletion of expired uploads")
	if s.provider == nil {
		return ErrProviderNotConfigured
	}

	for {
		// Get expired uploads that aren't used
		keys, err := s.repo.getExpiredUnusedUploadObjectKeys(ctx, s.db, defaultDeleteExpiredUploadsBatchSize)
		if err != nil {
			return fmt.Errorf("query expired uploads: %w", err)
		}
		if len(keys) == 0 {
			logger.Info("No expired uploads found for deletion, exiting")
			break
		}

		// Send keys to R2 (or other provider) for deletion
		failed, deleteErr := s.provider.DeleteObjects(ctx, keys)
		if failed == nil {
			failed = map[string]error{}
		}

		failedCount := 0
		for _, keyErr := range failed {
			if keyErr != nil {
				failedCount++
			}
		}

		// If provider returned a top-level error and no per-key failures,
		// treat the whole batch as failed.
		if deleteErr != nil && failedCount == 0 {
			for _, key := range keys {
				if _, exists := failed[key]; !exists {
					failed[key] = deleteErr
				}
			}
			failedCount = len(keys)
		}

		if failedCount > 0 && deleteErr == nil {
			deleteErr = fmt.Errorf("provider reported %d failed deletions", failedCount)
		}

		// Delete records from DB for successfully deleted objects
		var successfullyDeleted []string
		for _, key := range keys {
			if keyErr, exists := failed[key]; !exists || keyErr == nil {
				successfullyDeleted = append(successfullyDeleted, key)
			}
		}

		// Set deleted at for successfully deleted uploads in DB
		if len(successfullyDeleted) > 0 {
			if err := s.repo.deleteUploadsByObjectKeys(ctx, s.db, successfullyDeleted); err != nil {
				return fmt.Errorf("delete upload records from DB: %w", err)
			}
			logger.Info("successfully deleted objects", "deleted", len(successfullyDeleted))
		}

		if deleteErr != nil {
			return fmt.Errorf("delete objects from provider: %w", deleteErr)
		}
	}
	return nil
}
