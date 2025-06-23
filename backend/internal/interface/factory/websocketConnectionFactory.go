// アダプターをサービスに変換するファクトリインターフェース
// 具体実装はbackend/internal/infrastructure/factoryImpl/wsConnectionFactoryImpl.go
package factory

import (
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/adapter"
)

type WebSocketConnectionFactory interface {
	CreateWebSocketConnection(conn adapter.ConnAdapter) (service.WebSocketConnection, error)
}
