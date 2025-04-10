# gRPC Chat Room Service

A chat room service built with gRPC and Go, featuring authentication, room management, and real-time messaging.

## Features

- User authentication (register, login, token validation)
- Room management (create, join, leave rooms)
- Public and private rooms
- Real-time chat messaging within rooms
- Message history retrieval
- Secure JWT-based authentication

## Project Structure

```bash
.
├── Makefile                # Build and run commands
├── client                  # Simple client for testing
│   └── main.go
├── cmd                     # Command-line tools
│   ├── auth_client         # Authentication client
│   │   └── main.go
│   ├── chat_client         # Chat client
│   │   └── main.go
│   └── db_test             # Database test utility
│       └── main.go
├── config                  # Configuration files
│   └── config.toml
├── db                      # Database layer
│   ├── auth.go             # Authentication database operations
│   ├── chat.go             # Chat database operations
│   └── postgres.go         # PostgreSQL connection
├── docs                    # Documentation
│   ├── README.md
│   ├── frontend_integration.md
│   └── grpc_web_setup.md
├── proto                   # Protocol buffer definitions
│   ├── auth.proto          # Authentication service
│   ├── chat.proto          # Chat service
│   └── hello.proto         # Hello service (example)
└── server                  # Server implementation
    ├── auth_server.go      # Authentication service implementation
    ├── chat_server.go      # Chat service implementation
    └── main.go             # Main server
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL database
- Protocol Buffers compiler (protoc)

## Getting Started

1. Clone the repository
2. Configure your database in `config/config.toml`
3. Generate protocol buffer files:

   ```bash
   make proto
   ```

4. Start the server:

   ```bash
   make server
   ```

5. Run the room chat client:

   ```bash
   make room-chat-client
   ```

## Authentication

The service uses JWT tokens for authentication. To use authenticated endpoints:

1. Register a user account
2. Login to get a JWT token
3. Include the token in the request metadata

## Chat Features

- **Room Management**:
  - Create public or private rooms
  - Join existing rooms
  - Leave rooms
  - List available rooms

- **Messaging**:
  - Send messages to rooms
  - Retrieve message history for a room
  - Stream real-time messages in a room

## Frontend Integration

For frontend developers, see the documentation in the `docs` directory:

- [Frontend Integration Guide](docs/frontend_integration.md)
- [gRPC-Web Setup Guide](docs/grpc_web_setup.md)

## License

MIT
