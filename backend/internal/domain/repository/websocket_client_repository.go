package repository

import "example.com/webrtc-practice/internal/domain/entity"

type WebsocketClientRepository interface {
	CreateClient(entity.WebsocketClient) error
	DeleteClient(id entity.WebsocketClientID) error
	GetClientByID(id entity.WebsocketClientID) (*entity.WebsocketClient, error)
	GetClientsByRoomID(roomID entity.RoomID) ([]*entity.WebsocketClient, error)
}