package service

import "example.com/webrtc-practice/internal/domain/entity"

type WebSocketConnection interface {
	ReadMessage() (int, *entity.Message, error)
	WriteMessage(*entity.Message) error
	Close() error
}

// 内部的にBroadcaster（Adapter）を使う予定。chan と Redisの差し替えを可能にしたい
type WebsocketManager interface {
	// コネクションの登録・削除
	Register(conn WebSocketConnection, userID entity.UserID, roomID entity.RoomID) error
	Unregister(conn WebSocketConnection) error
	// コネクションの取得
	GetConnectionByUserID(userID entity.UserID) (WebSocketConnection, error)

	// 指定した部屋にいるユーザーにブロードキャスト
	BroadcastToRoom(roomID entity.RoomID, msg *entity.Message) error
}