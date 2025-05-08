package sqliteroomrepoimpl_test

import (
	"context"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	sqliteroomrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/sqlite"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// SQLiteインメモリDBをセットアップ
	initializer := sqlitegatewayimpl.NewSQLiteInitializer(&sqlitegatewayimpl.NewSQLiteInitializerParams{
		Path:           ":memory:",
		MigrationsPath: "file://migrations",
	})

	db, err := initializer.Init()
	require.NoError(t, err)

	// テスト用スキーマを作成
	_, err = db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME
	);
	CREATE TABLE rooms (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL
	);
	CREATE TABLE room_members (
		id TEXT PRIMARY KEY,
		room_id INTEGER NOT NULL,
		user_id TEXT NOT NULL,
		FOREIGN KEY (room_id) REFERENCES rooms(id)
	);`)
	require.NoError(t, err)

	return db
}

func TestRoomRepositoryImpl_GetUsersInRoom(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := sqliteroomrepoimpl.NewRoomRepositoryImpl(&sqliteroomrepoimpl.NewRoomRepositoryImplParams{
		DB: db,
	})

	// データ準備
	// ユーザー登録
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(`INSERT INTO users (id, name, email, password_hash, created_at) VALUES
		('user-public-1', 'Alice', 'alice@example.com', 'hash1', ?),
		('user-public-2', 'Bob', 'bob@example.com', 'hash2', ?)`,
		now, now,
	)
	require.NoError(t, err)

	// ルーム登録
	_, err = db.Exec(`INSERT INTO rooms (id, name) VALUES ('room-public-1', 'Test Room')`)
	require.NoError(t, err)

	// room_members登録
	_, err = db.Exec(`INSERT INTO room_members (room_id, user_id) VALUES (1, '1'), (1, '2')`)
	require.NoError(t, err)

	// テスト対象呼び出し
	users, err := repo.GetUsersInRoom(context.Background(), entity.RoomID("room-public-1"))
	require.NoError(t, err)
	require.Len(t, users, 2)

	// 名前で確認
	userNames := map[string]bool{
		users[0].GetName(): true,
		users[1].GetName(): true,
	}
	require.True(t, userNames["Alice"])
	require.True(t, userNames["Bob"])
}
