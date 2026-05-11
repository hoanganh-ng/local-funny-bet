package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: migrate <up|down>\n")
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintf(os.Stderr, "DATABASE_URL not set\n")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "pinging database: %v\n", err)
		os.Exit(1)
	}

	if err := ensureMigrationsTable(db); err != nil {
		fmt.Fprintf(os.Stderr, "creating schema_migrations table: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			fmt.Fprintf(os.Stderr, "running migrations: %v\n", err)
			os.Exit(1)
		}
	case "down":
		if err := migrateDown(db); err != nil {
			fmt.Fprintf(os.Stderr, "rolling back migration: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		os.Exit(1)
	}
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func migrateUp(db *sql.DB) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("getting applied migrations: %w", err)
	}

	files, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("reading migration files: %w", err)
	}

	var count int
	for _, file := range files {
		version := strings.TrimSuffix(file, ".sql")
		if applied[version] {
			continue
		}

		content, err := os.ReadFile(filepath.Join("backend/migrations", file))
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("starting transaction for %s: %w", file, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("executing migration %s: %w", file, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("recording migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing migration %s: %w", file, err)
		}

		fmt.Printf("Applied: %s\n", file)
		count++
	}

	if count == 0 {
		fmt.Println("No pending migrations")
	} else {
		fmt.Printf("Applied %d migration(s)\n", count)
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("getting applied migrations: %w", err)
	}

	files, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("reading migration files: %w", err)
	}

	var lastApplied string
	for i := len(files) - 1; i >= 0; i-- {
		version := strings.TrimSuffix(files[i], ".sql")
		if applied[version] {
			lastApplied = version
			break
		}
	}

	if lastApplied == "" {
		fmt.Println("No migrations to roll back")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if _, err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", getTableNameFromMigration(lastApplied))); err != nil {
		tx.Rollback()
		return fmt.Errorf("dropping table: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", lastApplied); err != nil {
		tx.Rollback()
		return fmt.Errorf("removing migration record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing rollback: %w", err)
	}

	fmt.Printf("Rolled back: %s.sql\n", lastApplied)
	return nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func getMigrationFiles() ([]string, error) {
	entries, err := os.ReadDir("backend/migrations")
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") && !strings.Contains(entry.Name(), "seed") {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}

func getTableNameFromMigration(version string) string {
	parts := strings.Split(version, "_")
	if len(parts) < 3 {
		return ""
	}
	return strings.Join(parts[2:], "_")
}
