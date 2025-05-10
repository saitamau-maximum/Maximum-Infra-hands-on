package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/usecase"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

/*
* パターン４つ
* 1. キャッシュから取得される正常系
* 2. DBから取得される正常系
* 3. DBから取得されるが、エラーになる
* 4. キャッシュから取得されるが、エラーになる
 */

func TestGetMessageHistoryInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 共通データ
	const defaultLimit = 10
	roomID := entity.RoomID("public_room_1")
	beforeSentAt := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)

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

	// モックのセットアップ
	mockMsgRepo := mock_repository.NewMockMessageRepository(ctrl)
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockMsgCache := mock_service.NewMockMessageCacheService(ctrl)
	
	params := usecase.NewMessageUseCaseParams{
		MsgRepo:  mockMsgRepo,
		MsgCache: mockMsgCache,
		RoomRepo: mockRoomRepo,
		UserRepo: mockUserRepo,
	}
	messageUseCase := usecase.NewMessageUseCase(params)
	
	t.Run("1. キャッシュ内で完結できる場合", func(t *testing.T) {
		cachedMessages := []*entity.Message{
			entity.NewMessage(entity.MessageParams{
				ID:      "msg1",
				RoomID:  roomID,
				Content: "Latest",
				SentAt:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			}),
			entity.NewMessage(entity.MessageParams{
				ID:      "msg2",
				RoomID:  roomID,
				Content: "Older",
				SentAt:  time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),
			}),
		}

		mockMsgCache.EXPECT().
			GetRecentMessages(context.Background(), roomID).
			Return(cachedMessages, nil)

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomID:       roomID,
			Limit:        service.DefaultRecentMessageLimit(),
			BeforeSentAt: time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC), // キャッシュより新しい
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, cachedMessages, resp.Messages)
		assert.False(t, resp.HasNext)
	})

	t.Run("2. DBからの取得が行われる正常系", func(t *testing.T) {
		// cacheの配列
		cachedMessages := []*entity.Message{
			entity.NewMessage(entity.MessageParams{
				ID:      "msg1",
				RoomID:  roomID,
				Content: "Hello",
				SentAt: beforeSentAt.Add(1 * time.Minute),// リクエストよりも新しい
			}),
		}
		
		mockMsgCache.EXPECT().
			GetRecentMessages(context.Background(), roomID).
			Return(cachedMessages, nil)
		mockMsgRepo.EXPECT().
			GetMessageHistoryInRoom(context.Background(), roomID, defaultLimit, beforeSentAt).
			Return(messages, nextBeforeSentAt, hasNext, nil)

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomID:       roomID,
			Limit:        defaultLimit,
			BeforeSentAt: beforeSentAt,
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, messages, resp.Messages)
		assert.Equal(t, nextBeforeSentAt, resp.NextBeforeSentAt)
		assert.Equal(t, hasNext, resp.HasNext)
	})

	t.Run("3. DBから取得されるエラー", func(t *testing.T) {
		expectedErr := errors.New("failed to fetch messages")
		// cacheの配列
		cachedMessages := []*entity.Message{
			entity.NewMessage(entity.MessageParams{
				ID:      "msg1",
				RoomID:  roomID,
				Content: "Hello",
				SentAt: beforeSentAt.Add(1 * time.Minute),// リクエストよりも新しい
			}),
		}
		mockMsgCache.EXPECT().
			GetRecentMessages(context.Background(), roomID).
			Return(cachedMessages, nil)

		mockMsgRepo.EXPECT().
			GetMessageHistoryInRoom(context.Background(), roomID, defaultLimit, beforeSentAt).
			Return(nil, time.Time{}, false, expectedErr)

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomID:       roomID,
			Limit:        defaultLimit,
			BeforeSentAt: beforeSentAt,
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(context.Background(), req)

		assert.Error(t, err)
		assert.Empty(t, resp.Messages)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("4. キャッシュエラー", func(t *testing.T) {
		mockMsgCache.EXPECT().
			GetRecentMessages(context.Background(), roomID).
			Return(nil, errors.New("cache error"))

		req := usecase.GetMessageHistoryInRoomRequest{
			RoomID:       roomID,
			Limit:        service.DefaultRecentMessageLimit(),
			BeforeSentAt: beforeSentAt, // キャッシュより古い
		}

		resp, err := messageUseCase.GetMessageHistoryInRoom(context.Background(), req)

		assert.Error(t, err)
		assert.Empty(t, resp.Messages)
	})
}
