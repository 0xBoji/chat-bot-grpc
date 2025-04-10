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

	// If we get here, the connection was successful
	log.Println("Database connection test successful!")

	// Test a simple query
	var version string
	err = conn.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	log.Printf("PostgreSQL version: %s", version)
}
