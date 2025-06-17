package roomhandler

import (
	"example.com/infrahandson/internal/infrastructure/validator"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	mock_roomcase "example.com/infrahandson/test/mocks/usecase/roomcase"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

type mockDeps struct {
	RoomUseCase   mock_roomcase.MockRoomUseCaseInterface
	UserIDFactory mock_factory.MockUserIDFactory
	RoomIDFactory mock_factory.MockRoomIDFactory
	Logger        mock_adapter.MockLoggerAdapter
}

func NewTestRoomHandler(
	ctrl *gomock.Controller,
) (RoomHandlerInterface, mockDeps, *echo.Echo) {
	mockRoomUseCase := mock_roomcase.NewMockRoomUseCaseInterface(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	params := NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:        mockLogger,
	}
	handler := NewRoomHandler(params)

	mockDeps := mockDeps{
		RoomUseCase:   *mockRoomUseCase,
		UserIDFactory: *mockUserIDFactory,
		RoomIDFactory: *mockRoomIDFactory,
		Logger:        *mockLogger,
	}

	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	return handler, mockDeps, e
}
