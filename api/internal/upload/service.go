package upload

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// TODO: move to config
const (
	defaultMaxImageSizeBytes = 5 * 1024 * 1024
)

type UploadService struct {
	allowedImageTypes map[string]struct{}
	maxImageSizeBytes int64
	provider          UploadProvider
}

func NewUploadServiceWithProvider(provider UploadProvider, maxImageSizeBytes int64) *UploadService {
	if maxImageSizeBytes <= 0 {
		maxImageSizeBytes = defaultMaxImageSizeBytes
	}
	return &UploadService{
		allowedImageTypes: map[string]struct{}{
			"image/jpeg": {},
			"image/png":  {},
		},
		maxImageSizeBytes: maxImageSizeBytes,
		provider:          provider,
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
	objectKey := buildObjectKey(req.OwnerUserID, uploadID, req.FileName)

	providerResponse, err := s.provider.CreateSignedUploadURL(ctx, CreateSignedUploadURLRequest{
		ObjectKey:     objectKey,
		FileType:      req.FileType,
		FileSizeBytes: req.FileSizeBytes,
	})
	if err != nil {
		return nil, err
	}

	return &CreateImageUploadURLResponse{
		UploadID:  uploadID,
		ObjectKey: objectKey,
		UploadURL: providerResponse.UploadURL,
		FileURL:   providerResponse.FileURL,
	}, nil
}

func (s *UploadService) validateCreateImageUploadURLRequest(req CreateImageUploadURLRequest) error {
	if req.OwnerUserID == uuid.Nil {
		return ErrInvalidOwnerUserID
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

func buildObjectKey(userID, uploadID uuid.UUID, fileName string) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(fileName)))
	if ext == "" {
		ext = ".bin"
	}

	return fmt.Sprintf("recipe-images/%s/%s%s", userID.String(), uploadID.String(), ext)
}
