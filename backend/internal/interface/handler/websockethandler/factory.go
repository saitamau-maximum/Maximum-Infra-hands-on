package websockethandler

import (
	"errors"

	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase/websocketcase"
)

type NewWebSocketHandlerParams struct {
	WsUseCase     websocketcase.WebsocketUseCaseInterface
	WsUpgrader    adapter.WebSocketUpgraderAdapter
	WsConnFactory factory.WebSocketConnectionFactory
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
}

func (p *NewWebSocketHandlerParams) Validate() error {
	if p.WsUseCase == nil {
		return errors.New("websocketUseCase is required")
	}
	if p.WsUpgrader == nil {
		return errors.New("websocketUpgrader is required")
	}
	if p.WsConnFactory == nil {
		return errors.New("wsConnFactory is required")
	}
	if p.UserIDFactory == nil {
		return errors.New("userIDFactory is required")
	}
	if p.RoomIDFactory == nil {
		return errors.New("roomIDFactory is required")
	}
	if p.Logger == nil {
		return errors.New("logger is required")
	}
	return nil
}

func NewWebSocketHandler(params NewWebSocketHandlerParams) WebSocketHandlerInterface {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &WebSocketHandler{
		WsUseCase:     params.WsUseCase,
		WsUpgrader:    params.WsUpgrader,
		WsConnFactory: params.WsConnFactory,
		UserIDFactory: params.UserIDFactory,
		RoomIDFactory: params.RoomIDFactory,
		Logger:        params.Logger,
	}
}
