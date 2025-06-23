package websocketcase_test

import (
	"context"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/websocketcase"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 正常系
// 異常系：ユーザ取得失敗
// 異常系：クライアントID生成失敗
// 異常系：クライアント作成失敗
// 異常系：WebSocket登録失敗

func TestConnectUserToRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := websocketcase.NewTestWebsocketUseCase(ctrl)

	// テストデータ
	userID := entity.UserID("user123")
	roomID := entity.RoomID("room123")
	clientID := entity.WsClientID("client123")
	testUser := entity.NewUser(entity.UserParams{
		ID:         userID,
		Name:       "John Doe",
		Email:      "test@mail.com",
		PasswdHash: "hashed_password",
		CreatedAt:  time.Now(),
		UpdatedAt:  nil,
	})
	// モックの作成
	mockConn := mock_service.NewMockWebSocketConnection(ctrl)
	t.Run("正常系", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetUserByID(context.Background(), userID).Return(testUser, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientID().Return(clientID, nil)
		mocks.WsClientRepo.EXPECT().CreateClient(context.Background(), gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().Register(context.Background(), mockConn, userID, roomID).Return(nil)

		// テスト実行
		request := websocketcase.ConnectUserToRoomRequest{
			UserID: userID,
			RoomID: roomID,
			Conn:   mockConn,
		}
		err := useCase.ConnectUserToRoom(context.Background(), request)

		// 検証
		assert.NoError(t, err)
	})

	t.Run("異常系：ユーザ取得失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetUserByID(context.Background(), userID).Return(nil, assert.AnError)
		// テスト実行
		request := websocketcase.ConnectUserToRoomRequest{
			UserID: userID,
			RoomID: roomID,
			Conn:   mockConn,
		}
		err := useCase.ConnectUserToRoom(context.Background(), request)

		// 検証
		assert.Error(t, err)
	})

	t.Run("異常系：クライアントID生成失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetUserByID(context.Background(), userID).Return(testUser, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientID().Return(entity.WsClientID(""), assert.AnError)
		// テスト実行
		request := websocketcase.ConnectUserToRoomRequest{
			UserID: userID,
			RoomID: roomID,
			Conn:   mockConn,
		}
		err := useCase.ConnectUserToRoom(context.Background(), request)

		// 検証
		assert.Error(t, err)
	})

	t.Run("異常系：クライアント作成失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetUserByID(context.Background(), userID).Return(testUser, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientID().Return(clientID, nil)
		mocks.WsClientRepo.EXPECT().CreateClient(context.Background(), gomock.Any()).Return(assert.AnError)
		// テスト実行
		request := websocketcase.ConnectUserToRoomRequest{
			UserID: userID,
			RoomID: roomID,
			Conn:   mockConn,
		}
		err := useCase.ConnectUserToRoom(context.Background(), request)

		// 検証
		assert.Error(t, err)
	})

	t.Run("異常系：WebSocket登録失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetUserByID(context.Background(), userID).Return(testUser, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientID().Return(clientID, nil)
		mocks.WsClientRepo.EXPECT().CreateClient(context.Background(), gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().Register(context.Background(), mockConn, userID, roomID).Return(assert.AnError)
		// テスト実行
		request := websocketcase.ConnectUserToRoomRequest{
			UserID: userID,
			RoomID: roomID,
			Conn:   mockConn,
		}
		err := useCase.ConnectUserToRoom(context.Background(), request)

		// 検証
		assert.Error(t, err)
	})
}
