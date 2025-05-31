package websocket

import (
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
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

func NewTestWebsocketUseCase(
	ctrl *gomock.Controller,
) (WebsocketUseCaseInterface, mockDeps) {
	// モックの作成
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockMsgRepo := mock_repository.NewMockMessageRepository(ctrl)
	mockMsgCache := mock_service.NewMockMessageCacheService(ctrl)
	mockWsClientRepo := mock_repository.NewMockWebsocketClientRepository(ctrl)
	mockWebsocketManager := mock_service.NewMockWebsocketManager(ctrl)
	mockClientIDFactory := mock_factory.NewMockWsClientIDFactory(ctrl)
	mockMsgIDFactory := mock_factory.NewMockMessageIDFactory(ctrl)

	params := NewWebsocketUseCaseParams{
		UserRepo:         mockUserRepo,
		RoomRepo:         mockRoomRepo,
		MsgRepo:          mockMsgRepo,
		MsgCache:         mockMsgCache,
		WsClientRepo:     mockWsClientRepo,
		WebsocketManager: mockWebsocketManager,
		MsgIDFactory:     mockMsgIDFactory,
		ClientIDFactory:  mockClientIDFactory,
	}
	useCase := NewWebsocketUseCase(params)

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
