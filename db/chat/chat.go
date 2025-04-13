package chat

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

// Message represents a chat message in the database
type Message struct {
	ID         int64
	Content    string
	SenderID   int64
	RoomID     int64
	SenderName string
	Timestamp  time.Time
}

// Repository handles database operations for chat
type Repository struct {
	db *sql.DB

	// For real-time messaging
	roomSubscriptions     map[int64][]chan Message
	roomSubscriptionMutex sync.RWMutex
}

// NewRepository creates a new chat repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db:                db,
		roomSubscriptions: make(map[int64][]chan Message),
	}
}

// SaveMessage saves a message to the database and notifies subscribers
func (r *Repository) SaveMessage(ctx context.Context, content string, senderID, roomID int64) (int64, error) {
	var messageID int64
	var senderName string
	var timestamp time.Time

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Get sender name
	err = tx.QueryRowContext(ctx, `SELECT username FROM users WHERE id = $1`, senderID).Scan(&senderName)
	if err != nil {
		return 0, err
	}

	// Insert message
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO messages (content, sender_id, room_id) VALUES ($1, $2, $3) RETURNING id, created_at`,
		content, senderID, roomID,
	).Scan(&messageID, &timestamp)
	if err != nil {
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// Notify subscribers
	message := Message{
		ID:         messageID,
		Content:    content,
		SenderID:   senderID,
		RoomID:     roomID,
		SenderName: senderName,
		Timestamp:  timestamp,
	}
	r.NotifyRoomSubscribers(roomID, message)

	return messageID, nil
}

// GetRoomMessages retrieves messages from a room
func (r *Repository) GetRoomMessages(ctx context.Context, roomID, limit, offset int64) ([]Message, error) {
	query := `
		SELECT m.id, m.content, m.sender_id, m.room_id, u.username, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.SenderID, &msg.RoomID, &msg.SenderName, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// IsRoomMember checks if a user is a member of a room
func (r *Repository) IsRoomMember(ctx context.Context, roomID, userID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`
	err := r.db.QueryRowContext(ctx, query, roomID, userID).Scan(&exists)
	return exists, err
}

// SubscribeToRoom subscribes to messages in a room
func (r *Repository) SubscribeToRoom(roomID int64, ch chan Message) {
	r.roomSubscriptionMutex.Lock()
	defer r.roomSubscriptionMutex.Unlock()

	r.roomSubscriptions[roomID] = append(r.roomSubscriptions[roomID], ch)
}

// UnsubscribeFromRoom unsubscribes from messages in a room
func (r *Repository) UnsubscribeFromRoom(roomID int64, ch chan Message) {
	r.roomSubscriptionMutex.Lock()
	defer r.roomSubscriptionMutex.Unlock()

	subs := r.roomSubscriptions[roomID]
	for i, sub := range subs {
		if sub == ch {
			// Remove the channel from the slice
			r.roomSubscriptions[roomID] = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	// If no more subscribers, remove the room from the map
	if len(r.roomSubscriptions[roomID]) == 0 {
		delete(r.roomSubscriptions, roomID)
	}
}

// NotifyRoomSubscribers notifies all subscribers of a new message
func (r *Repository) NotifyRoomSubscribers(roomID int64, message Message) {
	r.roomSubscriptionMutex.RLock()
	defer r.roomSubscriptionMutex.RUnlock()

	for _, ch := range r.roomSubscriptions[roomID] {
		// Use non-blocking send to avoid deadlocks
		select {
		case ch <- message:
			// Message sent successfully
		default:
			// Channel is full or closed, skip
		}
	}
}
