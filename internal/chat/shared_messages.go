package chat

import (
	"log"
	"sync"

	pb "grpc-messenger-core/proto/chat"
)

// Global shared messages store
var (
	// Global message store
	GlobalMessages     = make(map[int64][]*pb.MessageResponse)
	GlobalMessageMutex sync.RWMutex
	sharedLogger       *log.Logger
)

// AddMessage adds a message to the global message store
func AddMessage(message *pb.MessageResponse) {
	GlobalMessageMutex.Lock()
	defer GlobalMessageMutex.Unlock()

	// Initialize the slice if it doesn't exist
	if GlobalMessages[message.RoomId] == nil {
		GlobalMessages[message.RoomId] = make([]*pb.MessageResponse, 0)
	}

	// Add the message
	GlobalMessages[message.RoomId] = append(GlobalMessages[message.RoomId], message)

	if sharedLogger != nil {
		sharedLogger.Printf("Added message to global store for room %d. Total messages: %d",
			message.RoomId, len(GlobalMessages[message.RoomId]))
	}
}

// GetMessages returns all messages for a room
func GetMessages(roomID int64) []*pb.MessageResponse {
	GlobalMessageMutex.RLock()
	defer GlobalMessageMutex.RUnlock()

	// Return a copy of the messages
	messages := make([]*pb.MessageResponse, 0)
	if roomMessages, ok := GlobalMessages[roomID]; ok {
		messages = append(messages, roomMessages...)
	}

	if sharedLogger != nil {
		sharedLogger.Printf("Retrieved %d messages from global store for room %d",
			len(messages), roomID)
	}

	return messages
}
