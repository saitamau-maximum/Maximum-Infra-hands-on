package gorillaconn

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

func (g *GorillaConnAdapter) ReadJSON(message any) error {
	return g.conn.ReadJSON(message)
}

func (g *GorillaConnAdapter) WriteJSON(message any) error {
	return g.conn.WriteJSON(message)
}

// CloseFunc は WebSocket接続を閉じる
func (g *GorillaConnAdapter) CloseFunc() error {
	return g.conn.Close()
}
