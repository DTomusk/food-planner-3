package recipe

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("postgres", "postgres://user:password@db_test:5432/food_planner_test?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := testDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	code := m.Run() // run tests
	os.Exit(code)
}

func withTx(t *testing.T, fn func(tx *sql.Tx)) {
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatalf("Failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	fn(tx)
}

func TestCreateAndGetRecipe(t *testing.T) {
	r := NewRepo()

	withTx(t, func(tx *sql.Tx) {
		id := uuid.New()
		err := r.CreateRecipe(Recipe{
			ID:   id,
			Name: "Chocolate Cake",
		}, context.Background(), tx)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}

		got, err := r.GetRecipeByID(id.String(), context.Background(), tx)
		if err != nil {
			t.Fatalf("Failed to get recipe: %v", err)
		}

		if got.Name != "Chocolate Cake" {
			t.Errorf("Expected name %q, got %q", "Chocolate Cake", got.Name)
		}
	})
}
