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

	// Create a debug wrapper to log requests
	debugHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the request
			logger.Printf("Request: %s %s", r.Method, r.URL.Path)

			// Log the authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				logger.Printf("Authorization header: %s", authHeader)
			} else {
				logger.Printf("No Authorization header found")
			}

			h.ServeHTTP(w, r)
		})
	}

	// Create a CORS wrapper
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the origin from the request
			origin := r.Header.Get("Origin")
			if origin == "" {
				// Default to localhost:3000 if no origin is provided
				origin = "http://localhost:3000"
			}

			// Allow localhost:3000 and localhost:3001
			if origin == "http://localhost:3000" || origin == "http://localhost:3001" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// For other origins, use a wildcard
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *gatewayPort),
		Handler: corsHandler(debugHandler(mux)),
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
