package sqlite3

import (
	"errors"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type MessageRepositoryImpl struct {
	db *sqlx.DB
}

type NewMessageRepositoryParams struct {
	DB *sqlx.DB
}

func (p *NewMessageRepositoryParams) Validate() error {
	if p.DB == nil {
		return errors.New("DB is required")
	}
	return nil
}

func NewMessageRepository(params NewMessageRepositoryParams) repository.MessageRepository {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &MessageRepositoryImpl{
		db: params.DB,
	}
}

type MessageModel struct {
	ID      string    `db:"id"`
	RoomID  int       `db:"room_id"`
	UserID  string    `db:"user_id"`
	Content string    `db:"content"`
	SentAt  time.Time `db:"sent_at"`
}

func (r *MessageRepositoryImpl) CreateMessage(message *entity.Message) error {
	query := `INSERT INTO messages (id, room_id, sender_id, content, sent_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, message.GetID(), int(message.GetRoomID()), string(message.GetUserID()), message.GetContent(), message.GetSentAt())
	return err
}

func (r *MessageRepositoryImpl) GetMessagesByRoomID(roomID entity.RoomID) ([]*entity.Message, error) {
	query := `SELECT id, room_id, sender_id, content, sent_at FROM messages WHERE room_id = ? ORDER BY sent_at ASC`
	var messageModels []MessageModel
	err := r.db.Select(&messageModels, query, roomID)
	if err != nil {
		return nil, err
	}

	var messages []*entity.Message
	for _, model := range messageModels {
		messages = append(messages, entity.NewMessage(entity.MessageParams{
			ID:      entity.MessageID(model.ID),
			RoomID:  entity.RoomID(model.RoomID),
			UserID:  entity.UserID(model.UserID),
			Content: model.Content,
			SentAt:  model.SentAt,
		}))
	}
	return messages, nil
}

func (r *MessageRepositoryImpl) GetMessageHistoryInRoom(roomID entity.RoomID, limit int, beforeSentAt time.Time) (messages []*entity.Message, nextBeforeSentAt time.Time, hasNext bool, err error) {
	query := `SELECT id, room_id, sender_id, content, sent_at FROM messages WHERE room_id = ? AND sent_at < ? ORDER BY sent_at DESC LIMIT ?`
	var messageModels []MessageModel
	err = r.db.Select(&messageModels, query, roomID, beforeSentAt, limit)
	if err != nil {
		return nil, time.Time{}, false, err
	}

	for _, model := range messageModels {
		messages = append(messages, entity.NewMessage(entity.MessageParams{
			ID:       entity.MessageID(model.ID),
			RoomID:   entity.RoomID(model.RoomID),
			UserID: entity.UserID(model.UserID),
			Content:  model.Content,
			SentAt:   model.SentAt,
		}))
	}

	if len(messageModels) > 0 {
		nextBeforeSentAt = messageModels[len(messageModels)-1].SentAt
		hasNext = len(messageModels) == limit
	}

	return messages, nextBeforeSentAt, hasNext, nil
}
