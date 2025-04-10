# Frontend Integration Guide

This document provides guidelines for frontend developers to integrate with our gRPC authentication and Hello services.

## Overview

Our backend provides the following gRPC services:

1. **AuthService**: Handles user registration, login, and token validation
2. **HelloService**: A simple greeting service that supports authenticated requests

## Prerequisites

To interact with our gRPC services from a frontend application, you'll need:

- A gRPC-Web compatible client library for your frontend framework
- The proto definitions for our services
- A proxy that supports gRPC-Web (like Envoy)

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

### Hello Service

```protobuf
service HelloService {
  rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
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

To make authenticated requests to other services (like the HelloService), you need to include the JWT token in the request metadata.

Example:

```javascript
const token = localStorage.getItem('auth_token');
if (!token) {
  console.error('No token found');
  return;
}

const request = new HelloRequest();
request.setName('World');

// Create metadata with the authorization token
const metadata = {'authorization': 'Bearer ' + token};

// Include metadata in the request
helloClient.sayHello(request, metadata, (err, response) => {
  if (err) {
    console.error('Error calling SayHello:', err);
    return;
  }
  
  console.log('Greeting:', response.getMessage());
  // For authenticated users, the response will include their username
  // e.g., "Hello World (Authenticated as username)"
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

## Example Integration with React

Here's a simplified example of how you might integrate with our gRPC services in a React application:

```jsx
import React, { useState, useEffect } from 'react';
import { AuthServiceClient } from './generated/AuthServiceClientPb';
import { HelloServiceClient } from './generated/HelloServiceClientPb';
import { LoginRequest, RegisterRequest, ValidateTokenRequest } from './generated/auth_pb';
import { HelloRequest } from './generated/hello_pb';

// Create clients
const authClient = new AuthServiceClient('http://localhost:8080');
const helloClient = new HelloServiceClient('http://localhost:8080');

function App() {
  const [user, setUser] = useState(null);
  const [greeting, setGreeting] = useState('');
  const [loginForm, setLoginForm] = useState({ username: '', password: '' });

  useEffect(() => {
    // Check if user is already logged in
    validateToken();
  }, []);

  const validateToken = () => {
    const token = localStorage.getItem('auth_token');
    if (!token) return;

    const request = new ValidateTokenRequest();
    request.setToken(token);

    authClient.validateToken(request, {}, (err, response) => {
      if (err || !response.getValid()) {
        localStorage.removeItem('auth_token');
        return;
      }
      
      setUser({
        id: response.getUserId(),
        username: response.getUsername(),
        token
      });
    });
  };

  const handleLogin = (e) => {
    e.preventDefault();
    
    const request = new LoginRequest();
    request.setUsername(loginForm.username);
    request.setPassword(loginForm.password);

    authClient.login(request, {}, (err, response) => {
      if (err || !response.getSuccess()) {
        alert(response ? response.getMessage() : 'Login failed');
        return;
      }
      
      const token = response.getToken();
      localStorage.setItem('auth_token', token);
      
      setUser({
        id: response.getUserId(),
        username: loginForm.username,
        token
      });
    });
  };

  const handleLogout = () => {
    localStorage.removeItem('auth_token');
    setUser(null);
    setGreeting('');
  };

  const fetchGreeting = () => {
    const request = new HelloRequest();
    request.setName('World');

    const metadata = user ? {'authorization': 'Bearer ' + user.token} : {};

    helloClient.sayHello(request, metadata, (err, response) => {
      if (err) {
        alert('Error fetching greeting');
        return;
      }
      
      setGreeting(response.getMessage());
    });
  };

  return (
    <div className="App">
      <h1>gRPC Frontend Example</h1>
      
      {user ? (
        <div>
          <h2>Welcome, {user.username}!</h2>
          <button onClick={handleLogout}>Logout</button>
          <hr />
          <button onClick={fetchGreeting}>Get Greeting</button>
          {greeting && <p>{greeting}</p>}
        </div>
      ) : (
        <form onSubmit={handleLogin}>
          <h2>Login</h2>
          <div>
            <label>Username:</label>
            <input 
              type="text" 
              value={loginForm.username}
              onChange={(e) => setLoginForm({...loginForm, username: e.target.value})}
            />
          </div>
          <div>
            <label>Password:</label>
            <input 
              type="password" 
              value={loginForm.password}
              onChange={(e) => setLoginForm({...loginForm, password: e.target.value})}
            />
          </div>
          <button type="submit">Login</button>
        </form>
      )}
    </div>
  );
}

export default App;
```

## Next Steps

1. Generate client code from the proto definitions for your specific frontend framework
2. Set up a gRPC-Web proxy (like Envoy) to communicate with the backend
3. Implement the authentication flow in your application
4. Add error handling and user feedback
5. Secure token storage and management

For more detailed information, refer to the gRPC-Web documentation for your specific frontend framework.
