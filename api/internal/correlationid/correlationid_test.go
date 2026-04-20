package correlationid

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWithContextAndFromContext_RoundTripsCorrelationID(t *testing.T) {
	expectedID := uuid.New()
	ctx := WithContext(context.Background(), expectedID)

	actualID := FromContext(ctx)

	require.Equal(t, expectedID, actualID)
}

func TestFromContext_GeneratesCorrelationIDWhenMissing(t *testing.T) {
	actualID := FromContext(context.Background())

	require.NotEqual(t, uuid.Nil, actualID)
}
