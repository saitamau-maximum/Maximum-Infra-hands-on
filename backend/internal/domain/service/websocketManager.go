// 部屋に属するコネクションの管理とメッセージ送信ロジックのインターフェース
package service

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// コネクションの抽象化
type WebSocketConnection interface {
	ReadMessage() (*entity.Message, error)
	WriteMessage(*entity.Message) error
	Close() error
}

type WebsocketManager interface {
	// コネクションの登録・削除
	Register(ctx context.Context, conn WebSocketConnection, userID entity.UserID, roomID entity.RoomID) error
	Unregister(ctx context.Context, conn WebSocketConnection) error
	// コネクションの取得
	GetConnectionByUserID(ctx context.Context, userID entity.UserID) (WebSocketConnection, error)

	// 指定した部屋にいるユーザーにブロードキャスト
	BroadcastToRoom(ctx context.Context, roomID entity.RoomID, msg *entity.Message) error
}
