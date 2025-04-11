FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /bin/auth-service ./cmd/auth-service
RUN go build -o /bin/chat-service ./cmd/chat-service
RUN go build -o /bin/room-service ./cmd/room-service
RUN go build -o /bin/gateway ./cmd/gateway

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /bin/auth-service /app/auth-service
COPY --from=builder /bin/chat-service /app/chat-service
COPY --from=builder /bin/room-service /app/room-service
COPY --from=builder /bin/gateway /app/gateway

# Expose ports
EXPOSE 50051 50052 50053 8082

# Set environment variables
ENV AUTH_SERVICE_PORT=50051
ENV CHAT_SERVICE_PORT=50052
ENV ROOM_SERVICE_PORT=50053
ENV GATEWAY_PORT=8082

# Command to run
CMD ["sh", "-c", "echo 'Starting services...' && ./auth-service & ./chat-service & ./room-service & ./gateway --gateway-port=8082 && wait"]
