package chat

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	pb "grpc-messenger-core/proto/chat"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// Helper function to create a context with auth token
func createAuthContext() context.Context {
	// In a real test, this would be a valid token
	// For now, we'll use a mock token that the service will accept in test mode
	return metadata.NewIncomingContext(
		context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer mock.token.for.testing"}),
	)
}

func TestSendMessage(t *testing.T) {
	// Create a new chat service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewChatService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.SendMessageRequest
		expectError bool
	}{
		{
			name: "Valid message",
			request: &pb.SendMessageRequest{
				Content:  "Hello, world!",
				SenderId: 1,
				RoomId:   1,
			},
			expectError: false,
		},
		{
			name: "Empty content",
			request: &pb.SendMessageRequest{
				Content:  "",
				SenderId: 1,
				RoomId:   1,
			},
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with auth token
			ctx := createAuthContext()

			// Call SendMessage
			response, err := service.SendMessage(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				assert.NotEmpty(t, response.Message)
				assert.Greater(t, response.MessageId, int64(0))
			}
		})
	}
}

func TestGetRoomMessages(t *testing.T) {
	// Create a new chat service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewChatService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.GetRoomMessagesRequest
		expectError bool
	}{
		{
			name: "Valid request",
			request: &pb.GetRoomMessagesRequest{
				RoomId: 1,
				UserId: 1,
				Limit:  10,
				Offset: 0,
			},
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with auth token
			ctx := createAuthContext()

			// Call GetRoomMessages
			response, err := service.GetRoomMessages(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Messages)
				
				// Check the first message
				firstMessage := response.Messages[0]
				assert.Greater(t, firstMessage.Id, int64(0))
				assert.NotEmpty(t, firstMessage.Content)
				assert.Greater(t, firstMessage.SenderId, int64(0))
				assert.Equal(t, tc.request.RoomId, firstMessage.RoomId)
				assert.NotEmpty(t, firstMessage.SenderName)
				assert.NotEmpty(t, firstMessage.Timestamp)
				
				// Verify timestamp format
				_, err := time.Parse(time.RFC3339, firstMessage.Timestamp)
				assert.NoError(t, err)
			}
		})
	}
}

func TestStreamRoomMessages(t *testing.T) {
	// This is a simplified test for the streaming functionality
	// In a real test, you would use a mock stream and verify the messages sent
	
	// Create a new chat service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewChatService(nil, logger)
	
	// Create a request
	request := &pb.StreamRoomMessagesRequest{
		RoomId: 1,
		UserId: 1,
	}
	
	// Create a mock stream
	mockStream := &mockChatServiceStreamRoomMessagesServer{
		ctx: createAuthContext(),
	}
	
	// Call StreamRoomMessages in a goroutine
	errChan := make(chan error)
	go func() {
		err := service.StreamRoomMessages(request, mockStream)
		errChan <- err
	}()
	
	// Wait for a short time to allow some messages to be sent
	select {
	case err := <-errChan:
		// If the function returns, check for errors
		assert.NoError(t, err)
	case <-time.After(100 * time.Millisecond):
		// If it doesn't return, that's expected (it's a long-running stream)
		// Cancel the context to stop the stream
		mockStream.cancel()
	}
}

// Mock implementation of ChatService_StreamRoomMessagesServer
type mockChatServiceStreamRoomMessagesServer struct {
	pb.UnimplementedChatServiceServer
	ctx    context.Context
	cancel func()
	sent   []*pb.MessageResponse
}

func (m *mockChatServiceStreamRoomMessagesServer) Send(msg *pb.MessageResponse) error {
	m.sent = append(m.sent, msg)
	return nil
}

func (m *mockChatServiceStreamRoomMessagesServer) Context() context.Context {
	if m.ctx == nil {
		m.ctx, m.cancel = context.WithCancel(createAuthContext())
	}
	return m.ctx
}
