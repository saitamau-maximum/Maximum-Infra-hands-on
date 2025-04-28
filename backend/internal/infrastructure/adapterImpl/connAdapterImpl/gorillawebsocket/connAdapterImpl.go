package gorillawebsocketconnadapterimpl

import (
	"github.com/gorilla/websocket"
)

type GorillaConnAdapter struct {
	conn *websocket.Conn
}

func NewGorillaConnAdapter(conn *websocket.Conn) *GorillaConnAdapter {
	return &GorillaConnAdapter{conn: conn}
}

// ReadMessageFunc は WebSocketからメッセージを読み取る
func (g *GorillaConnAdapter) ReadMessageFunc() (messageType int, p []byte, err error) {
	return g.conn.ReadMessage()
}

// WriteMessageFunc は WebSocketへメッセージを書き込む
func (g *GorillaConnAdapter) WriteMessageFunc(messageType int, data []byte) error {
	return g.conn.WriteMessage(messageType, data)
}

// CloseFunc は WebSocket接続を閉じる
func (g *GorillaConnAdapter) CloseFunc() error {
	return g.conn.Close()
}
