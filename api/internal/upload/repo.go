package upload

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"

	"github.com/google/uuid"
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

const selectUploadByIDQuery = `
SELECT
	id,
	owner_user_id,
	object_key,
	file_name,
	file_type,
	file_size_bytes,
	purpose,
	expires_at,
	used_at IS NOT NULL
FROM uploads
WHERE id = $1
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

func (r *uploadRepo) getUploadByID(ctx context.Context, database db.DBTX, uploadID uuid.UUID) (*Upload, error) {
	var upload Upload
	var purpose string

	err := database.QueryRowContext(ctx, selectUploadByIDQuery, uploadID).Scan(
		&upload.ID,
		&upload.OwnerUserID,
		&upload.ObjectKey,
		&upload.FileName,
		&upload.FileType,
		&upload.FileSizeBytes,
		&purpose,
		&upload.ExpiresAt,
		&upload.Used,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	upload.Purpose = UploadPurpose(purpose)
	return &upload, nil
}
