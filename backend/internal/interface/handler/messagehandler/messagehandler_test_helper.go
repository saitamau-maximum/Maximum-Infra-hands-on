package messagehandler

import (
	"example.com/infrahandson/internal/infrastructure/validator"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_messagecase "example.com/infrahandson/test/mocks/usecase/messagecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

type mockDeps struct {
	MsgUseCase mock_messagecase.MockMessageUseCaseInterface
	Logger     mock_adapter.MockLoggerAdapter
}

func NewTestMessageHandler(
	ctrl *gomock.Controller,
) (MessageHandlerInterface, mockDeps, *echo.Echo) {
	mockMsgUseCase := mock_messagecase.NewMockMessageUseCaseInterface(ctrl)
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)
	params := NewMessageHandlerParams{
		MsgUseCase: mockMsgUseCase,
		Logger:     mockLogger,
	}
	handler := NewMessageHandler(params)

	mockDeps := mockDeps{
		MsgUseCase: *mockMsgUseCase,
		Logger:     *mockLogger,
	}

	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	return handler, mockDeps, e
}
