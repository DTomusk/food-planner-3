package recipe

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseRecipeCursor_NilCursor(t *testing.T) {
	parsed, err := ParseRecipeCursor(nil)
	require.NoError(t, err)
	require.Nil(t, parsed)
}

func TestParseRecipeCursor_EmptyCursor(t *testing.T) {
	empty := ""

	parsed, err := ParseRecipeCursor(&empty)
	require.NoError(t, err)
	require.Nil(t, parsed)
}

func TestRecipeCursor_RoundTrip(t *testing.T) {
	expectedCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
	expectedID := uuid.MustParse("ca51c083-e08e-40fa-8e0d-496590aead83")

	cursor := (&RecipeCursor{
		CreatedAt: expectedCreatedAt,
		ID:        expectedID,
	}).String()

	parsed, err := ParseRecipeCursor(&cursor)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.True(t, expectedCreatedAt.Equal(parsed.CreatedAt))
	require.Equal(t, expectedID, parsed.ID)
}

func TestParseRecipeCursor_InvalidCursor(t *testing.T) {
	invalid := "not-a-valid-cursor"

	parsed, err := ParseRecipeCursor(&invalid)
	require.ErrorIs(t, err, ErrInvalidCursor)
	require.Nil(t, parsed)
}

func TestRecipeCursor_String_InvalidValues(t *testing.T) {
	zeroTimeCursor := (&RecipeCursor{ID: uuid.New()}).String()
	require.Equal(t, "", zeroTimeCursor)

	nilIDCursor := (&RecipeCursor{CreatedAt: time.Now().UTC()}).String()
	require.Equal(t, "", nilIDCursor)
}
