package recipe

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewIngredientUsage_InvalidQuantity(t *testing.T) {
	// Arrange
	request := CreateIngredientUsageRequest{
		IngredientID: uuid.New().String(),
		Quantity:     -1,
		Unit:         1,
	}

	// Act
	_, err := NewIngredientUsage(request)

	// Assert
	require.ErrorIs(t, err, ErrInvalidQuantity)
}

func TestNewIngredientUsage_InvalidUnit(t *testing.T) {
	// Arrange
	request := CreateIngredientUsageRequest{
		IngredientID: uuid.New().String(),
		Quantity:     1.5,
		Unit:         -1,
	}

	// Act
	_, err := NewIngredientUsage(request)

	// Assert
	require.ErrorIs(t, err, ErrInvalidUnit)
}
