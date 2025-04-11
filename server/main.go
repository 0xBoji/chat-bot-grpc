package main

import (
	"log"
	"net"

	"grpc-messenger-core/db"
	pb "grpc-messenger-core/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

	// No Hello service anymore

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
