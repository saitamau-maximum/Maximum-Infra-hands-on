package inmemorymsgcacheimpl_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
	inmemorymsgcacheimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/messageCacheImpl/Inmemory"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TODO: テストの充実
func TestMessageCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockMsgRepo := mock_repository.NewMockMessageRepository(ctrl)
	cache := inmemorymsgcacheimpl.NewMessageCacheService(&inmemorymsgcacheimpl.NewMessageCacheServiceParams{
		MsgRepo: mockMsgRepo,
	})
	roomID := entity.RoomID("test_room")

	// メッセージを順に追加（21件）
	for i := 0; i < 21; i++ {
		msg := entity.NewMessage(entity.MessageParams{
			ID:      entity.MessageID("msg" + strconv.Itoa(i)),
			RoomID:  roomID,
			Content: "message " + string(strconv.Itoa(i)),
			SentAt:  time.Now().Add(time.Duration(i) * time.Second),
		})
		err := cache.AddMessage(context.Background(), roomID, msg)
		assert.NoError(t, err)
	}

	// 取得して件数チェック（20件であること）
	messages, err := cache.GetRecentMessages(context.Background(), roomID)
	assert.NoError(t, err)
	assert.Len(t, messages, service.DefaultRecentMessageLimit())

	// 最初の1件（最古）は捨てられている
	assert.Equal(t, "message 1", messages[0].GetContent())
	assert.Equal(t, "message 20", messages[len(messages)-1].GetContent())
}
