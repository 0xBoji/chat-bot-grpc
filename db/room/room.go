package room

import (
	"context"
	"database/sql"
)

// Room represents a chat room in the database
type Room struct {
	ID          int64
	Name        string
	Description string
	CreatorID   int64
}

// Repository handles database operations for rooms
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new room repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateRoom creates a new room in the database
func (r *Repository) CreateRoom(ctx context.Context, name, description string, creatorID int64) (int64, error) {
	var roomID int64
	query := `INSERT INTO rooms (name, description, creator_id) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, name, description, creatorID).Scan(&roomID)
	return roomID, err
}

// GetRoom retrieves a room by ID
func (r *Repository) GetRoom(ctx context.Context, roomID int64) (*Room, error) {
	room := &Room{}
	query := `SELECT id, name, description, creator_id FROM rooms WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&room.ID, &room.Name, &room.Description, &room.CreatorID)
	return room, err
}

// GetUserRooms retrieves all rooms a user is a member of
func (r *Repository) GetUserRooms(ctx context.Context, userID int64) ([]Room, error) {
	query := `
		SELECT r.id, r.name, r.description, r.creator_id
		FROM rooms r
		JOIN room_members rm ON r.id = rm.room_id
		WHERE rm.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.CreatorID); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// RoomExists checks if a room with the given ID exists
func (r *Repository) RoomExists(ctx context.Context, roomID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&exists)
	return exists, err
}

// IsRoomMember checks if a user is a member of a room
func (r *Repository) IsRoomMember(ctx context.Context, roomID, userID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`
	err := r.db.QueryRowContext(ctx, query, roomID, userID).Scan(&exists)
	return exists, err
}

// AddRoomMember adds a user to a room
func (r *Repository) AddRoomMember(ctx context.Context, roomID, userID int64) error {
	query := `INSERT INTO room_members (room_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, roomID, userID)
	return err
}

// RemoveRoomMember removes a user from a room
func (r *Repository) RemoveRoomMember(ctx context.Context, roomID, userID int64) error {
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, roomID, userID)
	return err
}
