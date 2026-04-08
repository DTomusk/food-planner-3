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

func TestGetUploadByID_ReturnsUpload(t *testing.T) {
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

		actualUpload, err := repo.getUploadByID(ctx, tx, expectedUpload.ID)
		require.NoError(t, err)
		require.NotNil(t, actualUpload)
		require.Equal(t, expectedUpload.ID, actualUpload.ID)
		require.Equal(t, expectedUpload.OwnerUserID, actualUpload.OwnerUserID)
		require.Equal(t, expectedUpload.ObjectKey, actualUpload.ObjectKey)
		require.Equal(t, expectedUpload.FileName, actualUpload.FileName)
		require.Equal(t, expectedUpload.FileType, actualUpload.FileType)
		require.Equal(t, expectedUpload.FileSizeBytes, actualUpload.FileSizeBytes)
		require.Equal(t, expectedUpload.Purpose, actualUpload.Purpose)
		require.WithinDuration(t, expectedUpload.ExpiresAt, actualUpload.ExpiresAt, time.Second)
		require.False(t, actualUpload.Used)
	})
}

func TestGetUploadByID_ReturnsNilWhenUploadDoesNotExist(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewUploadRepo()

		uploadRecord, err := repo.getUploadByID(context.Background(), tx, uuid.New())
		require.NoError(t, err)
		require.Nil(t, uploadRecord)
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

func TestMarkUploadAsUsedSuccess(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		uploadRecord := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, uploadRecord)
		require.NoError(t, err)

		entityID := uuid.New()
		entityType := "recipe-version"
		err = repo.markUploadAsUsed(ctx, tx, uploadRecord.ID, entityID, entityType)
		require.NoError(t, err)

		// Verify upload is marked as used with entity linkage
		var usedAt sql.NullTime
		var linkedEntityID sql.NullString
		var linkedEntityType sql.NullString

		err = tx.QueryRowContext(
			ctx,
			`SELECT used_at, linked_entity_id, linked_entity_type FROM uploads WHERE id = $1`,
			uploadRecord.ID,
		).Scan(&usedAt, &linkedEntityID, &linkedEntityType)
		require.NoError(t, err)
		require.True(t, usedAt.Valid, "used_at should be set")
		require.True(t, linkedEntityID.Valid, "linked_entity_id should be set")
		require.Equal(t, entityID.String(), linkedEntityID.String)
		require.True(t, linkedEntityType.Valid, "linked_entity_type should be set")
		require.Equal(t, entityType, linkedEntityType.String)
	})
}

func TestMarkUploadAsUsedDetectsRaceCondition(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		uploadRecord := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, uploadRecord)
		require.NoError(t, err)

		// Mark as used first time
		firstEntityID := uuid.New()
		err = repo.markUploadAsUsed(ctx, tx, uploadRecord.ID, firstEntityID, "recipe-version")
		require.NoError(t, err)

		// Try to mark as used again - should fail due to RowsAffected check
		secondEntityID := uuid.New()
		err = repo.markUploadAsUsed(ctx, tx, uploadRecord.ID, secondEntityID, "recipe-version")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadAlreadyUsed)
	})
}

func TestGetUploadByID_ReturnsNilWhenUploadSoftDeleted(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		uploadRecord := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, uploadRecord)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx, `UPDATE uploads SET deleted_at = NOW() WHERE id = $1`, uploadRecord.ID)
		require.NoError(t, err)

		found, err := repo.getUploadByID(ctx, tx, uploadRecord.ID)
		require.NoError(t, err)
		require.Nil(t, found)
	})
}

func TestGetExpiredUnusedUploadObjectKeys_ReturnsExpiredUnusedWithLimit(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		expiredOlder := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "expired-older.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}
		expiredNewer := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "expired-newer.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}
		usedExpired := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "used-expired.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}
		deletedExpired := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "deleted-expired.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}
		active := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "active.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(2 * time.Hour),
		}

		uploads := []*Upload{expiredOlder, expiredNewer, usedExpired, deletedExpired, active}
		for _, u := range uploads {
			err = repo.saveUploadMetadata(ctx, tx, u)
			require.NoError(t, err)
		}

		_, err = tx.ExecContext(ctx, `UPDATE uploads SET created_at = NOW() - INTERVAL '6 hours', expires_at = NOW() - INTERVAL '5 hours' WHERE id = $1`, expiredOlder.ID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE uploads SET created_at = NOW() - INTERVAL '4 hours', expires_at = NOW() - INTERVAL '3 hours' WHERE id = $1`, expiredNewer.ID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE uploads SET created_at = NOW() - INTERVAL '4 hours', expires_at = NOW() - INTERVAL '3 hours', used_at = NOW() WHERE id = $1`, usedExpired.ID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE uploads SET created_at = NOW() - INTERVAL '4 hours', expires_at = NOW() - INTERVAL '3 hours', deleted_at = NOW() WHERE id = $1`, deletedExpired.ID)
		require.NoError(t, err)

		keys, err := repo.getExpiredUnusedUploadObjectKeys(ctx, tx, 1)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, expiredOlder.ObjectKey, keys[0])

		keys, err = repo.getExpiredUnusedUploadObjectKeys(ctx, tx, 10)
		require.NoError(t, err)
		require.Equal(t, []string{expiredOlder.ObjectKey, expiredNewer.ObjectKey}, keys)
	})
}

func TestDeleteUploadsByObjectKeys_SoftDeletesMatchingRows(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		uploadA := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "a.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}
		uploadB := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "b.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}
		uploadC := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "c.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		uploads := []*Upload{uploadA, uploadB, uploadC}
		for _, u := range uploads {
			err = repo.saveUploadMetadata(ctx, tx, u)
			require.NoError(t, err)
		}

		err = repo.deleteUploadsByObjectKeys(ctx, tx, []string{uploadA.ObjectKey, uploadB.ObjectKey})
		require.NoError(t, err)

		var deletedA sql.NullTime
		var deletedB sql.NullTime
		var deletedC sql.NullTime

		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, uploadA.ID).Scan(&deletedA)
		require.NoError(t, err)
		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, uploadB.ID).Scan(&deletedB)
		require.NoError(t, err)
		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, uploadC.ID).Scan(&deletedC)
		require.NoError(t, err)

		require.True(t, deletedA.Valid)
		require.True(t, deletedB.Valid)
		require.False(t, deletedC.Valid)
	})
}
