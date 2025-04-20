package repository

import "example.com/webrtc-practice/internal/domain/entity"

type MessageRepository interface {
	CreateMessage(*entity.Message) error
	GetMessagesByRoomID(roomID entity.RoomID) ([]*entity.Message, error)
}
