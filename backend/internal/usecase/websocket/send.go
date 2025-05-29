package websocket

import (
	"context"
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

// SendMessageRequest構造体: メッセージ送信リクエスト
type SendMessageRequest struct {
	RoomID  entity.RoomID
	Sender  entity.UserID
	Content string
}

// SendMessage メッセージ送信
func (w *WebsocketUseCase) SendMessage(ctx context.Context, req SendMessageRequest) error {
	id, err := w.msgIDFactory.NewMessageID()
	if err != nil {
		return err
	}

	msg := entity.NewMessage(entity.MessageParams{
		ID:      id,
		RoomID:  req.RoomID,
		UserID:  req.Sender,
		Content: req.Content,
		SentAt:  time.Now(),
	})

	if err := w.msgRepo.CreateMessage(ctx, msg); err != nil {
		return err
	}

	if err := w.msgCache.AddMessage(ctx, req.RoomID, msg); err != nil {
		return err
	}

	err = w.websocketManager.BroadcastToRoom(ctx, req.RoomID, msg)
	if err != nil {
		return err
	}

	return nil
}
