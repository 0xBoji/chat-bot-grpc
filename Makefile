.PHONY: proto build run clean test test-coverage auth-service chat-service room-service gateway clients server start-all frontend init-db

# Generate protobuf files
proto:
	./scripts/proto-gen.sh

# Build all services
build: proto
	./scripts/build.sh

# Run all servers (in separate terminals)
server:
	@echo "Starting all services..."
	@echo "Starting Auth Service on port 50051"
	go run cmd/auth-service/main.go &
	@echo "Starting Chat Service on port 50052"
	go run cmd/chat-service/main.go &
	@echo "Starting Room Service on port 50053"
	go run cmd/room-service/main.go &
	@echo "Starting API Gateway on port 8082"
	go run cmd/gateway/main.go --gateway-port=8082 &
	@echo "All services are running. Press Ctrl+C to stop all services."
	@wait

# Run Next.js frontend
frontend:
	cd ../chatbox-next && npm run dev

# Run all services and frontend
start-all:
	@echo "Starting all services and frontend..."
	@echo "Starting backend services..."
	make server &
	@echo "Starting frontend..."
	make frontend &
	@echo "All services and frontend are running. Press Ctrl+C to stop all services."
	@wait

# Run individual services
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
	./scripts/run-tests.sh

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

# Initialize database
init-db:
	./scripts/init-db.sh

# Deploy services
deploy: build
	./scripts/deploy.sh