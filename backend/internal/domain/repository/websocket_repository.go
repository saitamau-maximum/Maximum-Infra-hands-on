package repository

import "example.com/webrtc-practice/internal/domain/entity"

type IWebsocketRepository interface {
	CreateClient(id entity.WebsocketClientID) error
	DeleteClient(id entity.WebsicketClientID) error
	GetClientByID(id entity.WebsicketClientID) (*entity.WebsocketClient, error)
	GetClientsByRoomID(roomID entity.RoomID) ([]*entity.WebsocketClient, error)
}