package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"first-grpc/db"
	pb "first-grpc/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
	db *sql.DB
	// For managing active streaming connections
	mu      sync.Mutex
	streams map[int64][]pb.ChatService_StreamRoomMessagesServer
}

// newChatServer creates a new chat server instance
func newChatServer(db *sql.DB) *chatServer {
	return &chatServer{
		db:      db,
		streams: make(map[int64][]pb.ChatService_StreamRoomMessagesServer),
	}
}

// CreateRoom handles creating a new chat room
func (s *chatServer) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.RoomResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the creator ID matches the authenticated user
	if userID != req.CreatorId {
		return nil, status.Errorf(codes.PermissionDenied, "creator ID does not match authenticated user")
	}

	// Create the room
	room, err := db.CreateRoom(ctx, s.db, req.Name, req.Description, req.CreatorId, req.IsPrivate)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create room: %v", err)
	}

	// Convert to proto response
	return &pb.RoomResponse{
		Id:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		CreatorId:   room.CreatorID,
		IsPrivate:   room.IsPrivate,
		CreatedAt:   room.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetRooms handles retrieving rooms that a user can access
func (s *chatServer) GetRooms(ctx context.Context, req *pb.GetRoomsRequest) (*pb.GetRoomsResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Set default values for limit and offset if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Get rooms from the database
	rooms, err := db.GetRooms(ctx, s.db, req.UserId, req.IncludePrivate, limit, offset)
	if err != nil {
		log.Printf("Error getting rooms: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve rooms")
	}

	// Convert to proto response
	var pbRooms []*pb.RoomResponse
	for _, room := range rooms {
		pbRooms = append(pbRooms, &pb.RoomResponse{
			Id:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			CreatorId:   room.CreatorID,
			IsPrivate:   room.IsPrivate,
			CreatedAt:   room.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetRoomsResponse{
		Rooms: pbRooms,
	}, nil
}

// JoinRoom handles a user joining a room
func (s *chatServer) JoinRoom(ctx context.Context, req *pb.JoinRoomRequest) (*pb.JoinRoomResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Join the room
	_, err = db.JoinRoom(ctx, s.db, req.RoomId, req.UserId)
	if err != nil {
		log.Printf("Error joining room: %v", err)
		return &pb.JoinRoomResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Get the room details
	room, err := db.GetRoomByID(ctx, s.db, req.RoomId)
	if err != nil {
		log.Printf("Error getting room: %v", err)
		return &pb.JoinRoomResponse{
			Success: true,
			Message: "Joined room successfully, but failed to retrieve room details",
		}, nil
	}

	return &pb.JoinRoomResponse{
		Success: true,
		Message: "Joined room successfully",
		Room: &pb.RoomResponse{
			Id:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			CreatorId:   room.CreatorID,
			IsPrivate:   room.IsPrivate,
			CreatedAt:   room.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// LeaveRoom handles a user leaving a room
func (s *chatServer) LeaveRoom(ctx context.Context, req *pb.LeaveRoomRequest) (*pb.LeaveRoomResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Leave the room
	err = db.LeaveRoom(ctx, s.db, req.RoomId, req.UserId)
	if err != nil {
		log.Printf("Error leaving room: %v", err)
		return &pb.LeaveRoomResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LeaveRoomResponse{
		Success: true,
		Message: "Left room successfully",
	}, nil
}

// SendMessage handles sending a new chat message
func (s *chatServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the sender ID matches the authenticated user
	if userID != req.SenderId {
		return nil, status.Errorf(codes.PermissionDenied, "sender ID does not match authenticated user")
	}

	// Save the message to the database
	messageID, err := db.SaveMessage(ctx, s.db, req.Content, req.SenderId, req.RoomId)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		return &pb.SendMessageResponse{
			Success: false,
			Message: "Failed to save message: " + err.Error(),
		}, nil
	}

	// Get sender's username for the notification
	senderName, err := db.GetUserNameByID(ctx, s.db, req.SenderId)
	if err != nil {
		log.Printf("Error getting sender name: %v", err)
		senderName = "Unknown User"
	}

	// Create message response for streaming
	timestamp := time.Now().Format(time.RFC3339)
	msgResponse := &pb.MessageResponse{
		Id:         messageID,
		Content:    req.Content,
		SenderId:   req.SenderId,
		RoomId:     req.RoomId,
		SenderName: senderName,
		Timestamp:  timestamp,
	}

	// Broadcast the message to all users in the room
	s.broadcastMessage(msgResponse, req.RoomId)

	return &pb.SendMessageResponse{
		Success:   true,
		Message:   "Message sent successfully",
		MessageId: messageID,
	}, nil
}

// GetRoomMessages handles retrieving messages for a room
func (s *chatServer) GetRoomMessages(ctx context.Context, req *pb.GetRoomMessagesRequest) (*pb.GetRoomMessagesResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Set default values for limit and offset if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Get messages from the database
	messages, err := db.GetRoomMessages(ctx, s.db, req.RoomId, req.UserId, limit, offset)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve messages")
	}

	// Convert database messages to protobuf messages
	var pbMessages []*pb.MessageResponse
	for _, msg := range messages {
		// Get sender's username
		senderName, err := db.GetUserNameByID(ctx, s.db, msg.SenderID)
		if err != nil {
			log.Printf("Error getting sender name: %v", err)
			senderName = "Unknown User"
		}

		pbMessages = append(pbMessages, &pb.MessageResponse{
			Id:         msg.ID,
			Content:    msg.Content,
			SenderId:   msg.SenderID,
			RoomId:     msg.RoomID,
			SenderName: senderName,
			Timestamp:  msg.Timestamp.Format(time.RFC3339),
		})
	}

	return &pb.GetRoomMessagesResponse{
		Messages: pbMessages,
	}, nil
}

// StreamRoomMessages establishes a streaming connection for real-time messages in a room
func (s *chatServer) StreamRoomMessages(req *pb.StreamRoomMessagesRequest, stream pb.ChatService_StreamRoomMessagesServer) error {
	// Get context from the stream
	ctx := stream.Context()

	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Check if the user is a member of the room
	isMember, err := db.IsRoomMember(ctx, s.db, req.RoomId, req.UserId)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check room membership: %v", err)
	}

	if !isMember {
		return status.Errorf(codes.PermissionDenied, "user is not a member of this room")
	}

	// Register this stream for the room
	s.mu.Lock()
	if _, ok := s.streams[req.RoomId]; !ok {
		s.streams[req.RoomId] = make([]pb.ChatService_StreamRoomMessagesServer, 0)
	}
	s.streams[req.RoomId] = append(s.streams[req.RoomId], stream)
	s.mu.Unlock()

	// Keep the stream open until the client disconnects or context is canceled
	<-ctx.Done()

	// Remove the stream when the client disconnects
	s.mu.Lock()
	defer s.mu.Unlock()

	if streams, ok := s.streams[req.RoomId]; ok {
		for i, str := range streams {
			if str == stream {
				// Remove this stream from the slice
				s.streams[req.RoomId] = append(streams[:i], streams[i+1:]...)
				break
			}
		}

		// If no more streams for this room, remove the room entry
		if len(s.streams[req.RoomId]) == 0 {
			delete(s.streams, req.RoomId)
		}
	}

	return nil
}

// broadcastMessage sends a message to all connected clients in a room
func (s *chatServer) broadcastMessage(msg *pb.MessageResponse, roomID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Send to all clients in the room
	if streams, ok := s.streams[roomID]; ok {
		for _, stream := range streams {
			if err := stream.Send(msg); err != nil {
				log.Printf("Error sending message to stream in room %d: %v", roomID, err)
			}
		}
	}
}

// authenticateRequest extracts and validates the JWT token from the request metadata
func authenticateRequest(ctx context.Context) (int64, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return 0, "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// The auth token should be in the format "Bearer <token>"
	authHeader := values[0]
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return 0, "", status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}

	token := authHeader[7:]
	userID, username, err := db.ValidateToken(token)
	if err != nil {
		return 0, "", status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return userID, username, nil
}
