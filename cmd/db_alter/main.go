package main

import (
	"first-grpc/db"
	"log"
)

func main() {
	// Connect to the database
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Alter the table to add the user_id column
	_, err = conn.Exec(`
		ALTER TABLE greetings 
		ADD COLUMN IF NOT EXISTS user_id BIGINT
	`)
	if err != nil {
		log.Fatalf("Failed to alter table: %v", err)
	}

	log.Println("Successfully altered greetings table to add user_id column")
}
