package auth

import (
	"context"
	"log"
	"os"
	"testing"

	pb "grpc-messenger-core/proto/auth"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestRegister(t *testing.T) {
	// Create a new auth service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewAuthService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.RegisterRequest
		expectError bool
	}{
		{
			name: "Valid registration",
			request: &pb.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Empty username",
			request: &pb.RegisterRequest{
				Username: "",
				Password: "password123",
			},
			expectError: true,
		},
		{
			name: "Empty password",
			request: &pb.RegisterRequest{
				Username: "testuser",
				Password: "",
			},
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context
			ctx := context.Background()

			// Call Register
			response, err := service.Register(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				assert.NotEmpty(t, response.Message)
				assert.Greater(t, response.UserId, int64(0))
			}
		})
	}
}

func TestLogin(t *testing.T) {
	// Create a new auth service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewAuthService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.LoginRequest
		expectError bool
	}{
		{
			name: "Valid login",
			request: &pb.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Empty username",
			request: &pb.LoginRequest{
				Username: "",
				Password: "password123",
			},
			expectError: true,
		},
		{
			name: "Empty password",
			request: &pb.LoginRequest{
				Username: "testuser",
				Password: "",
			},
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context
			ctx := context.Background()

			// Call Login
			response, err := service.Login(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				assert.NotEmpty(t, response.Message)
				assert.NotEmpty(t, response.Token)
				assert.Greater(t, response.UserId, int64(0))
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Create a new auth service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewAuthService(nil, logger)

	// First, get a valid token by logging in
	loginResp, err := service.Login(context.Background(), &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResp.Token)

	// Test cases
	testCases := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       loginResp.Token,
			expectError: false,
		},
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "Invalid token",
			token:       "invalid.token.here",
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with token in metadata
			ctx := metadata.NewIncomingContext(
				context.Background(),
				metadata.New(map[string]string{"authorization": "Bearer " + tc.token}),
			)

			// Call ValidateToken
			response, err := service.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: tc.token})

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Valid)
				assert.Greater(t, response.UserId, int64(0))
				assert.NotEmpty(t, response.Username)
			}
		})
	}
}
