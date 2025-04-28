package repository

import "example.com/infrahandson/internal/domain/entity"

type WebsocketClientRepository interface {
	CreateClient(client *entity.WebsocketClient) error
	DeleteClient(id entity.WsClientID) error
	GetClientByID(id entity.WsClientID) (*entity.WebsocketClient, error)
	GetClientsByRoomID(roomID entity.RoomID) ([]*entity.WebsocketClient, error)
	GetClientsByUserID(userID entity.UserID) (*entity.WebsocketClient, error)
}
