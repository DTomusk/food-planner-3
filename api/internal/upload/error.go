package upload

import "errors"

var (
	ErrInvalidOwnerUserID    = errors.New("owner user ID is required")
	ErrInvalidFileName       = errors.New("file name is required")
	ErrInvalidFileType       = errors.New("file type is required")
	ErrUnsupportedFileType   = errors.New("unsupported file type")
	ErrInvalidFileSize       = errors.New("file size must be greater than zero")
	ErrFileTooLarge          = errors.New("file size exceeds maximum allowed")
	ErrProviderNotConfigured = errors.New("upload provider is not configured")
)
