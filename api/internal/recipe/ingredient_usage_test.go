package recipe

import (
	"foodplanner/internal/unit"
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

func TestNewIngredientUsage_Success(t *testing.T) {
	// Arrange
	request := CreateIngredientUsageRequest{
		IngredientID: uuid.New().String(),
		Quantity:     2,
		Unit:         int(unit.Quantum),
	}

	// Act
	usage, err := NewIngredientUsage(request)

	// Assert
	require.NoError(t, err)
	require.Equal(t, uuid.MustParse(request.IngredientID), usage.IngredientID)
	require.Equal(t, request.Quantity, usage.Quantity)
	require.Equal(t, unit.Quantum, usage.Unit)
}
