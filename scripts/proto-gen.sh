#!/bin/bash

# Script to generate protobuf files

set -e

# Generate Go code
echo "Generating Go code from proto files..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/auth/auth.proto proto/chat/chat.proto proto/room/room.proto

# Generate gRPC-Web code for frontend
echo "Generating gRPC-Web code for frontend..."
protoc --js_out=import_style=commonjs:frontend/src/generated \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/src/generated \
  proto/auth/auth.proto proto/chat/chat.proto proto/room/room.proto

echo "Proto generation completed successfully!"
