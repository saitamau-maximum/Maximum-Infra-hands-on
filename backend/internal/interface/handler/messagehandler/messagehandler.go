package messagehandler

import (
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/usecase/messagecase"
)

type MessageHandler struct {
	MsgUseCase messagecase.MessageUseCaseInterface
	Logger     adapter.LoggerAdapter
}
