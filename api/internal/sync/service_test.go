package sync

import (
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/reference"
	"foodplanner/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSyncIngredientData(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := t.Context()
		filePath := "../../reference/ingredients_test.yaml"
		loader := reference.NewLoader(filePath)
		ingredientService := ingredient.NewIngredientService(tx, ingredient.NewIngredientRepo())
		service := NewSyncService(ingredientService, loader)

		// Act
		err := service.SyncIngredientData(ctx)

		// Assert
		require.NoError(t, err)

		// Fetch ingredient via ingredient service to ensure it's populate
	})
}
