package factory

import "example.com/webrtc-practice/internal/domain/entity"

type UserIDFactory interface {
	NewUserID() (entity.UserID, error)
}

type RoomIDFactory interface {
	NewRoomID() (entity.RoomID, error)
}

type RoomPublicIDFactory interface {
	NewRoomPublicID() (entity.RoomPublicID, error)
}

type MessageIDFactory interface {
	NewMessageID() (entity.MessageID, error)
}

type WebsocketClientIDFactory interface {
	NewWebsocketClientID() (entity.WebsocketClientID, error)
}
