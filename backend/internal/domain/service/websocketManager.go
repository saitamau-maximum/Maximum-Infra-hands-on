package service

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

type WebSocketConnection interface {
	ReadMessage() (*entity.Message, error)
	WriteMessage(*entity.Message) error
	Close() error
}

// 内部的にBroadcaster（Adapter）を使う予定。chan と Redisの差し替えを可能にしたい
type WebsocketManager interface {
	// コネクションの登録・削除
	Register(ctx context.Context, conn WebSocketConnection, userID entity.UserID, roomID entity.RoomID) error
	Unregister(ctx context.Context, conn WebSocketConnection) error
	// コネクションの取得
	GetConnectionByUserID(ctx context.Context, userID entity.UserID) (WebSocketConnection, error)

	// 指定した部屋にいるユーザーにブロードキャスト
	BroadcastToRoom(ctx context.Context, roomID entity.RoomID, msg *entity.Message) error
}
