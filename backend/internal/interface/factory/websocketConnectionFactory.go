// interface/factory/websocket_connection_factory.go
package factory

import (
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/adapter"
)

type WebSocketConnectionFactory interface {
	CreateWebSocketConnection(conn adapter.ConnAdapter) (service.WebSocketConnection, error)
}
