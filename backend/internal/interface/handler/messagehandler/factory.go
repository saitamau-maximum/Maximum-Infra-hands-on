package messagehandler

import (
	"errors"

	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/usecase/messagecase"
)

type NewMessageHandlerParams struct {
	MsgUseCase messagecase.MessageUseCaseInterface
	Logger     adapter.LoggerAdapter
}

func (p *NewMessageHandlerParams) Validate() error {
	if p.MsgUseCase == nil {
		return errors.New("messageUseCase is required")
	}
	if p.Logger == nil {
		return errors.New("logger is required")
	}
	return nil
}

func NewMessageHandler(params NewMessageHandlerParams) MessageHandlerInterface {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &MessageHandler{
		MsgUseCase: params.MsgUseCase,
		Logger:     params.Logger,
	}
}