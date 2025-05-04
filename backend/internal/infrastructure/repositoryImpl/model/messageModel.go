package model

import (
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

type MessageModel struct {
	ID      string    `db:"id"`
	RoomID  string    `db:"room_id"`
	UserID  string    `db:"user_id"`
	Content string    `db:"content"`
	SentAt  time.Time `db:"sent_at"`
}

func (m *MessageModel) FromEntity(message *entity.Message) {
	m.ID = string(message.GetID())
	m.RoomID = string(message.GetRoomID())
	m.UserID = string(message.GetUserID())
	m.Content = message.GetContent()
	m.SentAt = message.GetSentAt()
}

func (m *MessageModel) ToEntity() *entity.Message {
	return entity.NewMessage(entity.MessageParams{
		ID:      entity.MessageID(m.ID),
		RoomID:  entity.RoomID(m.RoomID),
		UserID:  entity.UserID(m.UserID),
		Content: m.Content,
		SentAt:  m.SentAt,
	})
}
