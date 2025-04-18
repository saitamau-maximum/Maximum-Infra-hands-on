package repository

import "example.com/webrtc-practice/internal/domain/entity"

type WebsocketClientRepository interface {
	CreateClient(client *entity.WebsocketClient) error
	DeleteClient(id entity.WebsocketClientID) error
	GetClientByID(id entity.WebsocketClientID) (*entity.WebsocketClient, error)
	GetClientsByRoomID(roomID entity.RoomID) ([]*entity.WebsocketClient, error)
	GetClientsByUserID(userID entity.UserID) (*entity.WebsocketClient, error)
}