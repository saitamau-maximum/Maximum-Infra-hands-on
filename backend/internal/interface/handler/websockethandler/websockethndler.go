package websockethandler

import (
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase/websocketcase"
)

type WebSocketHandler struct {
	WsUseCase     websocketcase.WebsocketUseCaseInterface
	WsUpgrader    adapter.WebSocketUpgraderAdapter
	WsConnFactory factory.WebSocketConnectionFactory
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
}
