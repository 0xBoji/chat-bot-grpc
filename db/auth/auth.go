package auth

import (
	"context"
	"database/sql"
)

// User represents a user in the database
type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

// Repository handles database operations for auth
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// UserExists checks if a user with the given username exists
func (r *Repository) UserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(ctx context.Context, username, passwordHash string) (int64, error) {
	var userID int64
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, username, passwordHash).Scan(&userID)
	return userID, err
}

// GetUserByUsername retrieves a user by username
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	return user, err
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password_hash FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.PasswordHash)
	return user, err
}
