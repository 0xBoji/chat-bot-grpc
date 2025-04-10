package main

import (
	"database/sql"
	"first-grpc/db"
	"fmt"
	"log"
	"time"
)

func main() {
	// Connect to the database
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Query the greetings table
	rows, err := conn.Query("SELECT id, name, user_id, timestamp FROM greetings ORDER BY timestamp DESC")
	if err != nil {
		log.Fatalf("Failed to query greetings: %v", err)
	}
	defer rows.Close()

	fmt.Println("Greetings stored in the database:")
	fmt.Println("----------------------------------")
	fmt.Printf("%-5s | %-20s | %-10s | %-30s\n", "ID", "Name", "User ID", "Timestamp")
	fmt.Println("----------------------------------")

	// Iterate through the results
	for rows.Next() {
		var id int
		var name string
		var userID sql.NullInt64
		var timestamp time.Time

		if err := rows.Scan(&id, &name, &userID, &timestamp); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		userIDStr := "-"
		if userID.Valid {
			userIDStr = fmt.Sprintf("%d", userID.Int64)
		}

		fmt.Printf("%-5d | %-20s | %-10s | %s\n", id, name, userIDStr, timestamp.Format(time.RFC3339))
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
}
