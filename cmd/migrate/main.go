package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"pawnshop/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := sql.Open("pgx", cfg.Database.URL())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Parse command
	flag.Parse()
	command := flag.Arg(0)

	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatal("Migration failed:", err)
		}
		fmt.Println("Migrations completed successfully")
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Println("Rollback completed successfully")
	case "status":
		if err := showStatus(db); err != nil {
			log.Fatal("Failed to show status:", err)
		}
	default:
		fmt.Println("Usage: migrate [up|down|status]")
		os.Exit(1)
	}
}

func migrateUp(db *sql.DB) error {
	// Create migrations table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		applied[version] = true
	}

	// Get migration files
	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}
	sort.Strings(files)

	// Apply pending migrations
	for _, file := range files {
		version := extractVersion(file)
		if applied[version] {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		fmt.Printf("Applying %s...\n", version)

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply %s: %w", file, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		fmt.Printf("Applied %s\n", version)
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	// Get last applied migration
	var version string
	err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No migrations to rollback")
			return nil
		}
		return err
	}

	// Find corresponding down file
	downFile := fmt.Sprintf("migrations/%s.down.sql", version)
	content, err := os.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", downFile, err)
	}

	fmt.Printf("Rolling back %s...\n", version)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(content)); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback %s: %w", version, err)
	}

	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("Rolled back %s\n", version)
	return nil
}

func showStatus(db *sql.DB) error {
	rows, err := db.Query("SELECT version, applied_at FROM schema_migrations ORDER BY version")
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Applied migrations:")
	for rows.Next() {
		var version string
		var appliedAt string
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return err
		}
		fmt.Printf("  %s (applied at %s)\n", version, appliedAt)
	}

	return nil
}

func extractVersion(filename string) string {
	base := filepath.Base(filename)
	// Remove .up.sql or .down.sql
	version := strings.TrimSuffix(base, ".up.sql")
	version = strings.TrimSuffix(version, ".down.sql")
	return version
}
