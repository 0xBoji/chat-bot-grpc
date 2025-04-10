package main

import (
	"first-grpc/db"
	"fmt"
	"log"
)

func main() {
	// Connect to the database
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Query the table schema
	rows, err := conn.Query(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'greetings'
	`)
	if err != nil {
		log.Fatalf("Failed to query schema: %v", err)
	}
	defer rows.Close()

	fmt.Println("Greetings table schema:")
	fmt.Println("------------------------")
	fmt.Printf("%-20s | %-20s\n", "Column Name", "Data Type")
	fmt.Println("------------------------")

	// Iterate through the results
	for rows.Next() {
		var columnName, dataType string
		
		if err := rows.Scan(&columnName, &dataType); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		
		fmt.Printf("%-20s | %-20s\n", columnName, dataType)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
}
