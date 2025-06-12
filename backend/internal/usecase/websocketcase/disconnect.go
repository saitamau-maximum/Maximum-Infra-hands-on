package websocketcase

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// DisconnectUserRequest構造体: 切断処理リクエスト
type DisconnectUserRequest struct {
	UserID entity.UserID
}

// DisconnectUser 切断処理
func (w *WebsocketUseCase) DisconnectUser(ctx context.Context, req DisconnectUserRequest) error {
	conn, err := w.websocketManager.GetConnectionByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	user, err := w.wsClientRepo.GetClientsByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	err = w.websocketManager.Unregister(ctx, conn)
	if err != nil {
		return err
	}

	err = w.wsClientRepo.DeleteClient(ctx, user.GetID())
	if err != nil {
		return err
	}

	return nil
}
