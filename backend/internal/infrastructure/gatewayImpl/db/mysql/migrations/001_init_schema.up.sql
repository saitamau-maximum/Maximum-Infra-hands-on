CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    public_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME
);

CREATE INDEX idx_users_name ON users(name);

CREATE TABLE rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    public_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    FULLTEXT INDEX ft_idx_rooms_name (name) WITH PARSER ngram
) ENGINE=InnoDB
CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- パーサーがある前提での FULLTEXT インデックス作成
-- ただし、ngram パーサーは MySQL 8.0 以降でのみ使用可能
CREATE INDEX idx_rooms_name ON rooms(name);

CREATE TABLE IF NOT EXISTS room_members (
    id INT AUTO_INCREMENT PRIMARY KEY,
    room_id INT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE INDEX idx_room_members_user_id ON room_members(user_id);
CREATE INDEX idx_room_members_room_id ON room_members(room_id);

CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(255) PRIMARY KEY,
    public_id VARCHAR(255) NOT NULL UNIQUE,
    room_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    sent_at DATETIME NOT NULL
);

CREATE INDEX idx_messages_room_id_sent_at ON messages(room_id, sent_at DESC);
