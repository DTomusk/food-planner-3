package upload

import "errors"

var (
	ErrInvalidOwnerUserID    = errors.New("owner user ID is required")
	ErrInvalidFileName       = errors.New("file name is required")
	ErrInvalidFileType       = errors.New("file type is required")
	ErrInvalidPurpose        = errors.New("upload purpose is invalid")
	ErrUnsupportedFileType   = errors.New("unsupported file type")
	ErrInvalidFileSize       = errors.New("file size must be greater than zero")
	ErrFileTooLarge          = errors.New("file size exceeds maximum allowed")
	ErrProviderNotConfigured = errors.New("upload provider is not configured")
	ErrMissingUploadExpiry   = errors.New("upload provider did not return an expiry time for the upload URL")
)
