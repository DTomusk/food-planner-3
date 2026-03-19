package refreshtokens

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	// Arrange
	secret := "my-secret-key"
	hasher := NewHasher(secret)
	token := "sample-refresh-token"

	// Act
	hashedToken := hasher.hash(token)

	// Assert
	require.NotEqual(t, token, hashedToken)
}
