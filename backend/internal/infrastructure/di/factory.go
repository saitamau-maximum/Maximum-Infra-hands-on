package di

import (
	"example.com/infrahandson/internal/infrastructure/factoryImpl"
	"example.com/infrahandson/internal/interface/factory"
)

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
