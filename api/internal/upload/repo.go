package upload

import (
	"context"
	"foodplanner/internal/db"
)

type uploadRepo struct{}

const insertUploadQuery = `
INSERT INTO uploads (
	id,
	owner_user_id,
	object_key,
	file_name,
	file_type,
	file_size_bytes,
	purpose,
	expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

func NewUploadRepo() *uploadRepo {
	return &uploadRepo{}
}

func (r *uploadRepo) saveUploadMetadata(ctx context.Context, database db.DBTX, upload *Upload) error {
	_, err := database.ExecContext(
		ctx,
		insertUploadQuery,
		upload.ID,
		upload.OwnerUserID,
		upload.ObjectKey,
		upload.FileName,
		upload.FileType,
		upload.FileSizeBytes,
		upload.Purpose,
		upload.ExpiresAt,
	)
	if err != nil {
		return err
	}

	return nil
}
