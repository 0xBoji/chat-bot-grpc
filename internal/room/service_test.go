package room

import (
	"context"
	"log"
	"os"
	"testing"

	pb "grpc-messenger-core/proto/room"

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

func TestCreateRoom(t *testing.T) {
	// Create a new room service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewRoomService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.CreateRoomRequest
		expectError bool
	}{
		{
			name: "Valid room creation",
			request: &pb.CreateRoomRequest{
				Name:        "Test Room",
				Description: "A test room",
				CreatorId:   1,
			},
			expectError: false,
		},
		{
			name: "Empty name",
			request: &pb.CreateRoomRequest{
				Name:        "",
				Description: "A test room",
				CreatorId:   1,
			},
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with auth token
			ctx := createAuthContext()

			// Call CreateRoom
			response, err := service.CreateRoom(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Greater(t, response.Id, int64(0))
				assert.Equal(t, tc.request.Name, response.Name)
				assert.Equal(t, tc.request.Description, response.Description)
				assert.Equal(t, tc.request.CreatorId, response.CreatorId)
			}
		})
	}
}

func TestGetRooms(t *testing.T) {
	// Create a new room service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewRoomService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.GetRoomsRequest
		expectError bool
	}{
		{
			name: "Valid request",
			request: &pb.GetRoomsRequest{
				UserId: 1,
			},
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with auth token
			ctx := createAuthContext()

			// Call GetRooms
			response, err := service.GetRooms(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Rooms)
				
				// Check the first room
				firstRoom := response.Rooms[0]
				assert.Greater(t, firstRoom.Id, int64(0))
				assert.NotEmpty(t, firstRoom.Name)
			}
		})
	}
}

func TestJoinRoom(t *testing.T) {
	// Create a new room service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewRoomService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.JoinRoomRequest
		expectError bool
	}{
		{
			name: "Valid join",
			request: &pb.JoinRoomRequest{
				RoomId: 1,
				UserId: 1,
			},
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with auth token
			ctx := createAuthContext()

			// Call JoinRoom
			response, err := service.JoinRoom(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				assert.NotEmpty(t, response.Message)
			}
		})
	}
}

func TestLeaveRoom(t *testing.T) {
	// Create a new room service with nil DB (will use mock mode)
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := NewRoomService(nil, logger)

	// Test cases
	testCases := []struct {
		name        string
		request     *pb.LeaveRoomRequest
		expectError bool
	}{
		{
			name: "Valid leave",
			request: &pb.LeaveRoomRequest{
				RoomId: 1,
				UserId: 1,
			},
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context with auth token
			ctx := createAuthContext()

			// Call LeaveRoom
			response, err := service.LeaveRoom(ctx, tc.request)

			// Check results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				assert.NotEmpty(t, response.Message)
			}
		})
	}
}
