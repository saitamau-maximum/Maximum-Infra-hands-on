package gorillaupgrader

import (
	"net/http"

	"example.com/infrahandson/internal/infrastructure/adapterImpl/connAdapterImpl/gorillaconn"
	"example.com/infrahandson/internal/interface/adapter"
	"github.com/gorilla/websocket"
)

// GorillaWebSocketUpgrader は gorilla/websocket を使った WebSocketUpgraderAdapter の実装です。
type GorillaWebSocketUpgrader struct {
	upgrader websocket.Upgrader
}

// NewGorillaWebSocketUpgrader は GorillaWebSocketUpgrader を初期化して返します。
func NewGorillaWebSocketUpgrader() *GorillaWebSocketUpgrader {
	return &GorillaWebSocketUpgrader{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 必要に応じてオリジンチェックを実装
			},
		},
	}
}

// Upgrade はHTTP接続をWebSocket接続にアップグレードします。
func (g *GorillaWebSocketUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (adapter.ConnAdapter, error) {
	conn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return gorillaconn.NewGorillaConnAdapter(conn), nil
}
