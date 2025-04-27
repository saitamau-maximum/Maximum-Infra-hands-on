package repository

import (
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

type MessageRepository interface {
	CreateMessage(*entity.Message) error
	GetMessagesByRoomID(roomID entity.RoomID) ([]*entity.Message, error)
	GetMessageHistoryInRoom(roomID entity.RoomID, limit int, beforeSentAt time.Time) (messages []*entity.Message, nextBeforeSentAt time.Time, hasNext bool, err error)
}
