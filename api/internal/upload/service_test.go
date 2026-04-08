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

func TestUploadServiceCreateImageUploadURLSuccess(t *testing.T) {
	t.Parallel()

	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		req := CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "dish.PNG",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
		}

		res, err := service.CreateImageUploadURL(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEqual(t, uuid.Nil, res.UploadID)

		expectedObjectKey := fmt.Sprintf("%s/%s/%s%s", req.Purpose, req.OwnerUserID.String(), res.UploadID.String(), ".png")
		require.Equal(t, expectedObjectKey, res.ObjectKey)
		require.Equal(t, "https://upload.example.com/"+expectedObjectKey, res.UploadURL)
		require.Equal(t, "https://cdn.example.com/"+expectedObjectKey, res.FileURL)
	})
}

func TestUploadServiceCreateImageUploadURLUsesBinExtensionWhenMissing(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		req := CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "avatar",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
		}

		res, err := service.CreateImageUploadURL(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Contains(t, res.ObjectKey, ".bin")
	})
}

func TestUploadServiceCreateImageUploadURLValidationErrors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 100, NewUploadRepo())

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		baseReq := CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "image.png",
			FileType:      "image/png",
			FileSizeBytes: 10,
			Purpose:       UploadPurposeRecipeImage,
		}

		tests := []struct {
			name    string
			mutate  func(*CreateImageUploadURLRequest)
			wantErr error
		}{
			{
				name: "missing owner user ID",
				mutate: func(req *CreateImageUploadURLRequest) {
					req.OwnerUserID = uuid.Nil
				},
				wantErr: ErrInvalidOwnerUserID,
			},
			{
				name: "invalid purpose",
				mutate: func(req *CreateImageUploadURLRequest) {
					req.Purpose = UploadPurpose("unexpected-purpose")
				},
				wantErr: ErrInvalidPurpose,
			},
			{
				name: "missing file name",
				mutate: func(req *CreateImageUploadURLRequest) {
					req.FileName = "  "
				},
				wantErr: ErrInvalidFileName,
			},
			{
				name: "missing file type",
				mutate: func(req *CreateImageUploadURLRequest) {
					req.FileType = " "
				},
				wantErr: ErrInvalidFileType,
			},
			{
				name: "unsupported file type",
				mutate: func(req *CreateImageUploadURLRequest) {
					req.FileType = "image/gif"
				},
				wantErr: ErrUnsupportedFileType,
			},
			{
				name: "invalid file size",
				mutate: func(req *CreateImageUploadURLRequest) {
					req.FileSizeBytes = 0
				},
				wantErr: ErrInvalidFileSize,
			},
			{
				name: "file too large",
				mutate: func(req *CreateImageUploadURLRequest) {
					req.FileSizeBytes = 101
				},
				wantErr: ErrFileTooLarge,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := baseReq
				tc.mutate(&req)

				res, err := service.CreateImageUploadURL(ctx, req)
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)
				require.Nil(t, res)
			})
		}
	})
}

func TestUploadServiceCreateImageUploadURLProviderNotConfigured(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		service := NewUploadServiceWithProvider(tx, nil, defaultMaxImageSizeBytes, NewUploadRepo())

		req := CreateImageUploadURLRequest{
			OwnerUserID:   uuid.New(),
			FileName:      "image.png",
			FileType:      "image/png",
			FileSizeBytes: 100,
			Purpose:       UploadPurposeRecipeImage,
		}

		res, err := service.CreateImageUploadURL(context.Background(), req)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrProviderNotConfigured)
		require.Nil(t, res)
	})
}

func TestNewUploadServiceWithProviderUsesDefaultMaxSize(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 0, NewUploadRepo())

		req := CreateImageUploadURLRequest{
			OwnerUserID:   uuid.New(),
			FileName:      "image.png",
			FileType:      "image/png",
			FileSizeBytes: defaultMaxImageSizeBytes + 1,
			Purpose:       UploadPurposeRecipeImage,
		}

		res, err := service.CreateImageUploadURL(context.Background(), req)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFileTooLarge)
		require.Nil(t, res)
	})
}

func TestUploadServiceValidateAndGetFileURLSuccess(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		repo := NewUploadRepo()
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, repo)

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

		fileURL, err := service.ValidateAndGetFileURL(ctx, ValidateAndGetFileURLRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: testUser.ID,
			Purpose:     UploadPurposeRecipeImage,
		})
		require.NoError(t, err)
		require.NotNil(t, fileURL)
		require.Equal(t, "https://cdn.example.com/"+uploadRecord.ObjectKey, *fileURL)
	})
}

func TestUploadServiceValidateAndGetFileURLValidationErrors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		validUserID := uuid.New()
		validUploadID := uuid.New()

		tests := []struct {
			name    string
			request ValidateAndGetFileURLRequest
			wantErr error
		}{
			{
				name: "missing upload ID",
				request: ValidateAndGetFileURLRequest{
					OwnerUserID: validUserID,
					Purpose:     UploadPurposeRecipeImage,
				},
				wantErr: ErrInvalidUploadID,
			},
			{
				name: "missing owner user ID",
				request: ValidateAndGetFileURLRequest{
					UploadID: validUploadID,
					Purpose:  UploadPurposeRecipeImage,
				},
				wantErr: ErrInvalidOwnerUserID,
			},
			{
				name: "invalid purpose",
				request: ValidateAndGetFileURLRequest{
					UploadID:    validUploadID,
					OwnerUserID: validUserID,
					Purpose:     UploadPurpose("bad-purpose"),
				},
				wantErr: ErrInvalidPurpose,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				fileURL, err := service.ValidateAndGetFileURL(ctx, tc.request)
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)
				require.Nil(t, fileURL)
			})
		}
	})
}

func TestUploadServiceValidateAndGetFileURLReturnsNotFoundWhenUploadDoesNotExist(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		fileURL, err := service.ValidateAndGetFileURL(ctx, ValidateAndGetFileURLRequest{
			UploadID:    uuid.New(),
			OwnerUserID: uuid.New(),
			Purpose:     UploadPurposeRecipeImage,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadNotFound)
		require.Nil(t, fileURL)
	})
}

func TestUploadServiceValidateAndGetFileURLReturnsErrorWhenOwnershipDoesNotMatch(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		repo := NewUploadRepo()
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, repo)

		owner, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)
		otherUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		uploadRecord := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   owner.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, owner.ID.String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, uploadRecord)
		require.NoError(t, err)

		fileURL, err := service.ValidateAndGetFileURL(ctx, ValidateAndGetFileURLRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: otherUser.ID,
			Purpose:     UploadPurposeRecipeImage,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadOwnershipMismatch)
		require.Nil(t, fileURL)
	})
}

func TestUploadServiceValidateAndGetFileURLReturnsErrorWhenUploadAlreadyUsed(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		repo := NewUploadRepo()
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, repo)

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

		_, err = tx.ExecContext(ctx, `UPDATE uploads SET used_at = NOW() WHERE id = $1`, uploadRecord.ID)
		require.NoError(t, err)

		fileURL, err := service.ValidateAndGetFileURL(ctx, ValidateAndGetFileURLRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: testUser.ID,
			Purpose:     UploadPurposeRecipeImage,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadAlreadyUsed)
		require.Nil(t, fileURL)
	})
}

func TestUploadServiceValidateAndGetFileURLReturnsErrorWhenUploadExpired(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		repo := NewUploadRepo()
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, repo)

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

		_, err = tx.ExecContext(
			ctx,
			`UPDATE uploads
			 SET created_at = NOW() - INTERVAL '2 minutes',
			     expires_at = NOW() - INTERVAL '1 minute'
			 WHERE id = $1`,
			uploadRecord.ID,
		)
		require.NoError(t, err)

		fileURL, err := service.ValidateAndGetFileURL(ctx, ValidateAndGetFileURLRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: testUser.ID,
			Purpose:     UploadPurposeRecipeImage,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadExpired)
		require.Nil(t, fileURL)
	})
}

func TestUploadServiceValidateAndGetFileURLProviderNotConfigured(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		service := NewUploadServiceWithProvider(tx, nil, defaultMaxImageSizeBytes, NewUploadRepo())

		fileURL, err := service.ValidateAndGetFileURL(context.Background(), ValidateAndGetFileURLRequest{
			UploadID:    uuid.New(),
			OwnerUserID: uuid.New(),
			Purpose:     UploadPurposeRecipeImage,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrProviderNotConfigured)
		require.Nil(t, fileURL)
	})
}

func TestUploadServiceMarkUploadAsUsedSuccess(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		repo := NewUploadRepo()
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, repo)

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
		claimReq := ClaimUploadRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: testUser.ID,
			EntityID:    entityID,
			EntityType:  "recipe-version",
		}

		err = service.MarkUploadAsUsed(ctx, tx, claimReq)
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
		require.Equal(t, "recipe-version", linkedEntityType.String)
	})
}

func TestUploadServiceMarkUploadAsUsedValidationErrors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		validEntityID := uuid.New()
		validUploadID := uuid.New()

		tests := []struct {
			name    string
			request ClaimUploadRequest
			wantErr error
		}{
			{
				name: "missing upload ID",
				request: ClaimUploadRequest{
					OwnerUserID: testUser.ID,
					EntityID:    validEntityID,
					EntityType:  "recipe-version",
				},
				wantErr: ErrInvalidUploadID,
			},
			{
				name: "missing owner user ID",
				request: ClaimUploadRequest{
					UploadID:   validUploadID,
					EntityID:   validEntityID,
					EntityType: "recipe-version",
				},
				wantErr: ErrInvalidOwnerUserID,
			},
			{
				name: "missing entity ID",
				request: ClaimUploadRequest{
					UploadID:    validUploadID,
					OwnerUserID: testUser.ID,
					EntityType:  "recipe-version",
				},
				wantErr: ErrInvalidEntityID,
			},
			{
				name: "missing entity type",
				request: ClaimUploadRequest{
					UploadID:    validUploadID,
					OwnerUserID: testUser.ID,
					EntityID:    validEntityID,
				},
				wantErr: ErrInvalidEntityType,
			},
			{
				name: "empty entity type",
				request: ClaimUploadRequest{
					UploadID:    validUploadID,
					OwnerUserID: testUser.ID,
					EntityID:    validEntityID,
					EntityType:  "  ",
				},
				wantErr: ErrInvalidEntityType,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := service.MarkUploadAsUsed(ctx, tx, tc.request)
				require.Error(t, err)
				require.ErrorIs(t, err, tc.wantErr)
			})
		}
	})
}

func TestUploadServiceMarkUploadAsUsedReturnsNotFoundWhenUploadDoesNotExist(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		err = service.MarkUploadAsUsed(ctx, tx, ClaimUploadRequest{
			UploadID:    uuid.New(),
			OwnerUserID: testUser.ID,
			EntityID:    uuid.New(),
			EntityType:  "recipe-version",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadNotFound)
	})
}

func TestUploadServiceMarkUploadAsUsedReturnsErrorWhenOwnershipDoesNotMatch(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		repo := NewUploadRepo()
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, repo)

		owner, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)
		otherUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		uploadRecord := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   owner.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, owner.ID.String(), uuid.New().String()),
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(10 * time.Minute),
		}

		err = repo.saveUploadMetadata(ctx, tx, uploadRecord)
		require.NoError(t, err)

		err = service.MarkUploadAsUsed(ctx, tx, ClaimUploadRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: otherUser.ID,
			EntityID:    uuid.New(),
			EntityType:  "recipe-version",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadOwnershipMismatch)
	})
}

func TestUploadServiceMarkUploadAsUsedReturnsErrorWhenUploadAlreadyUsed(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		repo := NewUploadRepo()
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, repo)

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

		// Mark as used by another entity first
		firstEntityID := uuid.New()
		err = service.MarkUploadAsUsed(ctx, tx, ClaimUploadRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: testUser.ID,
			EntityID:    firstEntityID,
			EntityType:  "recipe-version",
		})
		require.NoError(t, err)

		// Try to claim again with different entity
		err = service.MarkUploadAsUsed(ctx, tx, ClaimUploadRequest{
			UploadID:    uploadRecord.ID,
			OwnerUserID: testUser.ID,
			EntityID:    uuid.New(),
			EntityType:  "recipe-version",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUploadAlreadyUsed)
	})
}

type cleanupProviderStub struct {
	failed map[string]error
	err    error
}

func (p *cleanupProviderStub) CreateSignedUploadURL(_ context.Context, _ CreateSignedUploadURLRequest) (*CreateSignedUploadURLResponse, error) {
	return nil, ErrProviderNotConfigured
}

func (p *cleanupProviderStub) FileURLForObjectKey(objectKey string) string {
	return "https://cdn.example.com/" + objectKey
}

func (p *cleanupProviderStub) DeleteObjects(_ context.Context, _ []string) (map[string]error, error) {
	return p.failed, p.err
}

func TestUploadServiceDeleteExpiredUploadsProviderNotConfigured(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		service := NewUploadServiceWithProvider(tx, nil, defaultMaxImageSizeBytes, NewUploadRepo())

		err := service.DeleteExpiredUploads(context.Background())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrProviderNotConfigured)
	})
}

func TestUploadServiceDeleteExpiredUploadsSoftDeletesExpiredUnusedUploads(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, defaultMaxImageSizeBytes, repo)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		expiredA := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "expired-a.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(2 * time.Hour),
		}
		expiredB := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "expired-b.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(2 * time.Hour),
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

		err = repo.saveUploadMetadata(ctx, tx, expiredA)
		require.NoError(t, err)
		err = repo.saveUploadMetadata(ctx, tx, expiredB)
		require.NoError(t, err)
		err = repo.saveUploadMetadata(ctx, tx, active)
		require.NoError(t, err)

		_, err = tx.ExecContext(
			ctx,
			`UPDATE uploads
			 SET created_at = NOW() - INTERVAL '3 hours',
			     expires_at = NOW() - INTERVAL '2 hours'
			 WHERE id = $1 OR id = $2`,
			expiredA.ID,
			expiredB.ID,
		)
		require.NoError(t, err)

		err = service.DeleteExpiredUploads(ctx)
		require.NoError(t, err)

		var expiredADeletedAt sql.NullTime
		var expiredBDeletedAt sql.NullTime
		var activeDeletedAt sql.NullTime

		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, expiredA.ID).Scan(&expiredADeletedAt)
		require.NoError(t, err)
		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, expiredB.ID).Scan(&expiredBDeletedAt)
		require.NoError(t, err)
		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, active.ID).Scan(&activeDeletedAt)
		require.NoError(t, err)

		require.True(t, expiredADeletedAt.Valid)
		require.True(t, expiredBDeletedAt.Valid)
		require.False(t, activeDeletedAt.Valid)
	})
}

func TestUploadServiceDeleteExpiredUploadsSoftDeletesSuccessfulKeysWhenProviderPartiallyFails(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewUploadRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		expiredA := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "expired-a.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(2 * time.Hour),
		}
		expiredB := &Upload{
			ID:            uuid.New(),
			OwnerUserID:   testUser.ID,
			ObjectKey:     fmt.Sprintf("%s/%s/%s.png", UploadPurposeRecipeImage, testUser.ID.String(), uuid.New().String()),
			FileName:      "expired-b.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
			ExpiresAt:     time.Now().UTC().Add(2 * time.Hour),
		}

		err = repo.saveUploadMetadata(ctx, tx, expiredA)
		require.NoError(t, err)
		err = repo.saveUploadMetadata(ctx, tx, expiredB)
		require.NoError(t, err)

		_, err = tx.ExecContext(
			ctx,
			`UPDATE uploads
			 SET created_at = NOW() - INTERVAL '3 hours',
			     expires_at = NOW() - INTERVAL '2 hours'
			 WHERE id = $1 OR id = $2`,
			expiredA.ID,
			expiredB.ID,
		)
		require.NoError(t, err)

		provider := &cleanupProviderStub{
			failed: map[string]error{
				expiredB.ObjectKey: fmt.Errorf("simulated delete failure"),
			},
			err: fmt.Errorf("failed to delete 1 objects from R2"),
		}
		service := NewUploadServiceWithProvider(tx, provider, defaultMaxImageSizeBytes, repo)

		err = service.DeleteExpiredUploads(ctx)
		require.Error(t, err)

		var expiredADeletedAt sql.NullTime
		var expiredBDeletedAt sql.NullTime

		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, expiredA.ID).Scan(&expiredADeletedAt)
		require.NoError(t, err)
		err = tx.QueryRowContext(ctx, `SELECT deleted_at FROM uploads WHERE id = $1`, expiredB.ID).Scan(&expiredBDeletedAt)
		require.NoError(t, err)

		require.True(t, expiredADeletedAt.Valid)
		require.False(t, expiredBDeletedAt.Valid)
	})
}
