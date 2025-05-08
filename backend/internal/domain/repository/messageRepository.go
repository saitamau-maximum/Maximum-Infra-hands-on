package repository

import (
	"context"
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

type MessageRepository interface {
	CreateMessage(context.Context, *entity.Message) error
	GetMessageHistoryInRoom(ctx context.Context, roomID entity.RoomID, limit int, beforeSentAt time.Time) (messages []*entity.Message, nextBeforeSentAt time.Time, hasNext bool, err error)
}
