package upload

import "github.com/google/uuid"

type UploadPurpose string

const (
	UploadPurposeRecipeImage UploadPurpose = "recipe-images"
)

var validUploadPurposes = map[UploadPurpose]struct{}{
	UploadPurposeRecipeImage: {},
}

func (p UploadPurpose) IsValid() bool {
	_, ok := validUploadPurposes[p]
	return ok
}

type CreateImageUploadURLRequest struct {
	OwnerUserID   uuid.UUID
	FileName      string
	FileType      string
	FileSizeBytes int64
	Purpose       UploadPurpose
}

type CreateImageUploadURLResponse struct {
	UploadID  uuid.UUID
	ObjectKey string
	UploadURL string
	FileURL   string
}

type ValidateAndGetFileURLRequest struct {
	UploadID    uuid.UUID
	OwnerUserID uuid.UUID
	Purpose     UploadPurpose
}
