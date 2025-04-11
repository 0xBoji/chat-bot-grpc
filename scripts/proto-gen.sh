#!/bin/bash

# Script to generate protobuf files

set -e

# Generate Go code
echo "Generating Go code from proto files..."
protoc -I. -I./third_party \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
  proto/auth/auth.proto proto/chat/chat.proto proto/room/room.proto

# Generate TypeScript code for Next.js frontend using protobuf-ts
echo "Generating TypeScript code for Next.js frontend..."
PROTO_OUT_DIR="../chatbox-next/src/proto-generated"

# Create the output directory if it doesn't exist
mkdir -p "$PROTO_OUT_DIR"

# Generate TypeScript code using protobuf-ts
protoc -I. -I./third_party \
  --plugin=protoc-gen-ts_proto=../chatbox-next/node_modules/.bin/protoc-gen-ts_proto \
  --ts_proto_out="$PROTO_OUT_DIR" \
  --ts_proto_opt=esModuleInterop=true \
  --ts_proto_opt=outputServices=grpc-web \
  --ts_proto_opt=env=browser \
  --ts_proto_opt=useOptionals=true \
  proto/auth/auth.proto proto/chat/chat.proto proto/room/room.proto

echo "Proto generation completed successfully!"
