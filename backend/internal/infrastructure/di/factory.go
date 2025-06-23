package di

import (
	"example.com/infrahandson/internal/infrastructure/factoryImpl"
	"example.com/infrahandson/internal/interface/factory"
)

// InitializeFactory はファクトリーの初期化を行います。
// 返り値 factory.Factory はファクトリー層をまとめた構造体（詳細：internal/interface/factory/factory.go）です。
func InitializeFactory() factory.Factory {
	// Factoryの初期化
	userIDFactory := factoryimpl.NewUserIDFactory()
	roomIDFactory := factoryimpl.NewRoomIDFactory()
	MsgIDFactory := factoryimpl.NewMessageIDFactory()
	clientDFactory := factoryimpl.NewWsClientIDFactory()
	wsConnFactory := factoryimpl.NewWebSocketConnectionFactoryImpl()

	return factory.Factory{
		UserIDFactory:     userIDFactory,
		RoomIDFactory:     roomIDFactory,
		MessageIDFactory:  MsgIDFactory,
		WsClientIDFactory: clientDFactory,
		WsConnFactory:     wsConnFactory,
	}
}
