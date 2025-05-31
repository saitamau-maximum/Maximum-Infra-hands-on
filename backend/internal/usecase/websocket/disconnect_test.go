package websocket_test

import (
	"context"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	wsUC "example.com/infrahandson/internal/usecase/websocket"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2.GetConnectionByUserID失敗
// 3.GetClientsByUserID失敗
// 4.Unregister失敗
// 5. DeleteClient失敗

func TestDisconnectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, mocks := wsUC.NewTestWebsocketUseCase(ctrl)

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

		request := wsUC.DisconnectUserRequest{UserID: userID}
		err := useCase.DisconnectUser(context.Background(), request)

		assert.NoError(t, err)
	})

	t.Run("GetConnectionByUserID失敗", func(t *testing.T) {
		userID := entity.UserID("user123")
		expectedErr := assert.AnError

		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(context.Background(), userID).Return(nil, expectedErr)

		request := wsUC.DisconnectUserRequest{UserID: userID}
		err := useCase.DisconnectUser(context.Background(), request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("GetClientsByUserID失敗", func(t *testing.T) {
		userID := entity.UserID("user123")
		expectedErr := assert.AnError

		mockConn := mock_service.NewMockWebSocketConnection(ctrl)
		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(context.Background(), userID).Return(mockConn, nil)
		mocks.WsClientRepo.EXPECT().GetClientsByUserID(context.Background(), userID).Return(nil, expectedErr)

		request := wsUC.DisconnectUserRequest{UserID: userID}
		err := useCase.DisconnectUser(context.Background(), request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Unregister失敗", func(t *testing.T) {
		userID := entity.UserID("user123")
		expectedErr := assert.AnError

		mockConn := mock_service.NewMockWebSocketConnection(ctrl)
		mockClient := entity.NewWebsocketClient(entity.WebsocketClientParams{
			ID:     entity.WsClientID("client123"),
			UserID: userID,
			RoomID: entity.RoomID("room123"),
		})

		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(context.Background(), userID).Return(mockConn, nil)
		mocks.WsClientRepo.EXPECT().GetClientsByUserID(context.Background(), userID).Return(mockClient, nil)
		mocks.WebsocketManager.EXPECT().Unregister(context.Background(), mockConn).Return(expectedErr)

		request := wsUC.DisconnectUserRequest{UserID: userID}
		err := useCase.DisconnectUser(context.Background(), request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
	
	t.Run("DeleteClient失敗", func(t *testing.T) {
		userID := entity.UserID("user123")
		expectedErr := assert.AnError

		mockConn := mock_service.NewMockWebSocketConnection(ctrl)
		mockClient := entity.NewWebsocketClient(entity.WebsocketClientParams{
			ID:     entity.WsClientID("client123"),
			UserID: userID,
			RoomID: entity.RoomID("room123"),
		})

		mocks.WebsocketManager.EXPECT().GetConnectionByUserID(context.Background(), userID).Return(mockConn, nil)
		mocks.WsClientRepo.EXPECT().GetClientsByUserID(context.Background(), userID).Return(mockClient, nil)
		mocks.WebsocketManager.EXPECT().Unregister(context.Background(), mockConn).Return(nil)
		mocks.WsClientRepo.EXPECT().DeleteClient(context.Background(), mockClient.GetID()).Return(expectedErr)

		request := wsUC.DisconnectUserRequest{UserID: userID}
		err := useCase.DisconnectUser(context.Background(), request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
