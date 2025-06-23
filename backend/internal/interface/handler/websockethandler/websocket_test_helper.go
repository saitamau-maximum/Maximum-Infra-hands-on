package websockethandler

import (
	"example.com/infrahandson/internal/infrastructure/validator"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	mock_websocketcase "example.com/infrahandson/test/mocks/usecase/websocketcase"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

type mockDeps struct {
	WsUseCase     mock_websocketcase.MockWebsocketUseCaseInterface
	WsUpgrader    mock_adapter.MockWebSocketUpgraderAdapter
	WsConnFactory mock_factory.MockWebSocketConnectionFactory
	UserIDFactory mock_factory.MockUserIDFactory
	RoomIDFactory mock_factory.MockRoomIDFactory
	Logger        mock_adapter.MockLoggerAdapter
}

func NewTestWebsocketHandler(
	ctrl *gomock.Controller,
) (WebSocketHandlerInterface, mockDeps, *echo.Echo) {
	mockWsUseCase := mock_websocketcase.NewMockWebsocketUseCaseInterface(ctrl)
	mockWsUpGrader := mock_adapter.NewMockWebSocketUpgraderAdapter(ctrl)
	mockWsConnFactory := mock_factory.NewMockWebSocketConnectionFactory(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)
	mockDeps := mockDeps{
		WsUseCase:     *mockWsUseCase,
		WsUpgrader:    *mockWsUpGrader,
		WsConnFactory: *mockWsConnFactory,
		UserIDFactory: *mockUserIDFactory,
		RoomIDFactory: *mockRoomIDFactory,
		Logger:        *mockLogger,
	}
	handler := NewWebSocketHandler(NewWebSocketHandlerParams{
		WsUseCase:     mockWsUseCase,
		WsUpgrader:    mockWsUpGrader,
		WsConnFactory: mockWsConnFactory,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:        mockLogger,
	})

	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	return handler, mockDeps, e
}
