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

	// Login
	log.Println("Attempting to login...")
	resp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	if resp.Success {
		log.Printf("Login successful! User ID: %d", resp.UserId)
		log.Printf("Token: %s", resp.Token)

		// Now validate the token
		log.Println("Validating token...")
		validateResp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{
			Token: resp.Token,
		})

		if err != nil {
			log.Fatalf("Token validation failed: %v", err)
		}

		if validateResp.Valid {
			log.Printf("Token is valid! User ID: %d, Username: %s", validateResp.UserId, validateResp.Username)
		} else {
			log.Println("Token is invalid")
		}
	} else {
		log.Printf("Login failed: %s", resp.Message)
	}
}
