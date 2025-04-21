package usecase_test

import (
	"errors"
	"testing"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/usecase"
	mock_repository "example.com/webrtc-practice/mocks/domain/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetMessageHistoryInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockMsgRepo := mock_repository.NewMockMessageRepository(ctrl)
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)

	params := usecase.NewMessageUseCaseParams{
		MsgRepo:  mockMsgRepo,
		RoomRepo: mockRoomRepo,
	}

	messageUseCase := usecase.NewMessageUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomPublicID := entity.RoomPublicID("public_room_1")
		roomID := entity.RoomID(1)
		beforeSentAt := time.Now()
		limit := 10
		messages := []*entity.Message{
			entity.NewMessage(entity.MessageParams{
				ID:      "1",
				RoomID:  roomID,
				Content: "Hello",
				SentAt:  beforeSentAt.Add(-1 * time.Minute),
			}),
			entity.NewMessage(entity.MessageParams{
				ID:      "2",
				RoomID:  roomID,
				Content: "World",
				SentAt:  beforeSentAt.Add(-2 * time.Minute),
			}),
		}
		nextBeforeSentAt := beforeSentAt.Add(-2 * time.Minute)
		hasNext := true

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(roomPublicID).Return(roomID, nil)
		mockMsgRepo.EXPECT().GetMessageHistoryInRoom(roomID, limit, beforeSentAt).Return(messages, nextBeforeSentAt, hasNext, nil)

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomPublicID: roomPublicID,
			Limit:        limit,
			BeforeSentAt: beforeSentAt,
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(req)

		assert.NoError(t, err)
		assert.Equal(t, messages, resp.Messages)
		assert.Equal(t, nextBeforeSentAt, resp.NextBeforeSentAt)
		assert.Equal(t, hasNext, resp.HasNext)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		roomPublicID := entity.RoomPublicID("nonexistent_room")
		mockRoomRepo.EXPECT().GetRoomIDByPublicID(roomPublicID).Return(entity.RoomID(0), errors.New("room not found"))

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomPublicID: roomPublicID,
			Limit:        10,
			BeforeSentAt: time.Now(),
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(req)

		assert.Error(t, err)
		assert.Empty(t, resp.Messages)
		assert.Equal(t, "room not found", err.Error())
	})

	t.Run("メッセージ取得エラー", func(t *testing.T) {
		roomPublicID := entity.RoomPublicID("public_room_1")
		roomID := entity.RoomID(1)
		beforeSentAt := time.Now()
		limit := 10
		expectedErr := errors.New("failed to fetch messages")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(roomPublicID).Return(roomID, nil)
		mockMsgRepo.EXPECT().GetMessageHistoryInRoom(roomID, limit, beforeSentAt).Return(nil, time.Time{}, false, expectedErr)

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomPublicID: roomPublicID,
			Limit:        limit,
			BeforeSentAt: beforeSentAt,
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(req)

		assert.Error(t, err)
		assert.Empty(t, resp.Messages)
		assert.Equal(t, expectedErr, err)
	})
}
