package chat

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"grpc-messenger-core/db/chat"
	"grpc-messenger-core/internal/middleware"
	pb "grpc-messenger-core/proto/chat"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ChatService implements the ChatService gRPC service
type ChatService struct {
	pb.UnimplementedChatServiceServer
	db     *sql.DB
	logger *log.Logger
	repo   *chat.Repository

	// For testing purposes
	mockMessagesMutex  sync.Mutex
	mockMessages       map[int64][]*pb.MessageResponse                     // roomID -> messages
	activeStreams      map[int64][]pb.ChatService_StreamRoomMessagesServer // roomID -> streams
	activeStreamsMutex sync.Mutex
}

// NewChatService creates a new chat service
func NewChatService(db *sql.DB, logger *log.Logger) *ChatService {
	// Set the global logger
	sharedLogger = logger

	return &ChatService{
		db:            db,
		logger:        logger,
		repo:          chat.NewRepository(db),
		mockMessages:  make(map[int64][]*pb.MessageResponse),
		activeStreams: make(map[int64][]pb.ChatService_StreamRoomMessagesServer),
	}
}

// SendMessage sends a message to a room
func (s *ChatService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.SenderId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Validate request
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message content cannot be empty")
	}

	// For testing purposes, if db is nil, return success and simulate a message
	if s.db == nil {
		s.logger.Println("Database connection is nil, returning mock send message response")

		// Create a mock message ID
		mockMessageID := time.Now().Unix()

		// Get username from context if available
		_, username, _ := authenticateRequest(ctx)
		if username == "" {
			username = "User"
		}

		// Simulate notifying subscribers by broadcasting to all active streams
		go func() {
			// Create a mock message
			mockMessage := chat.Message{
				ID:         mockMessageID,
				Content:    req.Content,
				SenderID:   req.SenderId,
				RoomID:     req.RoomId,
				SenderName: username,
				Timestamp:  time.Now(),
			}

			// Notify subscribers if repository exists
			if s.repo != nil {
				s.repo.NotifyRoomSubscribers(req.RoomId, mockMessage)
			} else {
				// For testing, directly broadcast to all active streams
				s.logger.Printf("Broadcasting mock message: %s", req.Content)

				// Create a mock message response immediately
				mockResponse := &pb.MessageResponse{
					Id:         mockMessageID,
					Content:    req.Content,
					SenderId:   req.SenderId,
					RoomId:     req.RoomId,
					SenderName: username,
					Timestamp:  time.Now().Format(time.RFC3339),
				}

				s.logger.Printf("Created mock message: %+v", mockResponse)

				// Store the message in memory and broadcast to all clients
				s.storeMockMessageAndBroadcast(mockResponse)

				// Also directly add to the global store to ensure it's stored
				AddMessage(mockResponse)

				// Also log the current state of stored messages
				s.mockMessagesMutex.Lock()
				s.logger.Printf("Current stored messages for room %d: %d messages", req.RoomId, len(s.mockMessages[req.RoomId]))
				s.mockMessagesMutex.Unlock()
			}
		}()

		return &pb.SendMessageResponse{
			Success:   true,
			Message:   "message sent successfully",
			MessageId: mockMessageID,
		}, nil
	}

	// Check if the user is a member of the room
	isMember, err := s.repo.IsRoomMember(ctx, req.RoomId, req.SenderId)
	if err != nil {
		s.logger.Printf("Error checking room membership: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check room membership")
	}
	if !isMember {
		return &pb.SendMessageResponse{
			Success: false,
			Message: "user is not a member of the room",
		}, nil
	}

	// Save message to database
	messageID, err := s.repo.SaveMessage(ctx, req.Content, req.SenderId, req.RoomId)
	if err != nil {
		s.logger.Printf("Error saving message: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to save message")
	}

	return &pb.SendMessageResponse{
		Success:   true,
		Message:   "message sent successfully",
		MessageId: messageID,
	}, nil
}

// GetRoomMessages retrieves messages from a room
func (s *ChatService) GetRoomMessages(ctx context.Context, req *pb.GetRoomMessagesRequest) (*pb.GetRoomMessagesResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// For testing purposes, if db is nil, return mock messages
	if s.db == nil {
		s.logger.Println("Database connection is nil, returning mock messages")

		// Create a welcome message
		welcomeMessage := &pb.MessageResponse{
			Id:         1,
			Content:    "Welcome to the chat!",
			SenderId:   0,
			RoomId:     req.RoomId,
			SenderName: "System",
			Timestamp:  time.Now().Format(time.RFC3339),
		}

		// Get messages from the global shared store
		globalMessages := GetMessages(req.RoomId)
		s.logger.Printf("Found %d messages in global store for room %d", len(globalMessages), req.RoomId)

		// Start with the welcome message
		mockMessages := []*pb.MessageResponse{welcomeMessage}

		// Add messages from the global store
		if len(globalMessages) > 0 {
			mockMessages = append(mockMessages, globalMessages...)
		}

		return &pb.GetRoomMessagesResponse{
			Messages: mockMessages,
		}, nil
	}

	// Check if the user is a member of the room
	isMember, err := s.repo.IsRoomMember(ctx, req.RoomId, req.UserId)
	if err != nil {
		s.logger.Printf("Error checking room membership: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check room membership")
	}
	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of the room")
	}

	// Get messages from database
	messages, err := s.repo.GetRoomMessages(ctx, req.RoomId, req.Limit, req.Offset)
	if err != nil {
		s.logger.Printf("Error getting messages: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get messages")
	}

	// Convert to protobuf messages
	pbMessages := make([]*pb.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		pbMessages = append(pbMessages, &pb.MessageResponse{
			Id:         msg.ID,
			Content:    msg.Content,
			SenderId:   msg.SenderID,
			RoomId:     msg.RoomID,
			SenderName: msg.SenderName,
			Timestamp:  msg.Timestamp.Format(time.RFC3339),
		})
	}

	return &pb.GetRoomMessagesResponse{
		Messages: pbMessages,
	}, nil
}

// StreamRoomMessages establishes a streaming connection for real-time messages in a room
func (s *ChatService) StreamRoomMessages(req *pb.StreamRoomMessagesRequest, stream pb.ChatService_StreamRoomMessagesServer) error {
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

	// For testing purposes, if db is nil, just continue and set up mock streaming
	if s.db == nil {
		s.logger.Println("Database connection is nil, continuing with streaming")

		// Register this stream for the room
		s.activeStreamsMutex.Lock()
		s.activeStreams[req.RoomId] = append(s.activeStreams[req.RoomId], stream)
		streamIndex := len(s.activeStreams[req.RoomId]) - 1
		s.logger.Printf("Registered stream %d for room %d. Total streams for this room: %d", streamIndex, req.RoomId, len(s.activeStreams[req.RoomId]))
		s.activeStreamsMutex.Unlock()

		// Deregister on exit
		defer func() {
			s.activeStreamsMutex.Lock()
			s.logger.Printf("Deregistering stream %d for room %d", streamIndex, req.RoomId)
			if streamIndex < len(s.activeStreams[req.RoomId]) {
				// Remove this stream from the active streams
				s.activeStreams[req.RoomId] = append(
					s.activeStreams[req.RoomId][:streamIndex],
					s.activeStreams[req.RoomId][streamIndex+1:]...,
				)
				s.logger.Printf("Stream %d removed. Remaining streams for room %d: %d", streamIndex, req.RoomId, len(s.activeStreams[req.RoomId]))
			} else {
				s.logger.Printf("Stream %d already removed or index out of bounds. Current streams for room %d: %d", streamIndex, req.RoomId, len(s.activeStreams[req.RoomId]))
			}
			s.activeStreamsMutex.Unlock()
		}()

		// Send any existing mock messages
		s.mockMessagesMutex.Lock()
		for _, msg := range s.mockMessages[req.RoomId] {
			err := stream.Send(msg)
			if err != nil {
				s.logger.Printf("Error sending message to client: %v", err)
				s.mockMessagesMutex.Unlock()
				return status.Errorf(codes.Internal, "failed to send message to client")
			}
		}
		s.mockMessagesMutex.Unlock()

		// Wait for context to be done
		<-ctx.Done()
		return nil
	} else {
		// Check if the user is a member of the room
		isMember, err := s.repo.IsRoomMember(ctx, req.RoomId, req.UserId)
		if err != nil {
			s.logger.Printf("Error checking room membership: %v", err)
			return status.Errorf(codes.Internal, "failed to check room membership")
		}
		if !isMember {
			return status.Errorf(codes.PermissionDenied, "user is not a member of the room")
		}

		// Create a message channel
		messageChan := make(chan chat.Message)

		// Subscribe to room messages
		s.repo.SubscribeToRoom(req.RoomId, messageChan)
		defer s.repo.UnsubscribeFromRoom(req.RoomId, messageChan)

		// Stream messages to client
		for {
			select {
			case msg := <-messageChan:
				// Send message to client
				err := stream.Send(&pb.MessageResponse{
					Id:         msg.ID,
					Content:    msg.Content,
					SenderId:   msg.SenderID,
					RoomId:     msg.RoomID,
					SenderName: msg.SenderName,
					Timestamp:  msg.Timestamp.Format(time.RFC3339),
				})
				if err != nil {
					s.logger.Printf("Error sending message to client: %v", err)
					return status.Errorf(codes.Internal, "failed to send message to client")
				}
			case <-ctx.Done():
				// Client disconnected
				return nil
			}
		}
	}
}

// storeMockMessageAndBroadcast stores a mock message and broadcasts it to all active streams
func (s *ChatService) storeMockMessageAndBroadcast(message *pb.MessageResponse) {
	// Store the message in both the local and global stores
	s.mockMessagesMutex.Lock()
	// Initialize the slice if it doesn't exist
	if s.mockMessages[message.RoomId] == nil {
		s.mockMessages[message.RoomId] = make([]*pb.MessageResponse, 0)
	}
	s.mockMessages[message.RoomId] = append(s.mockMessages[message.RoomId], message)
	s.logger.Printf("Stored message in local store for room %d. Total messages: %d", message.RoomId, len(s.mockMessages[message.RoomId]))
	s.mockMessagesMutex.Unlock()

	// Also store in the global shared store
	// Create a copy of the message to store in the global store
	globalMessage := &pb.MessageResponse{
		Id:         message.Id,
		Content:    message.Content,
		SenderId:   message.SenderId,
		RoomId:     message.RoomId,
		SenderName: message.SenderName,
		Timestamp:  message.Timestamp,
	}

	AddMessage(globalMessage)
	s.logger.Printf("Stored message in global store for room %d with ID %d", message.RoomId, globalMessage.Id)

	// Broadcast to all active streams for this room
	s.activeStreamsMutex.Lock()
	streams := s.activeStreams[message.RoomId]
	s.logger.Printf("Broadcasting message to %d active streams for room %d", len(streams), message.RoomId)

	for i, stream := range streams {
		// Use non-blocking send to avoid deadlocks
		go func(st pb.ChatService_StreamRoomMessagesServer, idx int) {
			s.logger.Printf("Sending message to stream %d for room %d: %+v", idx, message.RoomId, message)
			err := st.Send(message)
			if err != nil {
				s.logger.Printf("Error sending message to stream %d: %v", idx, err)
			} else {
				s.logger.Printf("Successfully sent message to stream %d for room %d", idx, message.RoomId)
			}
		}(stream, i)
	}
	s.activeStreamsMutex.Unlock()
}

// Helper function to authenticate a request
func authenticateRequest(ctx context.Context) (int64, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return 0, "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// Extract token from "Bearer <token>"
	token := authHeader[0]
	if len(token) <= 7 || token[:7] != "Bearer " {
		return 0, "", status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}
	token = token[7:]

	// Validate token
	claims, err := middleware.ValidateToken(token)
	if err != nil {
		return 0, "", status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return claims.UserID, claims.Username, nil
}
