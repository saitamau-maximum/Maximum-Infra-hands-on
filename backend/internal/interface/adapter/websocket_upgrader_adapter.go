package adapter

import (
	"net/http"
)

// WebSocketUpgraderAdapter は WebSocket の接続アップグレードを行うインターフェースです。
type WebSocketUpgraderAdapter interface {
	// WebSocket接続をアップグレードします。
	Upgrade(w http.ResponseWriter, r *http.Request) (ConnAdapter, error)
}
