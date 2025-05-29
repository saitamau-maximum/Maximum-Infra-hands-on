package websocket

import (
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/factory"
)

type WebsocketUseCase struct {
	userRepo         repository.UserRepository
	roomRepo         repository.RoomRepository
	msgRepo          repository.MessageRepository
	msgCache         service.MessageCacheService
	wsClientRepo     repository.WebsocketClientRepository
	websocketManager service.WebsocketManager
	msgIDFactory     factory.MessageIDFactory
	clientIDFactory  factory.WsClientIDFactory
}
