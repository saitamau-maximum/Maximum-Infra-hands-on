package model

import (
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

type MessageModel struct {
	ID       int       `db:"id"`
	PublicID string    `db:"public_id"`
	RoomID   int       `db:"room_id"`
	UserID   int       `db:"user_id"`
	Content  string    `db:"content"`
	SentAt   time.Time `db:"sent_at"`
}

func (m *MessageModel) FromEntity(message *entity.Message) {
	m.PublicID = string(message.GetPublicID())
	m.RoomID = int(message.GetRoomID())
	m.UserID = int(message.GetUserID())
	m.Content = message.GetContent()
	m.SentAt = message.GetSentAt()
}

func (m *MessageModel) ToEntity() *entity.Message {
	return entity.NewMessage(entity.MessageParams{
		ID:       entity.MessageID(m.ID),
		PublicID: entity.MessagePublicID(m.PublicID),
		RoomID:   entity.RoomID(m.RoomID),
		UserID:   entity.UserID(m.UserID),
		Content:  m.Content,
		SentAt:   m.SentAt,
	})
}
