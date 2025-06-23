package factory

type Factory struct {
	// 各種IDを生成するファクトリー
	UserIDFactory     UserIDFactory
	RoomIDFactory     RoomIDFactory
	MessageIDFactory  MessageIDFactory
	WsClientIDFactory WsClientIDFactory

	// WebSocket接続を生成するファクトリー
	WsConnFactory     WebSocketConnectionFactory
}
