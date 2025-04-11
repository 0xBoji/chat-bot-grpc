package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Database connection parameters
const (
	defaultHost     = "localhost"
	defaultPort     = 5432
	defaultUser     = "postgres"
	defaultPassword = ""
	defaultDBName   = "postgres"
	defaultSSLMode  = "disable"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*sql.DB, error) {
	// Get connection parameters from environment variables or use defaults
	host := getEnv("DB_HOST", defaultHost)
	port := getEnv("DB_PORT", fmt.Sprintf("%d", defaultPort))
	user := getEnv("DB_USER", defaultUser)
	password := getEnv("DB_PASSWORD", defaultPassword)
	dbname := getEnv("DB_NAME", defaultDBName)
	sslmode := getEnv("DB_SSLMODE", defaultSSLMode)

	// Create connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
