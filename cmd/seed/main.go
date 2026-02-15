package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"pawnshop/internal/config"
	"pawnshop/pkg/auth"
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

	pm := auth.NewPasswordManager()

	// Generate hash for admin123
	hash, err := pm.HashPassword("admin123")
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	fmt.Println("Generated hash:", hash)

	// Update admin user password
	result, err := db.Exec("UPDATE users SET password_hash = $1 WHERE email = $2", hash, "admin@pawnshop.com")
	if err != nil {
		log.Fatal("Failed to update password:", err)
	}

	rows, _ := result.RowsAffected()
	fmt.Printf("Updated %d user(s)\n", rows)

	// Verify the password works
	var storedHash string
	err = db.QueryRow("SELECT password_hash FROM users WHERE email = $1", "admin@pawnshop.com").Scan(&storedHash)
	if err != nil {
		log.Fatal("Failed to get hash:", err)
	}

	valid, err := pm.VerifyPassword("admin123", storedHash)
	if err != nil {
		log.Fatal("Failed to verify:", err)
	}

	if valid {
		fmt.Println("Password verification successful!")
	} else {
		fmt.Println("Password verification FAILED!")
	}
}
