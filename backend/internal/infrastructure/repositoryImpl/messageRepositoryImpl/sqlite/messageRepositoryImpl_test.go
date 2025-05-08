package sqlitemsgrepoimpl_test

import (
	"context"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	sqlitemsgrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLiteドライバ
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	schema := `
CREATE TABLE messages (
	id TEXT PRIMARY KEY,
	public_id TEXT NOT NULL,
	room_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	content TEXT NOT NULL,
	sent_at DATETIME NOT NULL
);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestMessageRepositoryImpl_CreateAndGetMessageHistoryInRoom(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlitemsgrepoimpl.NewMessageRepositoryImpl(&sqlitemsgrepoimpl.NewMessageRepositoryImplParams{DB: db})

	// メッセージを作成
	now := time.Now().UTC()
	message := entity.NewMessage(entity.MessageParams{
		ID: entity.MessageID("test-public-id"),
		RoomID:   entity.RoomID("test-room-id"),
		UserID:   entity.UserID("test-user-id"),
		Content:  "Hello, World!",
		SentAt:   now,
	})

	err := repo.CreateMessage(context.Background(), message)
	assert.NoError(t, err)

	// メッセージ履歴を取得
	messages, nextBeforeSentAt, hasNext, err := repo.GetMessageHistoryInRoom(
		context.Background(),
		entity.RoomID("test-room-id"),
		10,
		time.Now().Add(1*time.Hour), // 未来時間を指定しているので、登録したメッセージが対象になる
	)

	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, message.GetContent(), messages[0].GetContent())
	assert.Equal(t, message.GetRoomID(), messages[0].GetRoomID())
	assert.WithinDuration(t, message.GetSentAt(), messages[0].GetSentAt(), time.Second)
	assert.Equal(t, now, nextBeforeSentAt)
	assert.False(t, hasNext)
}
