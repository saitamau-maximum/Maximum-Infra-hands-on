package service

import "example.com/webrtc-practice/internal/domain/entity"

type WebSocketConnection interface {
	ReadMessage() (int, entity.Message, error)
	WriteMessage(entity.Message) error
	Close() error
}

type WebsocketManager interface {
	// コネクションの登録・削除
	Register(conn *WebSocketConnection, userID entity.UserID, roomID entity.RoomID) error
	Unregister(conn *WebSocketConnection) error
	// コネクションの取得
	GetConnctionByUserID(userID entity.UserID) (*WebSocketConnection, error)

	// 指定した部屋にいるユーザーにブロードキャスト
	BroadcastToRoom(roomID entity.RoomID, msg *entity.Message) error
}