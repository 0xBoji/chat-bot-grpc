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

	pb "grpc-messenger-core/proto"

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
	var currentRoomID int64

	for {
		fmt.Println("\n=== Chat Room Client ===")
		if token == "" {
			fmt.Println("1. Login")
			fmt.Println("2. Register")
			fmt.Println("3. Exit")
		} else {
			fmt.Printf("Logged in as: %s (ID: %d)\n", username, userID)
			if currentRoomID != 0 {
				fmt.Printf("Current room: %d\n", currentRoomID)
				fmt.Println("1. Send Message")
				fmt.Println("2. Get Room Messages")
				fmt.Println("3. Start Message Stream")
				fmt.Println("4. Leave Room")
				fmt.Println("5. List Rooms")
				fmt.Println("6. Create Room")
				fmt.Println("7. Join Room")
				fmt.Println("8. Logout")
				fmt.Println("9. Exit")
			} else {
				fmt.Println("1. List Rooms")
				fmt.Println("2. Create Room")
				fmt.Println("3. Join Room")
				fmt.Println("4. Logout")
				fmt.Println("5. Exit")
			}
		}
		fmt.Print("Choose an option: ")

		scanner.Scan()
		option := scanner.Text()

		if token == "" {
			// Not logged in
			switch option {
			case "1":
				token, userID, username = login(authClient, scanner)
			case "2":
				register(authClient, scanner)
			case "3":
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println("Invalid option")
			}
		} else if currentRoomID != 0 {
			// In a room
			switch option {
			case "1":
				sendMessage(chatClient, scanner, token, userID, currentRoomID)
			case "2":
				getRoomMessages(chatClient, scanner, token, userID, currentRoomID)
			case "3":
				streamRoomMessages(chatClient, token, userID, currentRoomID)
			case "4":
				leaveRoom(chatClient, token, userID, currentRoomID)
				currentRoomID = 0
			case "5":
				listRooms(chatClient, token, userID)
			case "6":
				newRoomID := createRoom(chatClient, scanner, token, userID)
				if newRoomID != 0 {
					currentRoomID = newRoomID
				}
			case "7":
				newRoomID := joinRoom(chatClient, scanner, token, userID)
				if newRoomID != 0 {
					currentRoomID = newRoomID
				}
			case "8":
				token = ""
				userID = 0
				username = ""
				currentRoomID = 0
			case "9":
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println("Invalid option")
			}
		} else {
			// Logged in but not in a room
			switch option {
			case "1":
				listRooms(chatClient, token, userID)
			case "2":
				newRoomID := createRoom(chatClient, scanner, token, userID)
				if newRoomID != 0 {
					currentRoomID = newRoomID
				}
			case "3":
				newRoomID := joinRoom(chatClient, scanner, token, userID)
				if newRoomID != 0 {
					currentRoomID = newRoomID
				}
			case "4":
				token = ""
				userID = 0
				username = ""
			case "5":
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println("Invalid option")
			}
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

func createRoom(client pb.ChatServiceClient, scanner *bufio.Scanner, token string, userID int64) int64 {
	fmt.Print("Enter room name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Enter room description: ")
	scanner.Scan()
	description := scanner.Text()

	fmt.Print("Is this a private room? (y/n): ")
	scanner.Scan()
	isPrivateStr := scanner.Text()
	isPrivate := strings.ToLower(isPrivateStr) == "y"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.CreateRoom(ctx, &pb.CreateRoomRequest{
		Name:        name,
		Description: description,
		CreatorId:   userID,
		IsPrivate:   isPrivate,
	})

	if err != nil {
		log.Printf("Failed to create room: %v", err)
		return 0
	}

	fmt.Printf("Room created successfully! Room ID: %d\n", resp.Id)
	return resp.Id
}

func listRooms(client pb.ChatServiceClient, token string, userID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetRooms(ctx, &pb.GetRoomsRequest{
		UserId:         userID,
		IncludePrivate: true,
		Limit:          50,
		Offset:         0,
	})

	if err != nil {
		log.Printf("Failed to list rooms: %v", err)
		return
	}

	if len(resp.Rooms) == 0 {
		fmt.Println("No rooms found")
		return
	}

	fmt.Println("\n=== Available Rooms ===")
	fmt.Printf("%-5s | %-20s | %-40s | %-10s\n", "ID", "Name", "Description", "Private")
	fmt.Println(strings.Repeat("-", 80))

	for _, room := range resp.Rooms {
		privateStr := "No"
		if room.IsPrivate {
			privateStr = "Yes"
		}
		fmt.Printf("%-5d | %-20s | %-40s | %-10s\n", room.Id, room.Name, room.Description, privateStr)
	}
}

func joinRoom(client pb.ChatServiceClient, scanner *bufio.Scanner, token string, userID int64) int64 {
	fmt.Print("Enter room ID to join: ")
	scanner.Scan()
	roomIDStr := scanner.Text()

	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid room ID")
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.JoinRoom(ctx, &pb.JoinRoomRequest{
		RoomId: roomID,
		UserId: userID,
	})

	if err != nil {
		log.Printf("Failed to join room: %v", err)
		return 0
	}

	if !resp.Success {
		fmt.Printf("Failed to join room: %s\n", resp.Message)
		return 0
	}

	fmt.Printf("Joined room successfully! Room: %s\n", resp.Room.Name)
	return roomID
}

func leaveRoom(client pb.ChatServiceClient, token string, userID int64, roomID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.LeaveRoom(ctx, &pb.LeaveRoomRequest{
		RoomId: roomID,
		UserId: userID,
	})

	if err != nil {
		log.Printf("Failed to leave room: %v", err)
		return
	}

	if !resp.Success {
		fmt.Printf("Failed to leave room: %s\n", resp.Message)
		return
	}

	fmt.Println("Left room successfully")
}

func sendMessage(client pb.ChatServiceClient, scanner *bufio.Scanner, token string, senderID int64, roomID int64) {
	fmt.Print("Enter message: ")
	scanner.Scan()
	content := scanner.Text()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.SendMessage(ctx, &pb.SendMessageRequest{
		Content:  content,
		SenderId: senderID,
		RoomId:   roomID,
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	if !resp.Success {
		fmt.Printf("Failed to send message: %s\n", resp.Message)
		return
	}

	fmt.Println("Message sent successfully")
}

func getRoomMessages(client pb.ChatServiceClient, scanner *bufio.Scanner, token string, userID int64, roomID int64) {
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

	resp, err := client.GetRoomMessages(ctx, &pb.GetRoomMessagesRequest{
		RoomId: roomID,
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
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.SenderName, msg.Content)
	}
}

func streamRoomMessages(client pb.ChatServiceClient, token string, userID int64, roomID int64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add token to context
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.StreamRoomMessages(ctx, &pb.StreamRoomMessagesRequest{
		RoomId: roomID,
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

			fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.SenderName, msg.Content)
		}
	}()

	// Wait for user to press Enter to stop streaming
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println("Stopping message stream...")
}
