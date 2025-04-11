# Frontend Integration Guide for React with Vite

This document provides guidelines for frontend developers to integrate with our gRPC chat room service using React with Vite.

## Overview

Our backend provides the following gRPC services:

1. **AuthService**: Handles user registration, login, and token validation
2. **ChatService**: Manages chat rooms and messaging functionality

## Prerequisites

To interact with our gRPC services from a React + Vite application, you'll need:

- Node.js 16.x or later
- npm 7.x or later
- The proto definitions for our services (`auth.proto` and `chat.proto`)
- A proxy that supports gRPC-Web (like Envoy)
- Protocol Buffers compiler (protoc) and the gRPC-Web plugin

## Proto Definitions

### Auth Service

```protobuf
service AuthService {
  rpc Register (RegisterRequest) returns (RegisterResponse) {}
  rpc Login (LoginRequest) returns (LoginResponse) {}
  rpc ValidateToken (ValidateTokenRequest) returns (ValidateTokenResponse) {}
}

message RegisterRequest {
  string username = 1;
  string email = 2;
  string password = 3;
}

message RegisterResponse {
  bool success = 1;
  string message = 2;
  int64 user_id = 3;
}

message LoginRequest {
  string username = 1;
  string password = 2;
}

message LoginResponse {
  bool success = 1;
  string message = 2;
  string token = 3;
  int64 user_id = 4;
}

message ValidateTokenRequest {
  string token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  int64 user_id = 2;
  string username = 3;
}
```

### Chat Service

```protobuf
service ChatService {
  // Room management
  rpc CreateRoom (CreateRoomRequest) returns (RoomResponse) {}
  rpc GetRooms (GetRoomsRequest) returns (GetRoomsResponse) {}
  rpc JoinRoom (JoinRoomRequest) returns (JoinRoomResponse) {}
  rpc LeaveRoom (LeaveRoomRequest) returns (LeaveRoomResponse) {}

  // Messaging
  rpc SendMessage (SendMessageRequest) returns (SendMessageResponse) {}
  rpc GetRoomMessages (GetRoomMessagesRequest) returns (GetRoomMessagesResponse) {}
  rpc StreamRoomMessages (StreamRoomMessagesRequest) returns (stream MessageResponse) {}
}

// Room related messages
message CreateRoomRequest {
  string name = 1;
  string description = 2;
  int64 creator_id = 3;
  bool is_private = 4;
}

message RoomResponse {
  int64 id = 1;
  string name = 2;
  string description = 3;
  int64 creator_id = 4;
  bool is_private = 5;
  string created_at = 6;
}

message GetRoomsRequest {
  int64 user_id = 1;
  bool include_private = 2;
  int64 limit = 3;
  int64 offset = 4;
}

message GetRoomsResponse {
  repeated RoomResponse rooms = 1;
}

// Message related messages
message SendMessageRequest {
  string content = 1;
  int64 sender_id = 2;
  int64 room_id = 3;
}

message MessageResponse {
  int64 id = 1;
  string content = 2;
  int64 sender_id = 3;
  int64 room_id = 4;
  string sender_name = 5;
  string timestamp = 6;
}
```

## Authentication Flow

### User Registration

1. Call the `Register` method with username, email, and password
2. Check the `success` field in the response to determine if registration was successful
3. If successful, the response will include a `user_id`

Example:

```javascript
// Using a hypothetical gRPC-Web client
const request = new RegisterRequest();
request.setUsername('newuser');
request.setEmail('user@example.com');
request.setPassword('securepassword');

authClient.register(request, {}, (err, response) => {
  if (err) {
    console.error('Registration error:', err);
    return;
  }

  if (response.getSuccess()) {
    console.log('Registration successful! User ID:', response.getUserId());
  } else {
    console.error('Registration failed:', response.getMessage());
  }
});
```

### User Login

1. Call the `Login` method with username and password
2. Check the `success` field in the response to determine if login was successful
3. If successful, the response will include a JWT `token` that should be stored securely (e.g., in localStorage)
4. This token will be used for authenticated requests

Example:

```javascript
const request = new LoginRequest();
request.setUsername('existinguser');
request.setPassword('correctpassword');

authClient.login(request, {}, (err, response) => {
  if (err) {
    console.error('Login error:', err);
    return;
  }

  if (response.getSuccess()) {
    console.log('Login successful! User ID:', response.getUserId());
    // Store the token securely
    localStorage.setItem('auth_token', response.getToken());
  } else {
    console.error('Login failed:', response.getMessage());
  }
});
```

### Token Validation

1. Call the `ValidateToken` method with the stored token
2. Check the `valid` field in the response to determine if the token is valid
3. If valid, the response will include the `user_id` and `username`

Example:

```javascript
const token = localStorage.getItem('auth_token');
if (!token) {
  console.error('No token found');
  return;
}

const request = new ValidateTokenRequest();
request.setToken(token);

authClient.validateToken(request, {}, (err, response) => {
  if (err) {
    console.error('Token validation error:', err);
    return;
  }

  if (response.getValid()) {
    console.log('Token is valid!');
    console.log('User ID:', response.getUserId());
    console.log('Username:', response.getUsername());
  } else {
    console.error('Token is invalid');
    // Clear the invalid token
    localStorage.removeItem('auth_token');
  }
});
```

## Making Authenticated Requests

To make authenticated requests to the ChatService, you need to include the JWT token in the request metadata.

Example for creating a room:

```javascript
const token = localStorage.getItem('auth_token');
if (!token) {
  console.error('No token found');
  return;
}

const request = new CreateRoomRequest();
request.setName('My Chat Room');
request.setDescription('A room for discussing frontend integration');
request.setCreatorId(userId); // Use the user ID from login response
request.setIsPrivate(false);

// Create metadata with the authorization token
const metadata = {'authorization': 'Bearer ' + token};

// Include metadata in the request
chatClient.createRoom(request, metadata, (err, response) => {
  if (err) {
    console.error('Error creating room:', err);
    return;
  }

  console.log('Room created successfully!');
  console.log('Room ID:', response.getId());
  console.log('Room Name:', response.getName());
});
```

## Error Handling

Common errors you might encounter:

1. **Authentication Errors**:
   - Invalid credentials during login
   - Expired or invalid token
   - Missing token for authenticated endpoints

2. **Network Errors**:
   - Connection issues to the gRPC server
   - Timeout errors

Always implement proper error handling in your frontend application to provide a good user experience.

## Security Considerations

1. **Token Storage**: Store JWT tokens securely, preferably in memory or secure storage mechanisms
2. **HTTPS**: Ensure all communication is over HTTPS
3. **Token Expiration**: Handle token expiration gracefully by redirecting to the login page
4. **Sensitive Data**: Never log sensitive data like passwords or tokens

## Setting Up a React + Vite Project

### 1. Create a New Vite Project with React

```bash
npm create vite@latest my-chat-app -- --template react
cd my-chat-app
npm install
```

### 2. Install Required Dependencies

```bash
npm install google-protobuf grpc-web
```

### 3. Set Up Directory Structure for Proto Files

```bash
mkdir -p src/proto src/generated
```

### 4. Copy Proto Files

Copy the `auth.proto` and `chat.proto` files to the `src/proto` directory.

### 5. Generate JavaScript Client Code

```bash
protoc -I=src/proto src/proto/auth.proto src/proto/chat.proto \
  --js_out=import_style=commonjs:src/generated \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/generated
```

### 6. Create a Module for the Generated Files

Create a file `src/generated/index.js` to export all the generated clients:

```javascript
// Export Auth service
export { AuthServiceClient } from './auth_grpc_web_pb';
export { RegisterRequest, RegisterResponse, LoginRequest, LoginResponse, ValidateTokenRequest, ValidateTokenResponse } from './auth_pb';

// Export Chat service
export { ChatServiceClient } from './chat_grpc_web_pb';
export {
  CreateRoomRequest, RoomResponse, GetRoomsRequest, GetRoomsResponse,
  JoinRoomRequest, JoinRoomResponse, LeaveRoomRequest, LeaveRoomResponse,
  SendMessageRequest, SendMessageResponse, GetRoomMessagesRequest, GetRoomMessagesResponse,
  StreamRoomMessagesRequest, MessageResponse
} from './chat_pb';
```

## Example Integration with React and Vite

Here's a simplified example of how you might integrate with our gRPC services in a React application using Vite:

### Authentication Context (src/contexts/AuthContext.jsx)

```jsx
import { createContext, useState, useContext, useEffect } from 'react';
import { AuthServiceClient, ValidateTokenRequest, LoginRequest, RegisterRequest } from '../generated';

// Create the gRPC client
const authClient = new AuthServiceClient('http://localhost:8080');

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    validateToken();
  }, []);

  const validateToken = async () => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setLoading(false);
      return;
    }

    try {
      const request = new ValidateTokenRequest();
      request.setToken(token);

      authClient.validateToken(request, {}, (err, response) => {
        if (err || !response.getValid()) {
          localStorage.removeItem('auth_token');
          setUser(null);
        } else {
          setUser({
            id: response.getUserId(),
            username: response.getUsername(),
            token
          });
        }
        setLoading(false);
      });
    } catch (error) {
      console.error('Token validation error:', error);
      localStorage.removeItem('auth_token');
      setUser(null);
      setLoading(false);
    }
  };

  const login = async (username, password) => {
    return new Promise((resolve, reject) => {
      const request = new LoginRequest();
      request.setUsername(username);
      request.setPassword(password);

      authClient.login(request, {}, (err, response) => {
        if (err) {
          reject(err);
          return;
        }

        if (!response.getSuccess()) {
          reject(new Error(response.getMessage()));
          return;
        }

        const token = response.getToken();
        localStorage.setItem('auth_token', token);

        const userData = {
          id: response.getUserId(),
          username,
          token
        };

        setUser(userData);
        resolve(userData);
      });
    });
  };

  const register = async (username, email, password) => {
    return new Promise((resolve, reject) => {
      const request = new RegisterRequest();
      request.setUsername(username);
      request.setEmail(email);
      request.setPassword(password);

      authClient.register(request, {}, (err, response) => {
        if (err) {
          reject(err);
          return;
        }

        if (!response.getSuccess()) {
          reject(new Error(response.getMessage()));
          return;
        }

        resolve({
          success: true,
          userId: response.getUserId(),
          message: response.getMessage()
        });
      });
    });
  };

  const logout = () => {
    localStorage.removeItem('auth_token');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
```

### Chat Service Hook (src/hooks/useChat.js)

```javascript
import { useState, useEffect, useCallback } from 'react';
import {
  ChatServiceClient, CreateRoomRequest, GetRoomsRequest, JoinRoomRequest,
  SendMessageRequest, GetRoomMessagesRequest, StreamRoomMessagesRequest
} from '../generated';

// Create the gRPC client
const chatClient = new ChatServiceClient('http://localhost:8080');

export function useChat(user) {
  const [rooms, setRooms] = useState([]);
  const [currentRoom, setCurrentRoom] = useState(null);
  const [messages, setMessages] = useState([]);
  const [stream, setStream] = useState(null);

  // Get all available rooms
  const getRooms = useCallback(async () => {
    if (!user) return;

    const request = new GetRoomsRequest();
    request.setUserId(user.id);
    request.setIncludePrivate(true);
    request.setLimit(50);
    request.setOffset(0);

    const metadata = {'authorization': 'Bearer ' + user.token};

    return new Promise((resolve, reject) => {
      chatClient.getRooms(request, metadata, (err, response) => {
        if (err) {
          reject(err);
          return;
        }

        const roomsList = response.getRoomsList();
        setRooms(roomsList.map(room => ({
          id: room.getId(),
          name: room.getName(),
          description: room.getDescription(),
          creatorId: room.getCreatorId(),
          isPrivate: room.getIsPrivate(),
          createdAt: room.getCreatedAt()
        })));
        resolve(roomsList);
      });
    });
  }, [user]);

  // Create a new room
  const createRoom = useCallback(async (name, description, isPrivate = false) => {
    if (!user) return;

    const request = new CreateRoomRequest();
    request.setName(name);
    request.setDescription(description);
    request.setCreatorId(user.id);
    request.setIsPrivate(isPrivate);

    const metadata = {'authorization': 'Bearer ' + user.token};

    return new Promise((resolve, reject) => {
      chatClient.createRoom(request, metadata, (err, response) => {
        if (err) {
          reject(err);
          return;
        }

        const newRoom = {
          id: response.getId(),
          name: response.getName(),
          description: response.getDescription(),
          creatorId: response.getCreatorId(),
          isPrivate: response.getIsPrivate(),
          createdAt: response.getCreatedAt()
        };

        setRooms(prev => [...prev, newRoom]);
        resolve(newRoom);
      });
    });
  }, [user]);

  // Join a room
  const joinRoom = useCallback(async (roomId) => {
    if (!user) return;

    const request = new JoinRoomRequest();
    request.setRoomId(roomId);
    request.setUserId(user.id);

    const metadata = {'authorization': 'Bearer ' + user.token};

    return new Promise((resolve, reject) => {
      chatClient.joinRoom(request, metadata, (err, response) => {
        if (err) {
          reject(err);
          return;
        }

        if (!response.getSuccess()) {
          reject(new Error(response.getMessage()));
          return;
        }

        const room = response.getRoom();
        const roomData = {
          id: room.getId(),
          name: room.getName(),
          description: room.getDescription(),
          creatorId: room.getCreatorId(),
          isPrivate: room.getIsPrivate(),
          createdAt: room.getCreatedAt()
        };

        setCurrentRoom(roomData);
        resolve(roomData);

        // Load messages for this room
        loadMessages(roomId);
        // Start streaming messages
        startMessageStream(roomId);
      });
    });
  }, [user]);

  // Load messages for a room
  const loadMessages = useCallback(async (roomId) => {
    if (!user) return;

    const request = new GetRoomMessagesRequest();
    request.setRoomId(roomId);
    request.setUserId(user.id);
    request.setLimit(50);
    request.setOffset(0);

    const metadata = {'authorization': 'Bearer ' + user.token};

    return new Promise((resolve, reject) => {
      chatClient.getRoomMessages(request, metadata, (err, response) => {
        if (err) {
          reject(err);
          return;
        }

        const messagesList = response.getMessagesList();
        const formattedMessages = messagesList.map(msg => ({
          id: msg.getId(),
          content: msg.getContent(),
          senderId: msg.getSenderId(),
          roomId: msg.getRoomId(),
          senderName: msg.getSenderName(),
          timestamp: msg.getTimestamp()
        }));

        setMessages(formattedMessages);
        resolve(formattedMessages);
      });
    });
  }, [user]);

  // Start streaming messages for a room
  const startMessageStream = useCallback((roomId) => {
    if (!user || !roomId) return;

    // Close existing stream if any
    if (stream) {
      stream.cancel();
    }

    const request = new StreamRoomMessagesRequest();
    request.setRoomId(roomId);
    request.setUserId(user.id);

    const metadata = {'authorization': 'Bearer ' + user.token};

    const newStream = chatClient.streamRoomMessages(request, metadata);

    newStream.on('data', (response) => {
      const newMessage = {
        id: response.getId(),
        content: response.getContent(),
        senderId: response.getSenderId(),
        roomId: response.getRoomId(),
        senderName: response.getSenderName(),
        timestamp: response.getTimestamp()
      };

      setMessages(prev => [...prev, newMessage]);
    });

    newStream.on('error', (err) => {
      console.error('Stream error:', err);
    });

    newStream.on('end', () => {
      console.log('Stream ended');
    });

    setStream(newStream);

    return () => {
      if (newStream) {
        newStream.cancel();
      }
    };
  }, [user, stream]);

  // Send a message to the current room
  const sendMessage = useCallback(async (content) => {
    if (!user || !currentRoom) return;

    const request = new SendMessageRequest();
    request.setContent(content);
    request.setSenderId(user.id);
    request.setRoomId(currentRoom.id);

    const metadata = {'authorization': 'Bearer ' + user.token};

    return new Promise((resolve, reject) => {
      chatClient.sendMessage(request, metadata, (err, response) => {
        if (err) {
          reject(err);
          return;
        }

        if (!response.getSuccess()) {
          reject(new Error(response.getMessage()));
          return;
        }

        resolve({
          success: true,
          messageId: response.getMessageId()
        });
      });
    });
  }, [user, currentRoom]);

  // Clean up stream on unmount
  useEffect(() => {
    return () => {
      if (stream) {
        stream.cancel();
      }
    };
  }, [stream]);

  return {
    rooms,
    currentRoom,
    messages,
    getRooms,
    createRoom,
    joinRoom,
    sendMessage
  };
}
```

### Main App Component (src/App.jsx)

```jsx
import { useState } from 'react';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import LoginForm from './components/LoginForm';
import RegisterForm from './components/RegisterForm';
import ChatRoom from './components/ChatRoom';
import './App.css';

function AppContent() {
  const { user, loading, logout } = useAuth();
  const [showRegister, setShowRegister] = useState(false);

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="app-container">
      <header>
        <h1>gRPC Chat Room</h1>
        {user && (
          <div className="user-info">
            <span>Welcome, {user.username}!</span>
            <button onClick={logout}>Logout</button>
          </div>
        )}
      </header>

      <main>
        {user ? (
          <ChatRoom />
        ) : (
          <div className="auth-container">
            {showRegister ? (
              <>
                <RegisterForm />
                <p>
                  Already have an account?{' '}
                  <button onClick={() => setShowRegister(false)}>Login</button>
                </p>
              </>
            ) : (
              <>
                <LoginForm />
                <p>
                  Don't have an account?{' '}
                  <button onClick={() => setShowRegister(true)}>Register</button>
                </p>
              </>
            )}
          </div>
        )}
      </main>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
```

## Next Steps

1. Set up a gRPC-Web proxy (like Envoy) to communicate with the backend
2. Create the UI components for login, registration, and chat functionality
3. Implement error handling and user feedback
4. Add styling to improve the user experience
5. Consider adding features like typing indicators, read receipts, or file sharing

For more detailed information, refer to the [gRPC-Web documentation](https://github.com/grpc/grpc-web) and the [Vite documentation](https://vitejs.dev/guide/).
