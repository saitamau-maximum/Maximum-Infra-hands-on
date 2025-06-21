package repository

type Repository struct {
	UserRepository     UserRepository
	RoomRepository     RoomRepository
	MessageRepository  MessageRepository
	WsClientRepository WebsocketClientRepository
}