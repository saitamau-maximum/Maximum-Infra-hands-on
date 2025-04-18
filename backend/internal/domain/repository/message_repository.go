package repository

import "example.com/webrtc-practice/internal/domain/entity"

type MessageRepository interface {
	CreateMessage(entity.Message) (string, error)
	GetMessagesByRoomID(roomID entity.RoomID) ([]*entity.Message, error)
}
