# Setting Up gRPC-Web for Frontend Development

This guide explains how to set up gRPC-Web to connect your frontend application with our gRPC backend services.

## What is gRPC-Web?

gRPC-Web is a JavaScript client library that enables web applications to directly communicate with gRPC services. It allows you to generate JavaScript client stubs from your proto definitions and make gRPC calls from the browser.

## Prerequisites

- Node.js and npm/yarn
- Access to the proto definitions for our services
- A proxy that supports gRPC-Web (like Envoy)

## Installation Steps

### 1. Install Required Tools

First, install the Protocol Buffers compiler (protoc) and the gRPC-Web plugin:

#### For macOS (using Homebrew):

```bash
brew install protobuf
npm install -g protoc-gen-grpc-web
```

#### For Linux:

```bash
# Install protoc
PROTOC_ZIP=protoc-3.14.0-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP

# Install protoc-gen-grpc-web
npm install -g protoc-gen-grpc-web
```

#### For Windows:

Download the protoc compiler from [GitHub](https://github.com/protocolbuffers/protobuf/releases) and add it to your PATH.

```bash
npm install -g protoc-gen-grpc-web
```

### 2. Set Up Your Frontend Project

#### For React:

```bash
npx create-react-app my-grpc-app
cd my-grpc-app
npm install google-protobuf grpc-web
```

#### For Vue.js:

```bash
vue create my-grpc-app
cd my-grpc-app
npm install google-protobuf grpc-web
```

#### For Angular:

```bash
ng new my-grpc-app
cd my-grpc-app
npm install google-protobuf grpc-web
```

### 3. Generate JavaScript Client Code

Create a directory structure for your proto files and generated code:

```bash
mkdir -p src/proto
```

Copy the proto files (`auth.proto` and `hello.proto`) to the `src/proto` directory.

Generate the JavaScript client code:

```bash
protoc -I=src/proto src/proto/auth.proto src/proto/hello.proto \
  --js_out=import_style=commonjs:src/generated \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/generated
```

This will generate the following files in the `src/generated` directory:

- `auth_pb.js`: Contains message classes for Auth service
- `auth_grpc_web_pb.js`: Contains client stubs for Auth service
- `hello_pb.js`: Contains message classes for Hello service
- `hello_grpc_web_pb.js`: Contains client stubs for Hello service

### 4. Set Up Envoy Proxy

gRPC-Web requires a proxy to translate between the gRPC-Web protocol and the gRPC protocol. Envoy is the recommended proxy.

Create a file named `envoy.yaml`:

```yaml
admin:
  access_log_path: /tmp/admin_access.log
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }

static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address: { address: 0.0.0.0, port_value: 8080 }
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
          codec_type: auto
          stat_prefix: ingress_http
          route_config:
            name: local_route
            virtual_hosts:
            - name: local_service
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route:
                  cluster: grpc_service
                  timeout: 0s
                  max_stream_duration:
                    grpc_timeout_header_max: 0s
              cors:
                allow_origin_string_match:
                - prefix: "*"
                allow_methods: GET, PUT, DELETE, POST, OPTIONS
                allow_headers: keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout,authorization
                max_age: "1728000"
                expose_headers: custom-header-1,grpc-status,grpc-message
          http_filters:
          - name: envoy.filters.http.grpc_web
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.grpc_web.v3.GrpcWeb
          - name: envoy.filters.http.cors
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.cors.v3.Cors
          - name: envoy.filters.http.router
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router
  clusters:
  - name: grpc_service
    connect_timeout: 0.25s
    type: logical_dns
    http2_protocol_options: {}
    lb_policy: round_robin
    load_assignment:
      cluster_name: cluster_0
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: host.docker.internal  # Use 'localhost' if not using Docker
                port_value: 50051  # The port your gRPC server is running on
```

Run Envoy using Docker:

```bash
docker run --name envoy -d \
  -v $(pwd)/envoy.yaml:/etc/envoy/envoy.yaml:ro \
  -p 8080:8080 -p 9901:9901 \
  envoyproxy/envoy:v1.20.0
```

### 5. Using the Generated Client Code

Here's how to use the generated client code in your frontend application:

#### React Example:

```jsx
import React, { useState } from 'react';
import { AuthServiceClient } from './generated/auth_grpc_web_pb';
import { LoginRequest } from './generated/auth_pb';

// Create the client
const client = new AuthServiceClient('http://localhost:8080');

function LoginForm() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    
    const request = new LoginRequest();
    request.setUsername(username);
    request.setPassword(password);
    
    client.login(request, {}, (err, response) => {
      if (err) {
        setMessage('Error: ' + err.message);
        return;
      }
      
      if (response.getSuccess()) {
        setMessage('Login successful!');
        // Store the token
        localStorage.setItem('auth_token', response.getToken());
      } else {
        setMessage('Login failed: ' + response.getMessage());
      }
    });
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Username:</label>
          <input 
            type="text" 
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div>
          <label>Password:</label>
          <input 
            type="password" 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit">Login</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default LoginForm;
```

#### Vue.js Example:

```vue
<template>
  <div>
    <h2>Login</h2>
    <form @submit.prevent="handleSubmit">
      <div>
        <label>Username:</label>
        <input type="text" v-model="username" />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" v-model="password" />
      </div>
      <button type="submit">Login</button>
    </form>
    <p v-if="message">{{ message }}</p>
  </div>
</template>

<script>
import { AuthServiceClient } from '../generated/auth_grpc_web_pb';
import { LoginRequest } from '../generated/auth_pb';

// Create the client
const client = new AuthServiceClient('http://localhost:8080');

export default {
  data() {
    return {
      username: '',
      password: '',
      message: ''
    };
  },
  methods: {
    handleSubmit() {
      const request = new LoginRequest();
      request.setUsername(this.username);
      request.setPassword(this.password);
      
      client.login(request, {}, (err, response) => {
        if (err) {
          this.message = 'Error: ' + err.message;
          return;
        }
        
        if (response.getSuccess()) {
          this.message = 'Login successful!';
          // Store the token
          localStorage.setItem('auth_token', response.getToken());
        } else {
          this.message = 'Login failed: ' + response.getMessage();
        }
      });
    }
  }
};
</script>
```

#### Angular Example:

```typescript
// login.component.ts
import { Component } from '@angular/core';
import { AuthServiceClient } from '../generated/auth_grpc_web_pb';
import { LoginRequest } from '../generated/auth_pb';

// Create the client
const client = new AuthServiceClient('http://localhost:8080');

@Component({
  selector: 'app-login',
  template: `
    <div>
      <h2>Login</h2>
      <form (ngSubmit)="handleSubmit()">
        <div>
          <label>Username:</label>
          <input type="text" [(ngModel)]="username" name="username" />
        </div>
        <div>
          <label>Password:</label>
          <input type="password" [(ngModel)]="password" name="password" />
        </div>
        <button type="submit">Login</button>
      </form>
      <p *ngIf="message">{{ message }}</p>
    </div>
  `,
  styles: []
})
export class LoginComponent {
  username = '';
  password = '';
  message = '';

  handleSubmit() {
    const request = new LoginRequest();
    request.setUsername(this.username);
    request.setPassword(this.password);
    
    client.login(request, {}, (err, response) => {
      if (err) {
        this.message = 'Error: ' + err.message;
        return;
      }
      
      if (response.getSuccess()) {
        this.message = 'Login successful!';
        // Store the token
        localStorage.setItem('auth_token', response.getToken());
      } else {
        this.message = 'Login failed: ' + response.getMessage();
      }
    });
  }
}
```

## Making Authenticated Requests

To make authenticated requests, include the JWT token in the metadata:

```javascript
const token = localStorage.getItem('auth_token');
if (!token) {
  console.error('No token found');
  return;
}

// Create metadata with the authorization token
const metadata = {'authorization': 'Bearer ' + token};

// Include metadata in the request
client.someMethod(request, metadata, (err, response) => {
  // Handle response
});
```

## Troubleshooting

### Common Issues:

1. **CORS Errors**: Ensure your Envoy proxy is configured to handle CORS correctly.

2. **Connection Refused**: Make sure your gRPC server is running and the Envoy proxy is configured with the correct address and port.

3. **Invalid Token**: If you're getting authentication errors, check that your token is valid and properly formatted in the metadata.

4. **Generated Code Issues**: If you update your proto files, regenerate the client code to ensure it's up to date.

## Additional Resources

- [gRPC-Web GitHub Repository](https://github.com/grpc/grpc-web)
- [Envoy Proxy Documentation](https://www.envoyproxy.io/docs/envoy/latest/)
- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers/docs/overview)

## Support

If you encounter any issues with the gRPC-Web setup or integration, please contact our backend team for assistance.
