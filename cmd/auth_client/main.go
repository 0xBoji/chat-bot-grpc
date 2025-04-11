package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "grpc-messenger-core/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewAuthServiceClient(conn)

	// We'll create a new context for each request

	// Interactive menu
	for {
		fmt.Println("\n=== Auth Client ===")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Validate Token")
		fmt.Println("4. Exit")
		fmt.Print("Choose an option: ")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			register(context.Background(), client)
		case 2:
			login(context.Background(), client)
		case 3:
			validateToken(context.Background(), client)
		case 4:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func register(_ context.Context, client pb.AuthServiceClient) {
	// Create a context with timeout for this request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var username, email, password string

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter email: ")
	fmt.Scanln(&email)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	resp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}

	if resp.Success {
		fmt.Printf("Registration successful! User ID: %d\n", resp.UserId)
	} else {
		fmt.Printf("Registration failed: %s\n", resp.Message)
	}
}

func login(_ context.Context, client pb.AuthServiceClient) {
	// Create a context with timeout for this request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var username, password string

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	resp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	if resp.Success {
		fmt.Printf("Login successful! User ID: %d\n", resp.UserId)
		fmt.Printf("Token: %s\n", resp.Token)
		// Save token to a file for later use
		fmt.Println("Token saved for validation")
	} else {
		fmt.Printf("Login failed: %s\n", resp.Message)
	}
}

func validateToken(_ context.Context, client pb.AuthServiceClient) {
	// Create a context with timeout for this request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var token string

	fmt.Print("Enter token: ")
	fmt.Scanln(&token)

	resp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})

	if err != nil {
		log.Fatalf("Token validation failed: %v", err)
	}

	if resp.Valid {
		fmt.Printf("Token is valid! User ID: %d, Username: %s\n", resp.UserId, resp.Username)
	} else {
		fmt.Println("Token is invalid")
	}
}
