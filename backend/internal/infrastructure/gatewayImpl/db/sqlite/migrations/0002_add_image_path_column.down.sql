-- 1. 一時テーブルを作成（image_pathを含まない）
CREATE TABLE users_temp (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME
);

-- 2. データをコピー（image_pathを除く）
INSERT INTO users_temp (id, name, email, password_hash, created_at, updated_at)
SELECT id, name, email, password_hash, created_at, updated_at FROM users;

-- 3. 元のテーブルを削除
DROP TABLE users;

-- 4. 一時テーブルをリネーム
ALTER TABLE users_temp RENAME TO users;
