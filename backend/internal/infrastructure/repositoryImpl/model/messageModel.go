package model

import (
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"github.com/google/uuid"
)

type MessageModel struct {
	ID      uuid.UUID `db:"id"`
	RoomID  uuid.UUID `db:"room_id"`
	UserID  uuid.UUID `db:"user_id"`
	Content string    `db:"content"`
	SentAt  time.Time `db:"sent_at"`
}

func (m *MessageModel) FromEntity(message *entity.Message) error {
	id := message.GetID()
	idUUID, err := id.MessageID2UUID()
	if err != nil {
		return err
	}
	m.ID = idUUID
	roomID := message.GetRoomID()
	roomIDUUID, err := roomID.RoomID2UUID()
	if err != nil {
		return err
	}
	m.RoomID = roomIDUUID
	userID := message.GetUserID()
	userIDUUID, err := userID.UserID2UUID()
	if err != nil {
		return err
	}
	m.UserID = userIDUUID
	m.Content = message.GetContent()
	m.SentAt = message.GetSentAt()
	return nil
}

func (m *MessageModel) ToEntity() *entity.Message {
	return entity.NewMessage(entity.MessageParams{
		ID:      entity.MessageID(m.ID.String()), // UUID -> MessageID
		RoomID:  entity.RoomID(m.RoomID.String()),
		UserID:  entity.UserID(m.UserID.String()),
		Content: m.Content,
		SentAt:  m.SentAt,
	})
}
