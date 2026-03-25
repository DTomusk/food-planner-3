package recipe

import (
	"encoding/base64"
	"encoding/json"
	"math"
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

func TestRecipeCursor_RoundTrip_DefaultModeNewest(t *testing.T) {
	expectedCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
	expectedID := uuid.MustParse("ca51c083-e08e-40fa-8e0d-496590aead83")

	cursor := (&RecipeCursor{
		CreatedAt: expectedCreatedAt,
		ID:        expectedID,
	}).String()

	parsed, err := ParseRecipeCursor(&cursor)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, RecipeCursorModeNewest, parsed.Mode)
	require.True(t, expectedCreatedAt.Equal(parsed.CreatedAt))
	require.Equal(t, expectedID, parsed.ID)
	require.Nil(t, parsed.RelevanceScore)
}

func TestParseRecipeCursor_Version1Compatibility(t *testing.T) {
	expectedCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
	expectedID := uuid.MustParse("ca51c083-e08e-40fa-8e0d-496590aead86")

	payload := recipeCursorPayload{
		Version:   1,
		CreatedAt: expectedCreatedAt.Format(time.RFC3339Nano),
		ID:        expectedID.String(),
	}
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	cursor := base64.RawURLEncoding.EncodeToString(encoded)

	parsed, err := ParseRecipeCursor(&cursor)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, RecipeCursorModeNewest, parsed.Mode)
	require.True(t, expectedCreatedAt.Equal(parsed.CreatedAt))
	require.Equal(t, expectedID, parsed.ID)
	require.Nil(t, parsed.RelevanceScore)
}

func TestRecipeCursor_RoundTrip_NewestWithFilterHash(t *testing.T) {
	expectedCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
	expectedID := uuid.MustParse("ca51c083-e08e-40fa-8e0d-496590aead84")

	cursor := (&RecipeCursor{
		Mode:       RecipeCursorModeNewest,
		FilterHash: "newest|none",
		CreatedAt:  expectedCreatedAt,
		ID:         expectedID,
	}).String()

	parsed, err := ParseRecipeCursor(&cursor)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, RecipeCursorModeNewest, parsed.Mode)
	require.Equal(t, "newest|none", parsed.FilterHash)
	require.True(t, expectedCreatedAt.Equal(parsed.CreatedAt))
	require.Equal(t, expectedID, parsed.ID)
	require.Nil(t, parsed.RelevanceScore)
}

func TestRecipeCursor_RoundTrip_RelevanceWithScore(t *testing.T) {
	expectedCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
	expectedID := uuid.MustParse("ca51c083-e08e-40fa-8e0d-496590aead85")
	score := 0.9234

	cursor := (&RecipeCursor{
		Mode:           RecipeCursorModeRelevance,
		FilterHash:     "q=chicken|sort=relevance",
		CreatedAt:      expectedCreatedAt,
		ID:             expectedID,
		RelevanceScore: &score,
	}).String()

	parsed, err := ParseRecipeCursor(&cursor)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, RecipeCursorModeRelevance, parsed.Mode)
	require.Equal(t, "q=chicken|sort=relevance", parsed.FilterHash)
	require.True(t, expectedCreatedAt.Equal(parsed.CreatedAt))
	require.Equal(t, expectedID, parsed.ID)
	require.NotNil(t, parsed.RelevanceScore)
	require.InDelta(t, score, *parsed.RelevanceScore, 0.0000001)
}

func TestParseRecipeCursor_InvalidCursor(t *testing.T) {
	invalid := "not-a-valid-cursor"

	parsed, err := ParseRecipeCursor(&invalid)
	require.ErrorIs(t, err, ErrInvalidCursor)
	require.Nil(t, parsed)
}

func TestParseRecipeCursor_InvalidMode(t *testing.T) {
	payload := recipeCursorPayload{
		Version:   recipeCursorVersion,
		Mode:      "oldest",
		CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		ID:        uuid.New().String(),
	}
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	cursor := base64.RawURLEncoding.EncodeToString(encoded)

	parsed, err := ParseRecipeCursor(&cursor)
	require.ErrorIs(t, err, ErrInvalidCursor)
	require.Nil(t, parsed)
}

func TestParseRecipeCursor_RelevanceMissingScore(t *testing.T) {
	payload := recipeCursorPayload{
		Version:   recipeCursorVersion,
		Mode:      string(RecipeCursorModeRelevance),
		CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		ID:        uuid.New().String(),
	}
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	cursor := base64.RawURLEncoding.EncodeToString(encoded)

	parsed, err := ParseRecipeCursor(&cursor)
	require.ErrorIs(t, err, ErrInvalidCursor)
	require.Nil(t, parsed)
}

func TestRecipeCursor_String_InvalidValues(t *testing.T) {
	zeroTimeCursor := (&RecipeCursor{ID: uuid.New()}).String()
	require.Equal(t, "", zeroTimeCursor)

	nilIDCursor := (&RecipeCursor{CreatedAt: time.Now().UTC()}).String()
	require.Equal(t, "", nilIDCursor)

	invalidModeCursor := (&RecipeCursor{
		Mode:      RecipeCursorMode("oldest"),
		CreatedAt: time.Now().UTC(),
		ID:        uuid.New(),
	}).String()
	require.Equal(t, "", invalidModeCursor)

	relevanceNoScoreCursor := (&RecipeCursor{
		Mode:      RecipeCursorModeRelevance,
		CreatedAt: time.Now().UTC(),
		ID:        uuid.New(),
	}).String()
	require.Equal(t, "", relevanceNoScoreCursor)

	nan := math.NaN()
	relevanceNaNCursor := (&RecipeCursor{
		Mode:           RecipeCursorModeRelevance,
		CreatedAt:      time.Now().UTC(),
		ID:             uuid.New(),
		RelevanceScore: &nan,
	}).String()
	require.Equal(t, "", relevanceNaNCursor)
}

func TestRecipeCursor_Matches(t *testing.T) {
	cursor := &RecipeCursor{Mode: RecipeCursorModeRelevance, FilterHash: "q=tomato"}
	require.True(t, cursor.Matches(RecipeCursorModeRelevance, "q=tomato"))
	require.False(t, cursor.Matches(RecipeCursorModeNewest, "q=tomato"))
	require.False(t, cursor.Matches(RecipeCursorModeRelevance, "q=pepper"))
	require.True(t, ((*RecipeCursor)(nil)).Matches(RecipeCursorModeNewest, ""))
}
