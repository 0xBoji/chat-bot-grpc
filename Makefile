.PHONY: proto build run clean test auth-service chat-service room-service gateway clients

# Generate protobuf files
proto:
	./scripts/proto-gen.sh

# Build all services
build: proto
	./scripts/build.sh

# Run services
auth-service:
	go run cmd/auth-service/main.go

chat-service:
	go run cmd/chat-service/main.go

room-service:
	go run cmd/room-service/main.go

gateway:
	go run cmd/gateway/main.go

# Run clients
auth-client:
	go run cmd/clients/auth_client/main.go

chat-client:
	go run cmd/clients/chat_client/main.go

room-chat-client:
	go run cmd/clients/room_chat_client/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Deploy services
deploy: build
	./scripts/deploy.sh