package roomhandler

import (
	"errors"

	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase/roomcase"
)

type NewRoomHandlerParams struct {
	RoomUseCase   roomcase.RoomUseCaseInterface
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
}

func (p *NewRoomHandlerParams) Validate() error {
	if p.RoomUseCase == nil {
		return errors.New("roomUseCase is required")
	}
	if p.UserIDFactory == nil {
		return errors.New("userIDFactory is required")
	}
	if p.RoomIDFactory == nil {
		return errors.New("roomPubIDFactory is required")
	}
	if p.Logger == nil {
		return errors.New("logger is required")
	}
	return nil
}

func NewRoomHandler(params NewRoomHandlerParams) RoomHandlerInterface {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &RoomHandler{
		RoomUseCase:   params.RoomUseCase,
		UserIDFactory: params.UserIDFactory,
		RoomIDFactory: params.RoomIDFactory,
		Logger:        params.Logger,
	}
}
