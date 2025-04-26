package sqlite3_test

import (
	"testing"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	sqlite3 "example.com/webrtc-practice/internal/infrastructure/repository_impl/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestRoomRepositoryImpl_GetUsersInRoom(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := sqlite3.NewRoomRepositoryImpl(&sqlite3.NewRoomRepositoryImplParams{
		DB: db,
	})

	// データ準備
	// ユーザー登録
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(`INSERT INTO users (public_id, name, email, password_hash, created_at) VALUES
		('user-public-1', 'Alice', 'alice@example.com', 'hash1', ?),
		('user-public-2', 'Bob', 'bob@example.com', 'hash2', ?)`,
		now, now,
	)
	require.NoError(t, err)

	// ルーム登録
	_, err = db.Exec(`INSERT INTO rooms (public_id, name) VALUES ('room-public-1', 'Test Room')`)
	require.NoError(t, err)

	// room_members登録
	_, err = db.Exec(`INSERT INTO room_members (room_id, user_id) VALUES (1, '1'), (1, '2')`)
	require.NoError(t, err)

	// テスト対象呼び出し
	users, err := repo.GetUsersInRoom(entity.RoomID(1))
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
