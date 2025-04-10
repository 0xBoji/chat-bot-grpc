package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"strings"

	"first-grpc/db"
	pb "first-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
	db *sql.DB
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// Check for authentication token in metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata found in context")
	}

	// Extract token from metadata
	var userID int64
	var username string
	var authenticated bool

	auth := md.Get("authorization")
	if len(auth) > 0 && strings.HasPrefix(auth[0], "Bearer ") {
		token := strings.TrimPrefix(auth[0], "Bearer ")
		// Validate the token
		var err error
		userID, username, err = db.ValidateToken(token)
		if err == nil {
			authenticated = true
			log.Printf("Authenticated request from user: %s (ID: %d)", username, userID)
		} else {
			log.Printf("Invalid token: %v", err)
		}
	}

	// Use the database connection to log the request
	// This is just a simple example to demonstrate database usage
	if s.db != nil {
		// Create table if it doesn't exist
		_, createErr := s.db.ExecContext(ctx,
			"CREATE TABLE IF NOT EXISTS greetings (id SERIAL PRIMARY KEY, name TEXT, user_id BIGINT, timestamp TIMESTAMPTZ DEFAULT NOW())")
		if createErr != nil {
			log.Printf("Error creating table: %v", createErr)
		}

		// Insert the greeting request into the database with user ID if authenticated
		var insertErr error
		if authenticated {
			_, insertErr = s.db.ExecContext(ctx,
				"INSERT INTO greetings (name, user_id) VALUES ($1, $2)", in.GetName(), userID)
		} else {
			_, insertErr = s.db.ExecContext(ctx,
				"INSERT INTO greetings (name) VALUES ($1)", in.GetName())
		}

		if insertErr != nil {
			log.Printf("Error inserting greeting: %v", insertErr)
		} else {
			log.Printf("Saved greeting for: %s", in.GetName())
		}
	}

	// Customize the response based on authentication
	if authenticated {
		return &pb.HelloResponse{Message: "Hello " + in.GetName() + " (Authenticated as " + username + ")"}, nil
	} else {
		return &pb.HelloResponse{Message: "Hello " + in.GetName() + " (Unauthenticated)"}, nil
	}
}

func main() {
	// Connect to the database
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Create database tables if they don't exist
	if err := db.CreateUsersTable(dbConn); err != nil {
		log.Printf("Warning: failed to create users table: %v", err)
	}
	if err := db.CreateTablesIfNotExist(dbConn); err != nil {
		log.Printf("Warning: failed to create chat tables: %v", err)
	}

	// Create a gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the Hello service
	pb.RegisterHelloServiceServer(s, &helloServer{db: dbConn})

	// Register the Auth service
	auth := &authServer{db: dbConn}
	pb.RegisterAuthServiceServer(s, auth)

	// Register the Chat service
	chat := newChatServer(dbConn)
	pb.RegisterChatServiceServer(s, chat)

	// Register reflection service on gRPC server (useful for tools like grpcurl)
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
