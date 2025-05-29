package websocket

import "context"

type WebsocketUseCaseInterface interface {
	// ConnectUserToRoom: 接続・参加処理
	ConnectUserToRoom(ctx context.Context, req ConnectUserToRoomRequest) error

	// SendMessage: メッセージ送信
	SendMessage(ctx context.Context, req SendMessageRequest) error

	// DisconnectUser: 切断処理
	DisconnectUser(ctx context.Context, req DisconnectUserRequest) error
}
