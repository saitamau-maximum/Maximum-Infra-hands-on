package usecase_test

import (
	"context"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockDeps struct {
	UserRepo         *mock_repository.MockUserRepository
	RoomRepo         *mock_repository.MockRoomRepository
	MsgRepo          *mock_repository.MockMessageRepository
	MsgCache         *mock_service.MockMessageCacheService
	WsClientRepo     *mock_repository.MockWebsocketClientRepository
	WebsocketManager *mock_service.MockWebsocketManager
	ClientIDFactory  *mock_factory.MockWsClientIDFactory
	MsgIDFactory     *mock_factory.MockMessageIDFactory
}

func newTestWebsocketUseCase(
	ctrl *gomock.Controller,
) (usecase.WebsocketUseCaseInterface, mockDeps) {
	// モックの作成
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockMsgRepo := mock_repository.NewMockMessageRepository(ctrl)
	mockMsgCache := mock_service.NewMockMessageCacheService(ctrl)
	mockWsClientRepo := mock_repository.NewMockWebsocketClientRepository(ctrl)
	mockWebsocketManager := mock_service.NewMockWebsocketManager(ctrl)
	mockClientIDFactory := mock_factory.NewMockWsClientIDFactory(ctrl)
	mockMsgIDFactory := mock_factory.NewMockMessageIDFactory(ctrl)

	params := usecase.NewWebsocketUseCaseParams{
		UserRepo:         mockUserRepo,
		RoomRepo:         mockRoomRepo,
		MsgRepo:          mockMsgRepo,
		MsgCache:         mockMsgCache,
		WsClientRepo:     mockWsClientRepo,
		WebsocketManager: mockWebsocketManager,
		MsgIDFactory:     mockMsgIDFactory,
		ClientIDFactory:  mockClientIDFactory,
	}
	useCase := usecase.NewWebsocketUseCase(params)

	// モックをテスト内で使いたいため構造体で返す
	return useCase, mockDeps{
		UserRepo:         mockUserRepo,
		RoomRepo:         mockRoomRepo,
		MsgRepo:          mockMsgRepo,
		MsgCache:         mockMsgCache,
		WsClientRepo:     mockWsClientRepo,
		WebsocketManager: mockWebsocketManager,
		ClientIDFactory:  mockClientIDFactory,
		MsgIDFactory:     mockMsgIDFactory,
	}
}

func TestConnectUserToRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := newTestWebsocketUseCase(ctrl)

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
		request := usecase.ConnectUserToRoomRequest{
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
		request := usecase.ConnectUserToRoomRequest{
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
		request := usecase.ConnectUserToRoomRequest{
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
		request := usecase.ConnectUserToRoomRequest{
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
		request := usecase.ConnectUserToRoomRequest{
			UserID: userID,
			RoomID: roomID,
			Conn:   mockConn,
		}
		err := useCase.ConnectUserToRoom(context.Background(), request)

		// 検証
		assert.Error(t, err)
	})
}

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := newTestWebsocketUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("room123")
		senderID := entity.UserID("user123")
		content := "Hello, World!"
		messageID := entity.MessageID("msg123")

		mocks.MsgIDFactory.EXPECT().NewMessageID().Return(messageID, nil)
		mocks.MsgRepo.EXPECT().CreateMessage(context.Background(), gomock.Any()).Return(nil)
		mocks.MsgCache.EXPECT().AddMessage(context.Background(), roomID ,gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().BroadcastToRoom(context.Background(), roomID, gomock.Any()).Return(nil)

		request := usecase.SendMessageRequest{
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

		request := usecase.SendMessageRequest{
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

		request := usecase.SendMessageRequest{
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

		request := usecase.SendMessageRequest{
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

		request := usecase.SendMessageRequest{
			RoomID:  roomID,
			Sender:  senderID,
			Content: content,
		}
		err := useCase.SendMessage(context.Background(), request)

		assert.Error(t, err)
	})
}

func TestDisconnectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := newTestWebsocketUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		userID := entity.UserID("user123")
		mockConn := mock_service.NewMockWebSocketConnection(ctrl)
		mockClient := entity.NewWebsocketClient(entity.WebsocketClientParams{
			ID:     entity.WsClientID("client123"),
			UserID: userID,
			RoomID: entity.RoomID("room123"),
		})

		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(context.Background(), userID).Return(mockConn, nil)
		mocks.WsClientRepo.EXPECT().GetClientsByUserID(context.Background(), userID).Return(mockClient, nil)
		mocks.WebsocketManager.EXPECT().Unregister(context.Background(), mockConn).Return(nil)
		mocks.WsClientRepo.EXPECT().DeleteClient(context.Background(), mockClient.GetID()).Return(nil)

		request := usecase.DisconnectUserRequest{UserID: userID}
		err := useCase.DisconnectUser(context.Background(), request)

		assert.NoError(t, err)
	})

	t.Run("異常系：接続取得失敗", func(t *testing.T) {
		userID := entity.UserID("user123")

		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(context.Background(), userID).Return(nil, assert.AnError)

		request := usecase.DisconnectUserRequest{UserID: userID}
		err := useCase.DisconnectUser(context.Background(), request)

		assert.Error(t, err)
	})
}
