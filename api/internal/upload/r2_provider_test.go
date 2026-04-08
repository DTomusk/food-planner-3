package upload

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResolveR2EndpointURL(t *testing.T) {
	t.Parallel()

	t.Run("uses explicit endpoint URL when provided", func(t *testing.T) {
		t.Parallel()

		endpoint, err := resolveR2EndpointURL(R2UploadProviderConfig{
			EndpointURL: "https://custom.endpoint.example.com",
			AccountID:   "ignored-account",
		})
		require.NoError(t, err)
		require.Equal(t, "https://custom.endpoint.example.com", endpoint)
	})

	t.Run("builds endpoint URL from account ID", func(t *testing.T) {
		t.Parallel()

		endpoint, err := resolveR2EndpointURL(R2UploadProviderConfig{AccountID: "abc123"})
		require.NoError(t, err)
		require.Equal(t, "https://abc123.r2.cloudflarestorage.com", endpoint)
	})

	t.Run("returns error when endpoint and account ID are missing", func(t *testing.T) {
		t.Parallel()

		endpoint, err := resolveR2EndpointURL(R2UploadProviderConfig{})
		require.Error(t, err)
		require.Empty(t, endpoint)
		require.ErrorContains(t, err, "account ID or endpoint URL is required")
	})
}

func TestNewR2UploadProviderValidatesRequiredConfig(t *testing.T) {
	t.Parallel()

	base := validR2ProviderConfig()
	tests := []struct {
		name        string
		mutate      func(*R2UploadProviderConfig)
		errorSubstr string
	}{
		{
			name: "requires bucket name",
			mutate: func(cfg *R2UploadProviderConfig) {
				cfg.BucketName = ""
			},
			errorSubstr: "bucket name is required",
		},
		{
			name: "requires access key ID",
			mutate: func(cfg *R2UploadProviderConfig) {
				cfg.AccessKeyID = ""
			},
			errorSubstr: "access key ID is required",
		},
		{
			name: "requires secret access key",
			mutate: func(cfg *R2UploadProviderConfig) {
				cfg.SecretAccessKey = ""
			},
			errorSubstr: "secret access key is required",
		},
		{
			name: "requires public base URL",
			mutate: func(cfg *R2UploadProviderConfig) {
				cfg.PublicBaseURL = ""
			},
			errorSubstr: "public base URL is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := base
			tc.mutate(&cfg)

			provider, err := NewR2UploadProvider(context.Background(), cfg)
			require.Error(t, err)
			require.Nil(t, provider)
			require.ErrorContains(t, err, tc.errorSubstr)
		})
	}
}

func TestNewR2UploadProviderAppliesDefaults(t *testing.T) {
	t.Parallel()

	cfg := validR2ProviderConfig()
	cfg.Region = ""
	cfg.PresignExpiry = 0

	provider, err := NewR2UploadProvider(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, provider)
	require.Equal(t, defaultR2PresignExpiry, provider.presignExpiry)
}

func TestR2UploadProviderCreateSignedUploadURL(t *testing.T) {
	t.Parallel()

	t.Run("returns provider not configured for nil provider", func(t *testing.T) {
		t.Parallel()

		var provider *R2UploadProvider
		res, err := provider.CreateSignedUploadURL(context.Background(), CreateSignedUploadURLRequest{
			ObjectKey: "recipe-images/user/upload.png",
			FileType:  "image/png",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrProviderNotConfigured)
		require.Nil(t, res)
	})

	t.Run("validates object key", func(t *testing.T) {
		t.Parallel()

		provider, err := NewR2UploadProvider(context.Background(), validR2ProviderConfig())
		require.NoError(t, err)

		res, err := provider.CreateSignedUploadURL(context.Background(), CreateSignedUploadURLRequest{
			ObjectKey: "  ",
			FileType:  "image/png",
		})
		require.Error(t, err)
		require.Nil(t, res)
		require.ErrorContains(t, err, "object key is required")
	})

	t.Run("validates file type", func(t *testing.T) {
		t.Parallel()

		provider, err := NewR2UploadProvider(context.Background(), validR2ProviderConfig())
		require.NoError(t, err)

		res, err := provider.CreateSignedUploadURL(context.Background(), CreateSignedUploadURLRequest{
			ObjectKey: "recipe-images/user/upload.png",
			FileType:  "",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidFileType)
		require.Nil(t, res)
	})

	t.Run("returns presigned upload URL and public file URL", func(t *testing.T) {
		t.Parallel()

		cfg := validR2ProviderConfig()
		provider, err := NewR2UploadProvider(context.Background(), cfg)
		require.NoError(t, err)

		objectKey := "recipe-images/user-1/upload-1.png"
		res, err := provider.CreateSignedUploadURL(context.Background(), CreateSignedUploadURLRequest{
			ObjectKey:     objectKey,
			FileType:      "image/png",
			FileSizeBytes: 2048,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEmpty(t, res.UploadURL)
		require.Contains(t, res.UploadURL, "X-Amz-Algorithm=")
		require.True(t, strings.HasPrefix(res.UploadURL, cfg.EndpointURL), "expected signed URL to start with endpoint")
		require.Equal(t, cfg.PublicBaseURL+"/"+objectKey, res.FileURL)
	})
}

func validR2ProviderConfig() R2UploadProviderConfig {
	return R2UploadProviderConfig{
		EndpointURL:     "https://account-id.r2.cloudflarestorage.com",
		BucketName:      "food-smash-assets",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		PublicBaseURL:   "https://images.example.com",
		Region:          defaultR2Region,
		PresignExpiry:   2 * time.Minute,
	}
}
