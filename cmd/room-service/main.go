package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"grpc-messenger-core/db/postgres"
	"grpc-messenger-core/internal/room"
	pb "grpc-messenger-core/proto/room"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50053, "The server port")
)

func main() {
	flag.Parse()

	// Initialize logger
	logger := log.New(os.Stdout, "[ROOM-SERVICE] ", log.LstdFlags)
	logger.Println("Starting Room Service...")

	// Connect to database
	db, err := postgres.NewPostgresDB()
	if err != nil {
		logger.Printf("Warning: Failed to connect to database: %v", err)
		logger.Println("Continuing without database connection for testing purposes...")
	} else {
		defer db.Close()
	}

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()

	// Create room service
	roomService := room.NewRoomService(db, logger)

	// Register service
	pb.RegisterRoomServiceServer(s, roomService)

	// Register reflection service for development tools
	reflection.Register(s)

	// Start server in a goroutine
	go func() {
		logger.Printf("Room service listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down Room Service...")
	s.GracefulStop()
	logger.Println("Room Service stopped")
}
