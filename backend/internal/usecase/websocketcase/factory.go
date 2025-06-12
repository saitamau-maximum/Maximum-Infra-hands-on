package websocketcase

import (
	"errors"

	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/factory"
)

type NewWebsocketUseCaseParams struct {
	UserRepo         repository.UserRepository
	RoomRepo         repository.RoomRepository
	MsgRepo          repository.MessageRepository
	MsgCache         service.MessageCacheService
	WsClientRepo     repository.WebsocketClientRepository
	WebsocketManager service.WebsocketManager
	MsgIDFactory     factory.MessageIDFactory
	ClientIDFactory  factory.WsClientIDFactory
}

func (p *NewWebsocketUseCaseParams) Validate() error {
	if p.UserRepo == nil {
		return errors.New("UserRepo is required")
	}
	if p.RoomRepo == nil {
		return errors.New("RoomRepo is required")
	}
	if p.MsgRepo == nil {
		return errors.New("MsgRepo is required")
	}
	if p.MsgCache == nil {
		return errors.New("MsgCache is required")
	}
	if p.WsClientRepo == nil {
		return errors.New("WsClientRepo is required")
	}
	if p.WebsocketManager == nil {
		return errors.New("WebsocketManager is required")
	}
	if p.MsgIDFactory == nil {
		return errors.New("MsgIDFactory is required")
	}
	if p.ClientIDFactory == nil {
		return errors.New("ClientIDFactory is required")
	}
	return nil
}

func NewWebsocketUseCase(params NewWebsocketUseCaseParams) WebsocketUseCaseInterface {
	// Paramsのバリデーションを行う
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &WebsocketUseCase{
		userRepo:         params.UserRepo,
		roomRepo:         params.RoomRepo,
		msgRepo:          params.MsgRepo,
		msgCache:         params.MsgCache,
		wsClientRepo:     params.WsClientRepo,
		websocketManager: params.WebsocketManager,
		msgIDFactory:     params.MsgIDFactory,
		clientIDFactory:  params.ClientIDFactory,
	}
}
