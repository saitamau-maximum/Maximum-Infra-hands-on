package factory

import "example.com/webrtc-practice/internal/domain/entity"

type UserIDFactory interface {
	NewUserID() (entity.UserID, error)
	FromString(string) entity.UserID
}

type RoomIDFactory interface {
	NewRoomID() (entity.RoomID, error)
	FromInt(int) entity.RoomID
}

type RoomPublicIDFactory interface {
	NewRoomPublicID() (entity.RoomPublicID, error)
	FromString(string) entity.RoomPublicID
}

type MessageIDFactory interface {
	NewMessageID() (entity.MessageID, error)
	FromString(string) entity.MessageID
}

type WebsocketClientIDFactory interface {
	NewWebsocketClientID() (entity.WebsocketClientID, error)
	FromString(string) entity.WebsocketClientID
}
