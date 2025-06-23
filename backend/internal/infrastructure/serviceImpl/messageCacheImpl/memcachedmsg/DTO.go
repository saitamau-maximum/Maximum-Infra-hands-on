package memcachedmsg

import (
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

// MessageDTO は entity.Message のデータ転送オブジェクトです。
type MessageDTO struct {
	ID        string
	RoomID    string
	UserID    string
	Content   string
	CreatedAt time.Time
}

func fromEntityMessage(m *entity.Message) *MessageDTO {
	return &MessageDTO{
		ID:        string(m.GetID()),
		RoomID:    string(m.GetRoomID()),
		UserID:    string(m.GetUserID()),
		Content:   m.GetContent(),
		CreatedAt: m.GetSentAt(),
	}
}

func toEntityMessage(d *MessageDTO) (*entity.Message, error) {
	id := entity.MessageID(d.ID)
	roomID := entity.RoomID(d.RoomID)
	userID := entity.UserID(d.UserID)
	return entity.NewMessage(entity.MessageParams{
		ID:      id,
		RoomID:  roomID,
		UserID:  userID,
		Content: d.Content,
		SentAt:  d.CreatedAt,
	}), nil
}
