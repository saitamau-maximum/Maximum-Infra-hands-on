package websocketcase_test

import (
	"context"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/websocketcase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := websocketcase.NewTestWebsocketUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("room123")
		senderID := entity.UserID("user123")
		content := "Hello, World!"
		messageID := entity.MessageID("msg123")

		mocks.MsgIDFactory.EXPECT().NewMessageID().Return(messageID, nil)
		mocks.MsgRepo.EXPECT().CreateMessage(context.Background(), gomock.Any()).Return(nil)
		mocks.MsgCache.EXPECT().AddMessage(context.Background(), roomID, gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().BroadcastToRoom(context.Background(), roomID, gomock.Any()).Return(nil)

		request := websocketcase.SendMessageRequest{
			RoomID:  roomID,
			Sender:  senderID,
			Content: content,
		}
		err := useCase.SendMessage(context.Background(), request)

		assert.NoError(t, err)
	})

	t.Run("異常系：メッセージID生成失敗", func(t *testing.T) {
		roomID := entity.RoomID("room123")
		senderID := entity.UserID("user123")
		content := "Hello, World!"

		mocks.MsgIDFactory.EXPECT().NewMessageID().Return(entity.MessageID(""), assert.AnError)

		request := websocketcase.SendMessageRequest{
			RoomID:  roomID,
			Sender:  senderID,
			Content: content,
		}
		err := useCase.SendMessage(context.Background(), request)

		assert.Error(t, err)
	})

	t.Run("異常系：メッセージ作成失敗", func(t *testing.T) {
		roomID := entity.RoomID("room123")
		senderID := entity.UserID("user123")
		content := "Hello, World!"
		messageID := entity.MessageID("msg123")

		mocks.MsgIDFactory.EXPECT().NewMessageID().Return(messageID, nil)
		mocks.MsgRepo.EXPECT().CreateMessage(context.Background(), gomock.Any()).Return(assert.AnError)

		request := websocketcase.SendMessageRequest{
			RoomID:  roomID,
			Sender:  senderID,
			Content: content,
		}
		err := useCase.SendMessage(context.Background(), request)

		assert.Error(t, err)
	})

	t.Run("異常系：メッセージキャッシュ失敗", func(t *testing.T) {
		roomID := entity.RoomID("room123")
		senderID := entity.UserID("user123")
		content := "Hello, World!"
		messageID := entity.MessageID("msg123")

		mocks.MsgIDFactory.EXPECT().NewMessageID().Return(messageID, nil)
		mocks.MsgRepo.EXPECT().CreateMessage(context.Background(), gomock.Any()).Return(nil)
		mocks.MsgCache.EXPECT().AddMessage(context.Background(), roomID, gomock.Any()).Return(assert.AnError)

		request := websocketcase.SendMessageRequest{
			RoomID:  roomID,
			Sender:  senderID,
			Content: content,
		}
		err := useCase.SendMessage(context.Background(), request)

		assert.Error(t, err)
	})

	t.Run("異常系：メッセージ送信失敗", func(t *testing.T) {
		roomID := entity.RoomID("room123")
		senderID := entity.UserID("user123")
		content := "Hello, World!"
		messageID := entity.MessageID("msg123")

		mocks.MsgIDFactory.EXPECT().NewMessageID().Return(messageID, nil)
		mocks.MsgRepo.EXPECT().CreateMessage(context.Background(), gomock.Any()).Return(nil)
		mocks.MsgCache.EXPECT().AddMessage(context.Background(), roomID, gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().BroadcastToRoom(context.Background(), roomID, gomock.Any()).Return(assert.AnError)

		request := websocketcase.SendMessageRequest{
			RoomID:  roomID,
			Sender:  senderID,
			Content: content,
		}
		err := useCase.SendMessage(context.Background(), request)

		assert.Error(t, err)
	})
}
