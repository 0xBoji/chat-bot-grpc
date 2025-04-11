#!/bin/bash

# Build script for all services

set -e

# Generate protobuf files
echo "Generating protobuf files..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/auth/auth.proto proto/chat/chat.proto proto/room/room.proto

# Build auth service
echo "Building auth service..."
go build -o bin/auth-service cmd/auth-service/main.go

# Build chat service
echo "Building chat service..."
go build -o bin/chat-service cmd/chat-service/main.go

# Build room service
echo "Building room service..."
go build -o bin/room-service cmd/room-service/main.go

# Build gateway
echo "Building gateway..."
go build -o bin/gateway cmd/gateway/main.go

echo "Build completed successfully!"
