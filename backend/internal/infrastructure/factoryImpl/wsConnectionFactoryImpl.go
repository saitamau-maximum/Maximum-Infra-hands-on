package factoryimpl

import (
	"example.com/infrahandson/internal/domain/service"
	gorillawsconnectionimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/websocketConnectionImpl/gorillawebsocket"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
)

type WebSocketConnectionFactoryImpl struct{}

func NewWebSocketConnectionFactoryImpl() factory.WebSocketConnectionFactory {
	return &WebSocketConnectionFactoryImpl{}
}

func (f *WebSocketConnectionFactoryImpl) CreateWebSocketConnection(conn adapter.ConnAdapter) (service.WebSocketConnection, error) {
	return gorillawsconnectionimpl.NewGorillaWebSocketConnection(&gorillawsconnectionimpl.NewGorillaWebSocketConnectionParams{
		Conn: conn, // Use exported field
	}), nil
}
