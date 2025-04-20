// interface/factory/websocket_connection_factory.go
package factory

import (
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/interface/adapter"
)

type WebSocketConnectionFactory interface {
	CreateWebSocketConnection(conn adapter.ConnAdapter) (service.WebSocketConnection, error)
}
