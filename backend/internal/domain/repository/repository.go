package repository

// Repository : repositoryのインターフェースをまとめた構造体
// DI層での依存性注入のために使用される
type Repository struct {
	UserRepository     UserRepository
	RoomRepository     RoomRepository
	MessageRepository  MessageRepository
	WsClientRepository WebsocketClientRepository
}