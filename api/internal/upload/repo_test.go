package upload

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveUploadMetadata_InsertsUpload(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		expectedUpload := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, expectedUpload)
		require.NoError(t, err)

		var (
			actualOwnerUserID   uuid.UUID
			actualObjectKey     string
			actualFileName      string
			actualFileType      string
			actualFileSizeBytes int64
			actualPurpose       string
			actualExpiresAt     time.Time
		)

		err = tx.QueryRowContext(
			ctx,
			`SELECT owner_user_id, object_key, file_name, file_type, file_size_bytes, purpose, expires_at FROM uploads WHERE id = $1`,
			expectedUpload.ID,
		).Scan(
			&actualOwnerUserID,
			&actualObjectKey,
			&actualFileName,
			&actualFileType,
			&actualFileSizeBytes,
			&actualPurpose,
			&actualExpiresAt,
		)
		require.NoError(t, err)
		require.Equal(t, expectedUpload.OwnerUserID, actualOwnerUserID)
		require.Equal(t, expectedUpload.ObjectKey, actualObjectKey)
		require.Equal(t, expectedUpload.FileName, actualFileName)
		require.Equal(t, expectedUpload.FileType, actualFileType)
		require.Equal(t, expectedUpload.FileSizeBytes, actualFileSizeBytes)
		require.Equal(t, string(expectedUpload.Purpose), actualPurpose)
		require.WithinDuration(t, expectedUpload.ExpiresAt, actualExpiresAt, time.Second)
	})
}

func TestSaveUploadMetadata_ReturnsErrorWhenOwnerDoesNotExist(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		upload := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   uuid.New(),
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, uuid.New().String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err := repo.saveUploadMetadata(ctx, tx, upload)
		require.Error(t, err)
	})
}

func TestSaveUploadMetadata_ReturnsErrorWhenObjectKeyAlreadyExists(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		duplicateObjectKey := fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String())

		firstUpload := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     duplicateObjectKey,
			FileName:      "first.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, firstUpload)
		require.NoError(t, err)

		secondUpload := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     duplicateObjectKey,
			FileName:      "second.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, secondUpload)
		require.Error(t, err)
	})
}

func TestSaveUploadMetadata_ReturnsErrorWhenPurposeIsInvalid(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		upload := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurpose("invalid-purpose"),
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, upload)
		require.Error(t, err)
	})
}
