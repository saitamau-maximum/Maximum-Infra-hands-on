package sqlitemsgrepoimpl

import (
	"context"
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/model"
	"github.com/jmoiron/sqlx"
)

type MessageRepositoryImpl struct {
	DB *sqlx.DB
}
type NewMessageRepositoryImplParams struct {
	DB *sqlx.DB
}

func (p *NewMessageRepositoryImplParams) Validate() error {
	if p.DB == nil {
		return errors.New("DB is nil")
	}
	return nil
}

func NewMessageRepositoryImpl(params *NewMessageRepositoryImplParams) repository.MessageRepository {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &MessageRepositoryImpl{
		DB: params.DB,
	}
}

func (r *MessageRepositoryImpl) CreateMessage(ctx context.Context, message *entity.Message) error {
	var Message model.MessageModel
	Message.FromEntity(message)

	_, err := r.DB.NamedExecContext(ctx, "INSERT INTO messages (id, room_id, user_id, content, sent_at) VALUES (:id, :room_id, :user_id, :content, :sent_at)", &Message)
	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepositoryImpl) GetMessageHistoryInRoom(
	ctx context.Context,
	roomID entity.RoomID,
	limit int,
	beforeSentAt time.Time,
) (messages []*entity.Message, nextBeforeSentAt time.Time, hasNext bool, err error) {
	var MessageModels []model.MessageModel
	query := "SELECT * FROM messages WHERE room_id = ? AND sent_at < ? ORDER BY sent_at DESC LIMIT ?"
	err = r.DB.SelectContext(ctx, &MessageModels, query, roomID, beforeSentAt, limit)
	if err != nil {
		return nil, time.Now(), false, err
	}

	if len(MessageModels) == 0 {
		return nil, time.Now(), false, nil
	}

	messages = make([]*entity.Message, len(MessageModels))
	for i := range MessageModels {
		messages[i] = MessageModels[i].ToEntity()
	}

	nextBeforeSentAt = MessageModels[len(MessageModels)-1].SentAt
	hasNext = len(messages) == limit

	return messages, nextBeforeSentAt, hasNext, nil
}
