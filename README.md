# gRPC Messenger Core

A modern chat application built with gRPC and Go, featuring authentication, room management, and real-time messaging. This project uses a microservice architecture with separate services for authentication, chat, and room management.

## Features

- **User Authentication**:
  - Register and login users
  - JWT-based token validation
  - Secure password hashing

- **Room Management**:
  - Create, join, and leave rooms
  - Public and private rooms
  - Room membership management

- **Real-time Messaging**:
  - Send and receive messages in real-time
  - Message history retrieval
  - Streaming messages using gRPC streams

- **API Gateway**:
  - RESTful API endpoints via gRPC-Gateway
  - Frontend integration with gRPC-Web
  - Unified API access point

## Project Structure

```bash
.
├── Makefile                # Build and run commands
├── cmd                     # Entry points for services
│   ├── auth-service        # Authentication service
│   │   └── main.go
│   ├── chat-service        # Chat service
│   │   └── main.go
│   ├── room-service        # Room service
│   │   └── main.go
│   ├── gateway             # API Gateway
│   │   └── main.go
│   └── clients             # Test clients
│       ├── auth_client
│       ├── chat_client
│       └── room_chat_client
├── config                  # Configuration files
│   ├── auth-service
│   ├── chat-service
│   ├── room-service
│   ├── gateway
│   └── config.toml
├── db                      # Database layer
│   ├── postgres            # PostgreSQL connection
│   ├── auth                # Authentication database operations
│   ├── chat                # Chat database operations
│   └── room                # Room database operations
├── internal                # Internal packages
│   ├── auth                # Authentication service implementation
│   ├── chat                # Chat service implementation
│   ├── room                # Room service implementation
│   └── middleware          # Shared middleware
├── proto                   # Protocol buffer definitions
│   ├── auth                # Authentication service proto
│   ├── chat                # Chat service proto
│   ├── room                # Room service proto
│   └── google              # Google API proto imports
├── scripts                 # Utility scripts
│   ├── build.sh            # Build script
│   ├── deploy.sh           # Deployment script
│   ├── init-db.sh          # Database initialization
│   ├── init-db.sql         # SQL schema
│   ├── proto-gen.sh        # Proto generation script
│   └── run-tests.sh        # Test runner
├── docs                    # Documentation
│   ├── README.md
│   ├── frontend_integration.md
│   └── grpc_web_setup.md
└── chatbox-next            # Next.js frontend
    └── src                 # Frontend source code
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL database
- Protocol Buffers compiler (protoc)
- Node.js 18 or later (for frontend)
- Docker and Docker Compose (optional, for containerized deployment)

## Getting Started

1. Clone the repository
2. Configure your database in `config/config.toml`
3. Initialize the database:

   ```bash
   make init-db
   ```

4. Generate protocol buffer files:

   ```bash
   make proto
   ```

5. Start all services:

   ```bash
   make server
   ```

6. In a separate terminal, start the frontend:

   ```bash
   make frontend
   ```

7. Or start everything at once:

   ```bash
   make start-all
   ```

8. For development and testing, you can run individual services:

   ```bash
   make auth-service
   make chat-service
   make room-service
   make gateway
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
  - Room membership management

- **Messaging**:
  - Send messages to rooms
  - Retrieve message history for a room
  - Stream real-time messages in a room using gRPC streaming
  - Optimistic UI updates for better user experience

## Frontend Integration

The project includes a Next.js frontend in the `chatbox-next` directory. The frontend uses:

- Next.js 14 with App Router
- TypeScript for type safety
- gRPC-Web for communication with the backend
- Tailwind CSS for styling

For more details, see the documentation in the `docs` directory:

- [Frontend Integration Guide](docs/frontend_integration.md)
- [gRPC-Web Setup Guide](docs/grpc_web_setup.md)

## Testing

Run the tests with:

```bash
make test
```

Run tests with coverage report:

```bash
make test-coverage
```

## License

MIT
