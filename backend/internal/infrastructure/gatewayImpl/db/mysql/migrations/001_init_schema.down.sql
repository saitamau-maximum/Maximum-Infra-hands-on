-- rooms テーブルの FULLTEXT インデックス削除
DROP INDEX ft_idx_rooms_name ON rooms;

-- room_members テーブル削除（外部キー制約などを考慮）
DROP TABLE IF EXISTS room_members;

-- messages テーブル削除
DROP TABLE IF EXISTS messages;

-- rooms テーブル削除
DROP TABLE IF EXISTS rooms;

-- users テーブル削除
DROP TABLE IF EXISTS users;
