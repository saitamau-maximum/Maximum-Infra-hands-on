package websockethandler

import "github.com/labstack/echo/v4"

type WebSocketHandlerInterface interface {
	// ConnectToChatRoom はWebSocket接続を確立し．ソケットの監視ゴルーチンを生成する
	ConnectToChatRoom(c echo.Context) error
}
