package main

import (
	"context"
	"database/sql"
	"log"

	"first-grpc/db"
	pb "first-grpc/proto"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

// Register handles user registration
func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Received register request for user: %s", req.Username)
	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Username, email, and password are required",
		}, nil
	}

	// Create users table if it doesn't exist
	log.Println("Creating users table if it doesn't exist")
	if err := db.CreateUsersTable(s.db); err != nil {
		log.Printf("Error creating users table: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

	// Register the user
	log.Printf("Registering user: %s", req.Username)
	userID, err := db.RegisterUser(s.db, req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  userID,
	}, nil
}

// Login handles user login
func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return &pb.LoginResponse{
			Success: false,
			Message: "Username and password are required",
		}, nil
	}

	// Authenticate the user
	token, userID, err := db.AuthenticateUser(s.db, req.Username, req.Password)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		UserId:  userID,
	}, nil
}

// ValidateToken validates a JWT token
func (s *authServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// Validate input
	if req.Token == "" {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Validate the token
	userID, username, err := db.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:    true,
		UserId:   userID,
		Username: username,
	}, nil
}
