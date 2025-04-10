package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "first-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create clients
	authClient := pb.NewAuthServiceClient(conn)
	chatClient := pb.NewChatServiceClient(conn)

	// Interactive menu
	scanner := bufio.NewScanner(os.Stdin)
	var token string
	var userID int64
	var username string

	for {
		fmt.Println("\n=== Chat Bot Client ===")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Send Message")
		fmt.Println("4. Get Messages")
		fmt.Println("5. Start Message Stream")
		fmt.Println("6. Exit")
		fmt.Print("Choose an option: ")

		scanner.Scan()
		option := scanner.Text()

		switch option {
		case "1":
			token, userID, username = login(authClient, scanner)
		case "2":
			register(authClient, scanner)
		case "3":
			if token == "" {
				fmt.Println("You must login first")
				continue
			}
			sendMessage(chatClient, scanner, token, userID)
		case "4":
			if token == "" {
				fmt.Println("You must login first")
				continue
			}
			getMessages(chatClient, scanner, token, userID)
		case "5":
			if token == "" {
				fmt.Println("You must login first")
				continue
			}
			streamMessages(chatClient, token, userID)
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func login(client pb.AuthServiceClient, scanner *bufio.Scanner) (string, int64, string) {
	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Printf("Login failed: %v", err)
		return "", 0, ""
	}

	if !resp.Success {
		fmt.Printf("Login failed: %s\n", resp.Message)
		return "", 0, ""
	}

	fmt.Printf("Login successful! User ID: %d\n", resp.UserId)
	return resp.Token, resp.UserId, username
}

func register(client pb.AuthServiceClient, scanner *bufio.Scanner) {
	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter email: ")
	scanner.Scan()
	email := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Printf("Registration failed: %v", err)
		return
	}

	if resp.Success {
		fmt.Printf("Registration successful! User ID: %d\n", resp.UserId)
	} else {
		fmt.Printf("Registration failed: %s\n", resp.Message)
	}
}

func sendMessage(client pb.ChatServiceClient, scanner *bufio.Scanner, token string, senderID int64) {
	fmt.Print("Enter message content: ")
	scanner.Scan()
	content := scanner.Text()

	fmt.Print("Enter receiver ID (leave empty for broadcast): ")
	scanner.Scan()
	receiverIDStr := scanner.Text()

	var receiverID int64
	if receiverIDStr != "" {
		var err error
		receiverID, err = strconv.ParseInt(receiverIDStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid receiver ID, using broadcast")
			receiverID = 0
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.SendMessage(ctx, &pb.SendMessageRequest{
		Content:    content,
		SenderId:   senderID,
		ReceiverId: receiverID,
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	if resp.Success {
		fmt.Printf("Message sent successfully! Message ID: %d\n", resp.MessageId)
	} else {
		fmt.Printf("Failed to send message: %s\n", resp.Message)
	}
}

func getMessages(client pb.ChatServiceClient, scanner *bufio.Scanner, token string, userID int64) {
	fmt.Print("Enter limit (default 10): ")
	scanner.Scan()
	limitStr := scanner.Text()

	fmt.Print("Enter offset (default 0): ")
	scanner.Scan()
	offsetStr := scanner.Text()

	var limit, offset int64 = 10, 0
	if limitStr != "" {
		var err error
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit <= 0 {
			limit = 10
		}
	}
	if offsetStr != "" {
		var err error
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || offset < 0 {
			offset = 0
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetMessages(ctx, &pb.GetMessagesRequest{
		UserId: userID,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		log.Printf("Failed to get messages: %v", err)
		return
	}

	if len(resp.Messages) == 0 {
		fmt.Println("No messages found")
		return
	}

	fmt.Println("\n=== Messages ===")
	for _, msg := range resp.Messages {
		receiverStr := "Broadcast"
		if msg.ReceiverId != 0 {
			receiverStr = fmt.Sprintf("User %d", msg.ReceiverId)
		}
		fmt.Printf("[%s] %s -> %s: %s\n", msg.Timestamp, msg.SenderName, receiverStr, msg.Content)
	}
}

func streamMessages(client pb.ChatServiceClient, token string, userID int64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.StreamMessages(ctx, &pb.StreamMessagesRequest{
		UserId: userID,
	})
	if err != nil {
		log.Printf("Failed to start message stream: %v", err)
		return
	}

	fmt.Println("\n=== Live Message Stream ===")
	fmt.Println("(Press Enter to stop streaming)")

	// Start a goroutine to receive messages
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Stream closed by server")
				return
			}
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}

			receiverStr := "Broadcast"
			if msg.ReceiverId != 0 {
				receiverStr = fmt.Sprintf("User %d", msg.ReceiverId)
			}
			fmt.Printf("[%s] %s -> %s: %s\n", msg.Timestamp, msg.SenderName, receiverStr, msg.Content)
		}
	}()

	// Wait for user to press Enter to stop streaming
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println("Stopping message stream...")
}
