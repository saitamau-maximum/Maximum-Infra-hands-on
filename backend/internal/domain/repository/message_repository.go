package repository

import "example.com/webrtc-practice/internal/domain/entity"

type MessageRepository interface {
	CreateMessage(roomID string, userID string, content string) (string, error)
	GetMessagesByRoomID(roomID string) ([]entity.Message, error)
}
