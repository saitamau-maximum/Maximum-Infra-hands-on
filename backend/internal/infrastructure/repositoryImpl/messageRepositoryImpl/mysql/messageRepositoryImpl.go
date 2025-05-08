package mysqlmsgrepoimpl

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
	db *sqlx.DB
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
		db: params.DB,
	}
}

func (r *MessageRepositoryImpl) CreateMessage(ctx context.Context, message *entity.Message) error {
	var msg model.MessageModel
	msg.FromEntity(message)

	// UUIDを文字列で扱い、DB側でUUID_TO_BINに変換
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO messages (id, room_id, user_id, content, sent_at)
		VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?)`,
		msg.ID, msg.RoomID, msg.UserID, msg.Content, msg.SentAt)

	return err
}

func (r *MessageRepositoryImpl) GetMessageHistoryInRoom(
	ctx context.Context,
	roomID entity.RoomID,
	limit int,
	beforeSentAt time.Time,
) (messages []*entity.Message, nextBeforeSentAt time.Time, hasNext bool, err error) {
	var msgModels []model.MessageModel

	query := `
		SELECT 
			BIN_TO_UUID(id) AS id,
			BIN_TO_UUID(room_id) AS room_id,
			BIN_TO_UUID(user_id) AS user_id,
			content,
			sent_at
		FROM messages
		WHERE room_id = UUID_TO_BIN(?) AND sent_at < ?
		ORDER BY sent_at DESC
		LIMIT ?`

	// RoomID -> UUID
	roomIDUUID, err := roomID.RoomID2UUID()
	if err != nil {
		return nil, time.Now(), false, err
	}
	err = r.db.SelectContext(ctx, &msgModels, query, roomIDUUID, beforeSentAt, limit)
	if err != nil {
		return nil, time.Now(), false, err
	}

	if len(msgModels) == 0 {
		return nil, time.Now(), false, nil
	}

	messages = make([]*entity.Message, len(msgModels))
	for i := range msgModels {
		messages[i] = msgModels[i].ToEntity()
	}

	nextBeforeSentAt = msgModels[len(msgModels)-1].SentAt
	hasNext = len(msgModels) == limit

	return messages, nextBeforeSentAt, hasNext, nil
}
