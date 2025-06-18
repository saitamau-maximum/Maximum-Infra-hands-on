package factory

type Factory struct {
	UserIDFactory     UserIDFactory
	RoomIDFactory     RoomIDFactory
	MessageIDFactory  MessageIDFactory
	WsClientIDFactory WsClientIDFactory
	WsConnFactory     WebSocketConnectionFactory
}
