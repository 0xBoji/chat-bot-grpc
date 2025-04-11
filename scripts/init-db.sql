-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    creator_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create room_members table
CREATE TABLE IF NOT EXISTS room_members (
    room_id INTEGER REFERENCES rooms(id),
    user_id INTEGER REFERENCES users(id),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id)
);

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    sender_id INTEGER REFERENCES users(id),
    room_id INTEGER REFERENCES rooms(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_room_members_room_id ON room_members(room_id);
CREATE INDEX IF NOT EXISTS idx_room_members_user_id ON room_members(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);

-- Insert test users
-- Password is 'password' hashed with bcrypt
INSERT INTO users (id, username, password_hash)
VALUES
(1, 'testuser', '$2a$10$1/dXSLDvXBGXJdZD8W2Xz.4vL3SLj0Lc5y7jzTmZzU.JIYB8.ZmMW'),
(2, 'admin', '$2a$10$1/dXSLDvXBGXJdZD8W2Xz.4vL3SLj0Lc5y7jzTmZzU.JIYB8.ZmMW'),
(3, 'user1', '$2a$10$1/dXSLDvXBGXJdZD8W2Xz.4vL3SLj0Lc5y7jzTmZzU.JIYB8.ZmMW'),
(4, 'user2', '$2a$10$1/dXSLDvXBGXJdZD8W2Xz.4vL3SLj0Lc5y7jzTmZzU.JIYB8.ZmMW'),
(5, 'user3', '$2a$10$1/dXSLDvXBGXJdZD8W2Xz.4vL3SLj0Lc5y7jzTmZzU.JIYB8.ZmMW')
ON CONFLICT (id) DO NOTHING;

-- Reset the sequence
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
