package upload

import "github.com/google/uuid"

type CreateImageUploadURLRequest struct {
	OwnerUserID   uuid.UUID
	FileName      string
	FileType      string
	FileSizeBytes int64
}

type CreateImageUploadURLResponse struct {
	UploadID  uuid.UUID
	ObjectKey string
	UploadURL string
	FileURL   string
}
