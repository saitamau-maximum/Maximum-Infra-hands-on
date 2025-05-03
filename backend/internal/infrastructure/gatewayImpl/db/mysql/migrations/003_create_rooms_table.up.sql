CREATE TABLE rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    public_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    FULLTEXT INDEX ft_idx_rooms_name (name) WITH PARSER ngram
) ENGINE=InnoDB
CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;