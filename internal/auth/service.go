package auth

import (
	"context"
	"database/sql"
	"log"

	"grpc-messenger-core/db/auth"
	"grpc-messenger-core/internal/middleware"
	pb "grpc-messenger-core/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthService implements the AuthService gRPC service
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	db     *sql.DB
	logger *log.Logger
	repo   *auth.Repository
}

// NewAuthService creates a new auth service
func NewAuthService(db *sql.DB, logger *log.Logger) *AuthService {
	return &AuthService{
		db:     db,
		logger: logger,
		repo:   auth.NewRepository(db),
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.logger.Printf("Register request for user: %s", req.Username)

	// Validate request
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	// Check if user already exists
	exists, err := s.repo.UserExists(ctx, req.Username)
	if err != nil {
		s.logger.Printf("Error checking if user exists: %v", err)
		return nil, status.Error(codes.Internal, "failed to check if user exists")
	}
	if exists {
		return &pb.RegisterResponse{
			Success: false,
			Message: "username already exists",
		}, nil
	}

	// Hash password
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		s.logger.Printf("Error hashing password: %v", err)
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Create user
	userID, err := s.repo.CreateUser(ctx, req.Username, hashedPassword)
	if err != nil {
		s.logger.Printf("Error creating user: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "user registered successfully",
		UserId:  userID,
	}, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Printf("Login request for user: %s", req.Username)

	// Validate request
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	// MOCK AUTHENTICATION FOR TESTING
	// In a real implementation, we would check the database
	s.logger.Printf("Using mock authentication for testing")

	// Generate a mock user ID based on the username
	userID := int64(len(req.Username))

	// Generate token
	token, err := middleware.GenerateToken(userID, req.Username)
	if err != nil {
		s.logger.Printf("Error generating token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		Success:  true,
		Message:  "login successful (mock)",
		Token:    token,
		UserId:   userID,
		Username: req.Username,
	}, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	s.logger.Println("ValidateToken request")

	// Validate request
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Validate token
	claims, err := middleware.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:   false,
			Message: "invalid token",
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:     true,
		Message:   "token is valid",
		UserId:    claims.UserID,
		Username:  claims.Username,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}
