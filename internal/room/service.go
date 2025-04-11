package room

import (
	"context"
	"database/sql"
	"log"

	"grpc-messenger-core/db/room"
	"grpc-messenger-core/internal/middleware"
	pb "grpc-messenger-core/proto/room"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RoomService implements the RoomService gRPC service
type RoomService struct {
	pb.UnimplementedRoomServiceServer
	db     *sql.DB
	logger *log.Logger
	repo   *room.Repository
}

// NewRoomService creates a new room service
func NewRoomService(db *sql.DB, logger *log.Logger) *RoomService {
	return &RoomService{
		db:     db,
		logger: logger,
		repo:   room.NewRepository(db),
	}
}

// CreateRoom creates a new chat room
func (s *RoomService) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.RoomResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.CreatorId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Validate request
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "room name cannot be empty")
	}

	// Create room in database
	roomID, err := s.repo.CreateRoom(ctx, req.Name, req.Description, req.CreatorId)
	if err != nil {
		s.logger.Printf("Error creating room: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create room")
	}

	// Add creator as a member
	err = s.repo.AddRoomMember(ctx, roomID, req.CreatorId)
	if err != nil {
		s.logger.Printf("Error adding creator as member: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to add creator as member")
	}

	return &pb.RoomResponse{
		Id:          roomID,
		Name:        req.Name,
		Description: req.Description,
		CreatorId:   req.CreatorId,
	}, nil
}

// GetRooms retrieves all rooms the user is a member of
func (s *RoomService) GetRooms(ctx context.Context, req *pb.GetRoomsRequest) (*pb.GetRoomsResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Get rooms from database
	rooms, err := s.repo.GetUserRooms(ctx, req.UserId)
	if err != nil {
		s.logger.Printf("Error getting rooms: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get rooms")
	}

	// Convert to protobuf rooms
	pbRooms := make([]*pb.RoomResponse, 0, len(rooms))
	for _, r := range rooms {
		pbRooms = append(pbRooms, &pb.RoomResponse{
			Id:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			CreatorId:   r.CreatorID,
		})
	}

	return &pb.GetRoomsResponse{
		Rooms: pbRooms,
	}, nil
}

// JoinRoom adds a user to a room
func (s *RoomService) JoinRoom(ctx context.Context, req *pb.JoinRoomRequest) (*pb.JoinRoomResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Check if room exists
	roomExists, err := s.repo.RoomExists(ctx, req.RoomId)
	if err != nil {
		s.logger.Printf("Error checking if room exists: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check if room exists")
	}
	if !roomExists {
		return &pb.JoinRoomResponse{
			Success: false,
			Message: "room does not exist",
		}, nil
	}

	// Check if user is already a member
	isMember, err := s.repo.IsRoomMember(ctx, req.RoomId, req.UserId)
	if err != nil {
		s.logger.Printf("Error checking room membership: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check room membership")
	}
	if isMember {
		return &pb.JoinRoomResponse{
			Success: false,
			Message: "user is already a member of the room",
		}, nil
	}

	// Add user to room
	err = s.repo.AddRoomMember(ctx, req.RoomId, req.UserId)
	if err != nil {
		s.logger.Printf("Error adding user to room: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to add user to room")
	}

	return &pb.JoinRoomResponse{
		Success: true,
		Message: "user joined room successfully",
	}, nil
}

// LeaveRoom removes a user from a room
func (s *RoomService) LeaveRoom(ctx context.Context, req *pb.LeaveRoomRequest) (*pb.LeaveRoomResponse, error) {
	// Authenticate the user
	userID, _, err := authenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the user ID matches the authenticated user
	if userID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "user ID does not match authenticated user")
	}

	// Check if user is a member
	isMember, err := s.repo.IsRoomMember(ctx, req.RoomId, req.UserId)
	if err != nil {
		s.logger.Printf("Error checking room membership: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to check room membership")
	}
	if !isMember {
		return &pb.LeaveRoomResponse{
			Success: false,
			Message: "user is not a member of the room",
		}, nil
	}

	// Remove user from room
	err = s.repo.RemoveRoomMember(ctx, req.RoomId, req.UserId)
	if err != nil {
		s.logger.Printf("Error removing user from room: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to remove user from room")
	}

	return &pb.LeaveRoomResponse{
		Success: true,
		Message: "user left room successfully",
	}, nil
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
