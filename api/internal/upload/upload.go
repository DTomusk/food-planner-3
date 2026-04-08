package upload

import (
	"time"

	"github.com/google/uuid"
)

type Upload struct {
	ObjectKey     string
	URL           string
	Used          bool
	OwnerUserID   uuid.UUID
	ID            uuid.UUID
	FileName      string
	FileType      string
	FileSizeBytes int64
	Purpose       UploadPurpose
	ExpiresAt     time.Time
}
