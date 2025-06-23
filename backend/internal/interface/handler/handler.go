package handler

import (
	"example.com/infrahandson/internal/interface/handler/messagehandler"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"example.com/infrahandson/internal/interface/handler/websockethandler"
)

// Handler はアプリケーションのハンドラーをまとめた構造体です。
// DI 層で扱う他の層を組み立てるために必要な構造体を定義します。
type Handler struct {
	UserHandler userhandler.UserHandlerInterface
	RoomHandler roomhandler.RoomHandlerInterface
	WsHandler   websockethandler.WebSocketHandlerInterface
	MsgHandler  messagehandler.MessageHandlerInterface
}
