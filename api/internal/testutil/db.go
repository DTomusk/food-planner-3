package testutil

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetTestDB() *sql.DB {
	once.Do(func() {
		db_url := os.Getenv("TEST_DB_URL")
		if db_url == "" {
			log.Fatal("TEST_DB_URL environment variable is not set")
		}

		var err error
		db, err = sql.Open("postgres", db_url)
		if err != nil {
			log.Fatalf("Failed to connect to test database: %v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping test database: %v", err)
		}
	})
	return db
}
