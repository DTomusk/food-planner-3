package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/auth"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/upload"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func setupUploadMutationResolver(tx *sql.Tx) *mutationResolver {
	provider := upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com")
	service := upload.NewUploadServiceWithProvider(tx, provider, 0, upload.NewUploadRepo())
	r := &Resolver{
		UploadService: service,
	}
	return &mutationResolver{r}
}

func TestUploadResolver_CreateImageUploadURL_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		mutResolver := setupUploadMutationResolver(tx)

		claims := auth.Claims{UserID: testUser.ID.String()}
		ctx = auth.ContextWithClaims(ctx, &claims)

		input := model.CreateImageUploadURLInput{
			FileName: "dish.png",
			FileType: "image/png",
			FileSize: 1024,
		}

		payload, err := mutResolver.CreateImageUploadURL(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, payload)
		require.NotEmpty(t, payload.UploadID)
		_, parseErr := uuid.Parse(payload.UploadID)
		require.NoError(t, parseErr, "UploadID should be a valid UUID")
		require.Contains(t, payload.UploadURL, "https://upload.example.com/")
		require.Contains(t, payload.FileURL, "https://cdn.example.com/")
		require.Contains(t, payload.UploadURL, payload.UploadID)
		require.Contains(t, payload.FileURL, payload.UploadID)
	})
}

func TestUploadResolver_CreateImageUploadURL_Unauthenticated(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutResolver := setupUploadMutationResolver(tx)

		input := model.CreateImageUploadURLInput{
			FileName: "dish.png",
			FileType: "image/png",
			FileSize: 1024,
		}

		payload, err := mutResolver.CreateImageUploadURL(context.Background(), input)

		require.Error(t, err)
		require.Nil(t, payload)

		gqlErr, ok := err.(*gqlerror.Error)
		require.True(t, ok, "Expected GraphQL error type")
		require.Equal(t, "user is not authenticated", gqlErr.Message)
		require.Equal(t, "UNAUTHENTICATED", gqlErr.Extensions["code"])
	})
}

func TestUploadResolver_CreateImageUploadURL_ServiceNotConfigured(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		r := &Resolver{UploadService: nil}
		mutResolver := &mutationResolver{r}

		claims := auth.Claims{UserID: testUser.ID.String()}
		ctx = auth.ContextWithClaims(ctx, &claims)

		input := model.CreateImageUploadURLInput{
			FileName: "dish.png",
			FileType: "image/png",
			FileSize: 1024,
		}

		payload, err := mutResolver.CreateImageUploadURL(ctx, input)

		require.Error(t, err)
		require.Nil(t, payload)

		gqlErr, ok := err.(*gqlerror.Error)
		require.True(t, ok, "Expected GraphQL error type")
		require.Equal(t, "INTERNAL_ERROR", gqlErr.Extensions["code"])
	})
}

func TestUploadResolver_CreateImageUploadURL_InvalidClaimsUserID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()

		mutResolver := setupUploadMutationResolver(tx)

		claims := auth.Claims{UserID: "not-a-uuid"}
		ctx = auth.ContextWithClaims(ctx, &claims)

		input := model.CreateImageUploadURLInput{
			FileName: "dish.png",
			FileType: "image/png",
			FileSize: 1024,
		}

		payload, err := mutResolver.CreateImageUploadURL(ctx, input)

		require.Error(t, err)
		require.Nil(t, payload)

		gqlErr, ok := err.(*gqlerror.Error)
		require.True(t, ok, "Expected GraphQL error type")
		require.Equal(t, "INTERNAL_ERROR", gqlErr.Extensions["code"])
	})
}

func TestUploadResolver_CreateImageUploadURL_UnsupportedFileType(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		mutResolver := setupUploadMutationResolver(tx)

		claims := auth.Claims{UserID: testUser.ID.String()}
		ctx = auth.ContextWithClaims(ctx, &claims)

		input := model.CreateImageUploadURLInput{
			FileName: "animation.gif",
			FileType: "image/gif",
			FileSize: 1024,
		}

		payload, err := mutResolver.CreateImageUploadURL(ctx, input)

		require.Error(t, err)
		require.Nil(t, payload)
	})
}
