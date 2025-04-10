package main

import (
	"context"
	"log"
	"time"

	pb "first-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create auth client
	authClient := pb.NewAuthServiceClient(conn)
	
	// Create hello client
	helloClient := pb.NewHelloServiceClient(conn)

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Login to get a token
	log.Println("Logging in...")
	loginResp, err := authClient.Login(ctx, &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	if !loginResp.Success {
		log.Fatalf("Login failed: %s", loginResp.Message)
	}

	log.Printf("Login successful! User ID: %d", loginResp.UserId)
	log.Printf("Token: %s", loginResp.Token)

	// Create a new context with the token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + loginResp.Token,
	})
	authCtx := metadata.NewOutgoingContext(ctx, md)

	// Call the Hello service with the authenticated context
	log.Println("Calling SayHello with authentication...")
	helloResp, err := helloClient.SayHello(authCtx, &pb.HelloRequest{
		Name: "Authenticated User",
	})

	if err != nil {
		log.Fatalf("SayHello failed: %v", err)
	}

	log.Printf("Greeting: %s", helloResp.GetMessage())
}
