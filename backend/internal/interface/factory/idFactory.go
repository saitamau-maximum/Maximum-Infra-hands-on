package factory

import "example.com/infrahandson/internal/domain/entity"

type UserIDFactory interface {
	NewUserID() (entity.UserID, error)
}

type RoomIDFactory interface {
	NewRoomID() (entity.RoomID, error)
}

type MessageIDFactory interface {
	NewMessageID() (entity.MessageID, error)
}

type WsClientIDFactory interface {
	NewWsClientID() (entity.WsClientID, error)
}
