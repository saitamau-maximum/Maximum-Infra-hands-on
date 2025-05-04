package usecase_test

import (
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetMessageHistoryInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockMsgRepo := mock_repository.NewMockMessageRepository(ctrl)
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)

	params := usecase.NewMessageUseCaseParams{
		MsgRepo:  mockMsgRepo,
		RoomRepo: mockRoomRepo,
		UserRepo: mockUserRepo,
	}

	messageUseCase := usecase.NewMessageUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")
		beforeSentAt := time.Now()
		limit := 10
		messages := []*entity.Message{
			entity.NewMessage(entity.MessageParams{
				ID:      "msg1",
				RoomID:  roomID,
				Content: "Hello",
				SentAt:  beforeSentAt.Add(-1 * time.Minute),
			}),
			entity.NewMessage(entity.MessageParams{
				ID:      "msg2",
				RoomID:  roomID,
				Content: "World",
				SentAt:  beforeSentAt.Add(-2 * time.Minute),
			}),
		}
		nextBeforeSentAt := beforeSentAt.Add(-2 * time.Minute)
		hasNext := true

		mockMsgRepo.EXPECT().GetMessageHistoryInRoom(roomID, limit, beforeSentAt).Return(messages, nextBeforeSentAt, hasNext, nil)

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomID:       roomID,
			Limit:        limit,
			BeforeSentAt: beforeSentAt,
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(req)

		assert.NoError(t, err)
		assert.Equal(t, messages, resp.Messages)
		assert.Equal(t, nextBeforeSentAt, resp.NextBeforeSentAt)
		assert.Equal(t, hasNext, resp.HasNext)
	})

	t.Run("メッセージ取得エラー", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")
		beforeSentAt := time.Now()
		limit := 10
		expectedErr := errors.New("failed to fetch messages")

		mockMsgRepo.EXPECT().GetMessageHistoryInRoom(roomID, limit, beforeSentAt).Return(nil, time.Time{}, false, expectedErr)

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomID:       roomID,
			Limit:        limit,
			BeforeSentAt: beforeSentAt,
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(req)

		assert.Error(t, err)
		assert.Empty(t, resp.Messages)
		assert.Equal(t, expectedErr, err)
	})
}
