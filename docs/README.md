# gRPC Authentication Service Documentation

Welcome to the documentation for our gRPC authentication service. This documentation is intended for frontend developers who need to integrate with our backend services.

## Table of Contents

1. [Frontend Integration Guide](frontend_integration.md)
   - Overview of available services
   - Authentication flow
   - Making authenticated requests
   - Error handling
   - Security considerations
   - Example integration with React

2. [gRPC-Web Setup Guide](grpc_web_setup.md)
   - Installation steps for different platforms
   - Setting up your frontend project
   - Generating JavaScript client code
   - Setting up Envoy proxy
   - Using the generated client code
   - Troubleshooting

## Getting Started

If you're new to this project, we recommend starting with the [Frontend Integration Guide](frontend_integration.md) to understand the available services and authentication flow.

Once you're ready to set up your frontend project, follow the [gRPC-Web Setup Guide](grpc_web_setup.md) for detailed instructions.

## Available Services

Our backend provides the following gRPC services:

1. **AuthService**: Handles user registration, login, and token validation
2. **HelloService**: A simple greeting service that supports authenticated requests

## Authentication Flow

Our authentication system uses JWT tokens. The typical flow is:

1. Register a new user account
2. Login to receive a JWT token
3. Include the token in the metadata for authenticated requests
4. Validate the token when needed

## Need Help?

If you have any questions or encounter issues while integrating with our services, please contact the backend team for assistance.
