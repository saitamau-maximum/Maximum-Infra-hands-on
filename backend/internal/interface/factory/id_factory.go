package factory

type UserIDFactory interface {
	CreateUserID() (string, error)
}

type RoomIDFactory interface {
	CreateRoomID() (int, error)
}

type RoomPublicIDFactory interface {
	CreateRoomPublicID() (string, error)
}

type MessageIDFactory interface {
	CreateMessageID() (string, error)
}

type WebsocketClientIDFactory interface {
	CreateWebsocketClientID() (string, error)
}