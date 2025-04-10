package main

import (
	"context"
	"log"
	"time"

	pb "first-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewAuthServiceClient(conn)

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Register a user
	log.Println("Attempting to register a user...")
	resp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}

	if resp.Success {
		log.Printf("Registration successful! User ID: %d", resp.UserId)
	} else {
		log.Printf("Registration failed: %s", resp.Message)
	}
}
