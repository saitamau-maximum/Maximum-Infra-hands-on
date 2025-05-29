// httpリクエストをWebSocket接続にアップグレードするためのアダプターインターフェースです。
// 具体実装はbackend/internal/infrastructure/adapterImpl/upgraderAdapterImpl
package adapter

import (
	"net/http"
)

// WebSocketUpgraderAdapter は WebSocket の接続アップグレードを行うインターフェースです。
type WebSocketUpgraderAdapter interface {
	// Upgrade はWebSocket接続をアップグレードします。
	Upgrade(w http.ResponseWriter, r *http.Request) (ConnAdapter, error)
}
