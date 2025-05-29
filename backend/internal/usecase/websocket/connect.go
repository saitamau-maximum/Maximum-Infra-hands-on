package websocket

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
)

// ConnectUserToRoomRequest構造体: 接続・参加処理のリクエスト
type ConnectUserToRoomRequest struct {
	UserID entity.UserID
	RoomID entity.RoomID
	Conn   service.WebSocketConnection
}

// ConnectUserToRoom 接続・参加処理
func (w *WebsocketUseCase) ConnectUserToRoom(ctx context.Context, req ConnectUserToRoomRequest) error {
	user, err := w.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	id, err := w.clientIDFactory.NewWsClientID()
	if err != nil {
		return err
	}

	client := entity.NewWebsocketClient(entity.WebsocketClientParams{
		ID:     id,
		UserID: user.GetID(),
		RoomID: req.RoomID,
	})

	err = w.wsClientRepo.CreateClient(ctx, client)
	if err != nil {
		return err
	}

	err = w.websocketManager.Register(ctx, req.Conn, req.UserID, req.RoomID)
	if err != nil {
		return err
	}

	return nil
}
