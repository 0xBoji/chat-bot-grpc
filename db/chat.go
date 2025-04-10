package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Room represents a chat room in the database
type Room struct {
	ID          int64
	Name        string
	Description string
	CreatorID   int64
	IsPrivate   bool
	CreatedAt   time.Time
}

// RoomMember represents a user's membership in a room
type RoomMember struct {
	RoomID   int64
	UserID   int64
	JoinedAt time.Time
}

// Message represents a chat message in the database
type Message struct {
	ID        int64
	Content   string
	SenderID  int64
	RoomID    int64
	Timestamp time.Time
}

// CreateTablesIfNotExist creates all necessary tables if they don't exist
func CreateTablesIfNotExist(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Create rooms table
	roomsQuery := `
	CREATE TABLE IF NOT EXISTS rooms (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		creator_id BIGINT NOT NULL REFERENCES users(id),
		is_private BOOLEAN NOT NULL DEFAULT false,
		created_at TIMESTAMPTZ DEFAULT NOW()
	)`

	_, err := db.Exec(roomsQuery)
	if err != nil {
		return fmt.Errorf("error creating rooms table: %w", err)
	}

	// Create room_members table
	membersQuery := `
	CREATE TABLE IF NOT EXISTS room_members (
		room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		joined_at TIMESTAMPTZ DEFAULT NOW(),
		PRIMARY KEY (room_id, user_id)
	)`

	_, err = db.Exec(membersQuery)
	if err != nil {
		return fmt.Errorf("error creating room_members table: %w", err)
	}

	// Create messages table
	messagesQuery := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		sender_id BIGINT NOT NULL REFERENCES users(id),
		room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
		timestamp TIMESTAMPTZ DEFAULT NOW()
	)`

	_, err = db.Exec(messagesQuery)
	if err != nil {
		return fmt.Errorf("error creating messages table: %w", err)
	}

	return nil
}

// Room operations

// CreateRoom creates a new chat room
func CreateRoom(ctx context.Context, db *sql.DB, name, description string, creatorID int64, isPrivate bool) (*Room, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var room Room
	query := `
	INSERT INTO rooms (name, description, creator_id, is_private)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, description, creator_id, is_private, created_at`

	err := db.QueryRowContext(
		ctx, query, name, description, creatorID, isPrivate,
	).Scan(&room.ID, &room.Name, &room.Description, &room.CreatorID, &room.IsPrivate, &room.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("error creating room: %w", err)
	}

	// Automatically add the creator as a member
	_, err = JoinRoom(ctx, db, room.ID, creatorID)
	if err != nil {
		return nil, fmt.Errorf("error adding creator to room: %w", err)
	}

	return &room, nil
}

// GetRooms retrieves rooms that a user can access
func GetRooms(ctx context.Context, db *sql.DB, userID int64, includePrivate bool, limit, offset int64) ([]Room, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var query string
	var args []interface{}

	if includePrivate {
		// Include all public rooms and private rooms the user is a member of
		query = `
		SELECT r.id, r.name, r.description, r.creator_id, r.is_private, r.created_at
		FROM rooms r
		WHERE NOT r.is_private OR r.creator_id = $1 OR EXISTS (
			SELECT 1 FROM room_members rm WHERE rm.room_id = r.id AND rm.user_id = $1
		)
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3`
		args = []interface{}{userID, limit, offset}
	} else {
		// Include only public rooms
		query = `
		SELECT r.id, r.name, r.description, r.creator_id, r.is_private, r.created_at
		FROM rooms r
		WHERE NOT r.is_private
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying rooms: %w", err)
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(
			&room.ID, &room.Name, &room.Description,
			&room.CreatorID, &room.IsPrivate, &room.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning room row: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating room rows: %w", err)
	}

	return rooms, nil
}

// GetRoomByID retrieves a room by its ID
func GetRoomByID(ctx context.Context, db *sql.DB, roomID int64) (*Room, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var room Room
	query := `
	SELECT id, name, description, creator_id, is_private, created_at
	FROM rooms
	WHERE id = $1`

	err := db.QueryRowContext(ctx, query, roomID).Scan(
		&room.ID, &room.Name, &room.Description,
		&room.CreatorID, &room.IsPrivate, &room.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("error querying room: %w", err)
	}

	return &room, nil
}

// JoinRoom adds a user to a room
func JoinRoom(ctx context.Context, db *sql.DB, roomID, userID int64) (*RoomMember, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Check if the room exists and if it's private
	var isPrivate bool
	checkQuery := "SELECT is_private FROM rooms WHERE id = $1"
	err := db.QueryRowContext(ctx, checkQuery, roomID).Scan(&isPrivate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("error checking room: %w", err)
	}

	// If the room is private, check if the user is the creator
	if isPrivate {
		var creatorID int64
		creatorQuery := "SELECT creator_id FROM rooms WHERE id = $1"
		err := db.QueryRowContext(ctx, creatorQuery, roomID).Scan(&creatorID)
		if err != nil {
			return nil, fmt.Errorf("error checking room creator: %w", err)
		}

		// If the user is not the creator, they need an invitation (not implemented here)
		if creatorID != userID {
			// Check if the user is already a member
			var count int
			memberQuery := "SELECT COUNT(*) FROM room_members WHERE room_id = $1 AND user_id = $2"
			err := db.QueryRowContext(ctx, memberQuery, roomID, userID).Scan(&count)
			if err != nil {
				return nil, fmt.Errorf("error checking membership: %w", err)
			}

			if count == 0 {
				return nil, fmt.Errorf("cannot join private room without invitation")
			}
		}
	}

	var member RoomMember
	query := `
	INSERT INTO room_members (room_id, user_id)
	VALUES ($1, $2)
	ON CONFLICT (room_id, user_id) DO NOTHING
	RETURNING room_id, user_id, joined_at`

	err = db.QueryRowContext(ctx, query, roomID, userID).Scan(
		&member.RoomID, &member.UserID, &member.JoinedAt,
	)

	// If there was a conflict (user already a member), get the existing record
	if err == sql.ErrNoRows {
		existingQuery := `
		SELECT room_id, user_id, joined_at
		FROM room_members
		WHERE room_id = $1 AND user_id = $2`

		err = db.QueryRowContext(ctx, existingQuery, roomID, userID).Scan(
			&member.RoomID, &member.UserID, &member.JoinedAt,
		)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error joining room: %w", err)
	}

	return &member, nil
}

// LeaveRoom removes a user from a room
func LeaveRoom(ctx context.Context, db *sql.DB, roomID, userID int64) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Check if the user is the creator
	var creatorID int64
	creatorQuery := "SELECT creator_id FROM rooms WHERE id = $1"
	err := db.QueryRowContext(ctx, creatorQuery, roomID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("room not found")
		}
		return fmt.Errorf("error checking room creator: %w", err)
	}

	// Creators cannot leave their own rooms (they should delete them instead)
	if creatorID == userID {
		return fmt.Errorf("room creator cannot leave the room")
	}

	// Remove the user from the room
	query := "DELETE FROM room_members WHERE room_id = $1 AND user_id = $2"
	result, err := db.ExecContext(ctx, query, roomID, userID)
	if err != nil {
		return fmt.Errorf("error leaving room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user is not a member of this room")
	}

	return nil
}

// IsRoomMember checks if a user is a member of a room
func IsRoomMember(ctx context.Context, db *sql.DB, roomID, userID int64) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("database connection is nil")
	}

	// Check if the room is public
	var isPrivate bool
	roomQuery := "SELECT is_private FROM rooms WHERE id = $1"
	err := db.QueryRowContext(ctx, roomQuery, roomID).Scan(&isPrivate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("room not found")
		}
		return false, fmt.Errorf("error checking room: %w", err)
	}

	// If the room is public, anyone can access it
	if !isPrivate {
		return true, nil
	}

	// For private rooms, check if the user is a member
	var count int
	query := "SELECT COUNT(*) FROM room_members WHERE room_id = $1 AND user_id = $2"
	err = db.QueryRowContext(ctx, query, roomID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking membership: %w", err)
	}

	return count > 0, nil
}

// Message operations

// SaveMessage saves a new message to the database
func SaveMessage(ctx context.Context, db *sql.DB, content string, senderID, roomID int64) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	// Check if the user is a member of the room
	isMember, err := IsRoomMember(ctx, db, roomID, senderID)
	if err != nil {
		return 0, err
	}

	if !isMember {
		return 0, fmt.Errorf("user is not a member of this room")
	}

	var messageID int64
	query := `
	INSERT INTO messages (content, sender_id, room_id)
	VALUES ($1, $2, $3)
	RETURNING id`

	err = db.QueryRowContext(ctx, query, content, senderID, roomID).Scan(&messageID)
	if err != nil {
		return 0, fmt.Errorf("error saving message: %w", err)
	}

	return messageID, nil
}

// GetRoomMessages retrieves messages for a room
func GetRoomMessages(ctx context.Context, db *sql.DB, roomID, userID int64, limit, offset int64) ([]Message, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Check if the user is a member of the room
	isMember, err := IsRoomMember(ctx, db, roomID, userID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, fmt.Errorf("user is not a member of this room")
	}

	query := `
	SELECT m.id, m.content, m.sender_id, m.room_id, m.timestamp
	FROM messages m
	WHERE m.room_id = $1
	ORDER BY m.timestamp DESC
	LIMIT $2 OFFSET $3`

	rows, err := db.QueryContext(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.SenderID, &msg.RoomID, &msg.Timestamp); err != nil {
			return nil, fmt.Errorf("error scanning message row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	return messages, nil
}

// GetUserNameByID retrieves a username by user ID
func GetUserNameByID(ctx context.Context, db *sql.DB, userID int64) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database connection is nil")
	}

	var username string
	query := "SELECT username FROM users WHERE id = $1"
	err := db.QueryRowContext(ctx, query, userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("error querying username: %w", err)
	}

	return username, nil
}
