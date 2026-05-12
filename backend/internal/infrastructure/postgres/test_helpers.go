package postgres

import (
	"database/sql"
	"os"
	"testing"

	"wc2026/migrations"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgresql:postgresql@localhost:5432/wc2026?sslmode=disable"
		// t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	db, err := NewDB(dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Drop all tables for clean slate
	_, err = db.Exec(`
		DROP SCHEMA IF EXISTS public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		t.Fatalf("failed to reset schema: %v", err)
	}

	// Run all migrations
	if err := Up(db, migrations.FS); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(`
			DROP SCHEMA IF EXISTS public CASCADE;
			CREATE SCHEMA public;
		`)
		if err != nil {
			t.Logf("failed to cleanup schema: %v", err)
		}
		db.Close()
	})

	return db
}
