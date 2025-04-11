package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "grpc-messenger-core/proto/auth"
	chatpb "grpc-messenger-core/proto/chat"
	roompb "grpc-messenger-core/proto/room"
)

var (
	// gRPC service addresses
	authServiceAddr = flag.String("auth-service", "localhost:50051", "Auth service address")
	chatServiceAddr = flag.String("chat-service", "localhost:50052", "Chat service address")
	roomServiceAddr = flag.String("room-service", "localhost:50053", "Room service address")
	
	// Gateway port
	gatewayPort = flag.Int("gateway-port", 8080, "Gateway server port")
)

func main() {
	flag.Parse()

	// Initialize logger
	logger := log.New(os.Stdout, "[GATEWAY] ", log.LstdFlags)
	logger.Println("Starting API Gateway...")

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// Register handlers
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register Auth service
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, *authServiceAddr, opts)
	if err != nil {
		logger.Fatalf("Failed to register auth service handler: %v", err)
	}

	// Register Chat service
	err = chatpb.RegisterChatServiceHandlerFromEndpoint(ctx, mux, *chatServiceAddr, opts)
	if err != nil {
		logger.Fatalf("Failed to register chat service handler: %v", err)
	}

	// Register Room service
	err = roompb.RegisterRoomServiceHandlerFromEndpoint(ctx, mux, *roomServiceAddr, opts)
	if err != nil {
		logger.Fatalf("Failed to register room service handler: %v", err)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *gatewayPort),
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("API Gateway listening at %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down API Gateway...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}
	logger.Println("API Gateway stopped")
}
