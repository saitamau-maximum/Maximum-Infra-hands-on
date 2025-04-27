package usecase_test

import (
	"testing"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/usecase"
	mock_repository "example.com/webrtc-practice/test/mocks/domain/repository"
	mock_service "example.com/webrtc-practice/test/mocks/domain/service"
	mock_factory "example.com/webrtc-practice/test/mocks/interface/factory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockDeps struct {
	UserRepo         *mock_repository.MockUserRepository
	RoomRepo         *mock_repository.MockRoomRepository
	MsgRepo          *mock_repository.MockMessageRepository
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
	mockWsClientRepo := mock_repository.NewMockWebsocketClientRepository(ctrl)
	mockWebsocketManager := mock_service.NewMockWebsocketManager(ctrl)
	mockClientIDFactory := mock_factory.NewMockWsClientIDFactory(ctrl)
	mockMsgIDFactory := mock_factory.NewMockMessageIDFactory(ctrl)

	params := usecase.NewWebsocketUseCaseParams{
		UserRepo:         mockUserRepo,
		RoomRepo:         mockRoomRepo,
		MsgRepo:          mockMsgRepo,
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
	userID := entity.UserID(1)
	userPublicID := entity.UserPublicID("user123")
	publicRoomID := entity.RoomPublicID("room123")
	roomID := entity.RoomID(123)
	clientID := entity.WsClientPublicID("client123")
	testUser := entity.NewUser(entity.UserParams{
		ID:         userID,
		PublicID:   userPublicID,
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
		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.UserRepo.EXPECT().GetUserByID(userID).Return(testUser, nil)
		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(publicRoomID).Return(roomID, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientPublicID().Return(clientID, nil)
		mocks.WsClientRepo.EXPECT().CreateClient(gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().Register(mockConn, userID, roomID).Return(nil)

		// テスト実行
		request := usecase.ConnectUserToRoomRequest{
			UserPublicID: userPublicID,
			RoomPublicID: publicRoomID,
			Conn:         mockConn,
		}
		err := useCase.ConnectUserToRoom(request)

		// 検証
		assert.NoError(t, err)
	})

	t.Run("異常系：ユーザ取得失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.UserRepo.EXPECT().GetUserByID(userID).Return(nil, assert.AnError)
		// テスト実行
		request := usecase.ConnectUserToRoomRequest{
			UserPublicID: userPublicID,
			RoomPublicID: publicRoomID,
			Conn:         mockConn,
		}
		err := useCase.ConnectUserToRoom(request)

		// 検証
		assert.Error(t, err)
	})

	t.Run("異常系：部屋取得失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.UserRepo.EXPECT().GetUserByID(userID).Return(testUser, nil)
		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(publicRoomID).Return(entity.RoomID(0), assert.AnError)
		// テスト実行
		request := usecase.ConnectUserToRoomRequest{
			UserPublicID: userPublicID,
			RoomPublicID: publicRoomID,
			Conn:         mockConn,
		}
		err := useCase.ConnectUserToRoom(request)

		// 検証
		assert.Error(t, err)
	})

	t.Run("異常系：クライアントID生成失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.UserRepo.EXPECT().GetUserByID(userID).Return(testUser, nil)
		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(publicRoomID).Return(roomID, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientPublicID().Return(entity.WsClientPublicID(""), assert.AnError)
		// テスト実行
		request := usecase.ConnectUserToRoomRequest{
			UserPublicID: userPublicID,
			RoomPublicID: publicRoomID,
			Conn:         mockConn,
		}
		err := useCase.ConnectUserToRoom(request)

		// 検証
		assert.Error(t, err)
	})

	t.Run("異常系：クライアント作成失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.UserRepo.EXPECT().GetUserByID(userID).Return(testUser, nil)
		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(publicRoomID).Return(roomID, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientPublicID().Return(clientID, nil)
		mocks.WsClientRepo.EXPECT().CreateClient(gomock.Any()).Return(assert.AnError)
		// テスト実行
		request := usecase.ConnectUserToRoomRequest{
			UserPublicID: userPublicID,
			RoomPublicID: publicRoomID,
			Conn:         mockConn,
		}
		err := useCase.ConnectUserToRoom(request)

		// 検証
		assert.Error(t, err)
	})

	t.Run("異常系：WebSocket登録失敗", func(t *testing.T) {
		// モックの期待値設定
		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.UserRepo.EXPECT().GetUserByID(userID).Return(testUser, nil)
		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(publicRoomID).Return(roomID, nil)
		mocks.ClientIDFactory.EXPECT().NewWsClientPublicID().Return(clientID, nil)
		mocks.WsClientRepo.EXPECT().CreateClient(gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().Register(mockConn, userID, roomID).Return(assert.AnError)
		// テスト実行
		request := usecase.ConnectUserToRoomRequest{
			UserPublicID: userPublicID,
			RoomPublicID: publicRoomID,
			Conn:         mockConn,
		}
		err := useCase.ConnectUserToRoom(request)

		// 検証
		assert.Error(t, err)
	})
}

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := newTestWebsocketUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		roomPublicID := entity.RoomPublicID("room123")
		roomID := entity.RoomID(123)
		senderID := entity.UserID(1)
		senderPublicID := entity.UserPublicID("user123")
		content := "Hello, World!"
		messageID := entity.MessagePublicID("msg123")

		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(roomPublicID).Return(roomID, nil)
		mocks.MsgIDFactory.EXPECT().NewMessagePublicID().Return(messageID, nil)
		mocks.UserRepo.EXPECT().GetIDByPublicID(senderPublicID).Return(senderID, nil)
		mocks.MsgRepo.EXPECT().CreateMessage(gomock.Any()).Return(nil)
		mocks.WebsocketManager.EXPECT().BroadcastToRoom(roomID, gomock.Any()).Return(nil)

		request := usecase.SendMessageRequest{
			RoomPublicID: roomPublicID,
			Sender:       senderPublicID,
			Content:      content,
		}
		err := useCase.SendMessage(request)

		assert.NoError(t, err)
	})

	t.Run("異常系：部屋取得失敗", func(t *testing.T) {
		roomPublicID := entity.RoomPublicID("room123")
		senderPublicID := entity.UserPublicID("user123")
		content := "Hello, World!"

		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(roomPublicID).Return(entity.RoomID(0), assert.AnError)

		request := usecase.SendMessageRequest{
			RoomPublicID: roomPublicID,
			Sender:       senderPublicID,
			Content:      content,
		}
		err := useCase.SendMessage(request)

		assert.Error(t, err)
	})
}

func TestDisconnectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := newTestWebsocketUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		userID := entity.UserID(1)
		userPublicID := entity.UserPublicID("user123")
		mockConn := mock_service.NewMockWebSocketConnection(ctrl)
		mockClient := entity.NewWebsocketClient(entity.WebsocketClientParams{
			PublicID: entity.WsClientPublicID("client123"),
			UserID:   userID,
			RoomID:   entity.RoomID(123),
		})

		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(userID).Return(mockConn, nil)
		mocks.WsClientRepo.EXPECT().GetClientsByUserID(userID).Return(mockClient, nil)
		mocks.WebsocketManager.EXPECT().Unregister(mockConn).Return(nil)
		mocks.WsClientRepo.EXPECT().DeleteClient(mockClient.GetID()).Return(nil)

		request := usecase.DisconnectUserRequest{UserID: userPublicID}
		err := useCase.DisconnectUser(request)

		assert.NoError(t, err)
	})

	t.Run("異常系：接続取得失敗", func(t *testing.T) {
		userPublicID := entity.UserPublicID("user123")
		userID := entity.UserID(1)

		mocks.UserRepo.EXPECT().GetIDByPublicID(userPublicID).Return(userID, nil)
		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(userID).Return(nil, assert.AnError)

		request := usecase.DisconnectUserRequest{UserID: userPublicID}
		err := useCase.DisconnectUser(request)

		assert.Error(t, err)
	})
}

func TestGetMessageHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := newTestWebsocketUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		publicRoomID := entity.RoomPublicID("room123")
		roomID := entity.RoomID(123)
		messages := []*entity.Message{
			entity.NewMessage(entity.MessageParams{
				ID:       entity.MessageID(1),
				PublicID: entity.MessagePublicID("msg1"),
				RoomID:   roomID,
				UserID:   entity.UserID(1),
				Content:  "Hello",
				SentAt:   time.Now(),
			}),
		}

		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(publicRoomID).Return(roomID, nil)
		mocks.MsgRepo.EXPECT().GetMessagesByRoomID(roomID).Return(messages, nil)

		request := usecase.GetMessageHistoryRequest{PublicRoomID: publicRoomID}
		response, err := useCase.GetMessageHistory(request)

		assert.NoError(t, err)
		assert.Equal(t, messages, response.Messages)
	})

	t.Run("異常系：部屋取得失敗", func(t *testing.T) {
		publicRoomID := entity.RoomPublicID("room123")

		mocks.RoomRepo.EXPECT().GetRoomIDByPublicID(publicRoomID).Return(entity.RoomID(0), assert.AnError)

		request := usecase.GetMessageHistoryRequest{PublicRoomID: publicRoomID}
		_, err := useCase.GetMessageHistory(request)

		assert.Error(t, err)
	})
}
