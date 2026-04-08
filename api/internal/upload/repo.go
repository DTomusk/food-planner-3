package upload

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
AND deleted_at IS NULL
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

func (r *uploadRepo) markUploadAsUsed(ctx context.Context, database db.DBTX, uploadID uuid.UUID, entityID uuid.UUID, entityType string) error {
	query := `UPDATE uploads
SET used_at = NOW(), linked_entity_id = $2, linked_entity_type = $3
WHERE id = $1 AND (used_at IS NULL)`
	result, err := database.ExecContext(ctx, query, uploadID, entityID, entityType)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUploadAlreadyUsed
	}

	return nil
}

func (r *uploadRepo) getExpiredUnusedUploadObjectKeys(ctx context.Context, database db.DBTX, batchSize int) ([]string, error) {
	if batchSize <= 0 {
		batchSize = 500
	}

	// Get expired (including grace period) and unused uploads.
	query := `SELECT object_key FROM uploads
	WHERE used_at IS NULL
	AND deleted_at IS NULL
	AND expires_at < NOW() - INTERVAL '1 hour'
	ORDER BY expires_at ASC
	LIMIT $1`
	rows, err := database.QueryContext(ctx, query, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objectKeys []string
	for rows.Next() {
		var objectKey string
		if err := rows.Scan(&objectKey); err != nil {
			return nil, err
		}
		objectKeys = append(objectKeys, objectKey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objectKeys, nil
}

func (r *uploadRepo) deleteUploadsByObjectKeys(ctx context.Context, database db.DBTX, objectKeys []string) error {
	if len(objectKeys) == 0 {
		return nil
	}

	// Mark uploads as deleted by setting deleted_at (soft delete)
	query := `UPDATE uploads SET deleted_at = NOW() WHERE object_key = ANY($1)`
	_, err := database.ExecContext(ctx, query, pq.Array(objectKeys))
	return err
}
