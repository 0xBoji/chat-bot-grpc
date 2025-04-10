package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the database
type User struct {
	ID       int64
	Username string
	Email    string
	Password string
}

// JWT secret key - in production, this should be stored securely
const jwtSecret = "your-secret-key-change-this-in-production"

// CreateUsersTable creates the users table if it doesn't exist
func CreateUsersTable(db *sql.DB) error {
	if db == nil {
		return errors.New("database connection is nil")
	}
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	return err
}

// RegisterUser registers a new user
func RegisterUser(db *sql.DB, username, email, password string) (int64, error) {
	if db == nil {
		return 0, errors.New("database connection is nil")
	}
	// Check if username already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking username: %w", err)
	}
	if exists {
		return 0, errors.New("username already exists")
	}

	// Check if email already exists
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking email: %w", err)
	}
	if exists {
		return 0, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error hashing password: %w", err)
	}

	// Insert the new user
	var userID int64
	err = db.QueryRow(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		username, email, string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("error inserting user: %w", err)
	}

	return userID, nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func AuthenticateUser(db *sql.DB, username, password string) (string, int64, error) {
	var user User
	err := db.QueryRow(
		"SELECT id, username, password FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, errors.New("invalid username or password")
		}
		return "", 0, fmt.Errorf("error querying user: %w", err)
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", 0, errors.New("invalid username or password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, user.ID, nil
}

// ValidateToken validates a JWT token and returns the user ID and username
func ValidateToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["user_id"].(float64))
		username := claims["username"].(string)
		return userID, username, nil
	}

	return 0, "", errors.New("invalid token")
}
