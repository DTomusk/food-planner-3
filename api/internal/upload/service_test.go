package upload

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUploadServiceCreateImageUploadURLSuccess(t *testing.T) {
	t.Parallel()

	testutil.WithTx(t, func(tx *sql.Tx) {
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		req := CreateImageUploadURLRequest{
			OwnerUserID:   uuid.New(),
			FileName:      "dish.PNG",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
		}

		res, err := service.CreateImageUploadURL(context.Background(), req)
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
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 10*1024*1024, NewUploadRepo())

		req := CreateImageUploadURLRequest{
			OwnerUserID:   uuid.New(),
			FileName:      "avatar",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       UploadPurposeRecipeImage,
		}

		res, err := service.CreateImageUploadURL(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Contains(t, res.ObjectKey, ".bin")
	})
}

func TestUploadServiceCreateImageUploadURLValidationErrors(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		provider := NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
		service := NewUploadServiceWithProvider(tx, provider, 100, NewUploadRepo())

		baseReq := CreateImageUploadURLRequest{
			OwnerUserID:   uuid.New(),
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

				res, err := service.CreateImageUploadURL(context.Background(), req)
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
