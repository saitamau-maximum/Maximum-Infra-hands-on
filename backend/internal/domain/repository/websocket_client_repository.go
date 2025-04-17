package repository

import "example.com/webrtc-practice/internal/domain/entity"

type IWebsocketClientRepository interface {
	CreateClient(id entity.WebsocketClientID) error
	DeleteClient(id entity.WebsocketClientID) error
	GetClientByID(id entity.WebsocketClientID) (*entity.WebsocketClient, error)
	GetClientsByRoomID(roomID entity.RoomID) ([]*entity.WebsocketClient, error)
}