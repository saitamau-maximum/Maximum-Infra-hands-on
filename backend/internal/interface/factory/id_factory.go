package factory

import "example.com/webrtc-practice/internal/domain/entity"

type UserIDFactory interface {
	NewUserPublicID() (entity.UserPublicID, error)
	FromInt(int) entity.UserID
	FromString(string) entity.UserPublicID
}

type RoomIDFactory interface {
	NewRoomPublicID() (entity.RoomPublicID, error)
	FromInt(int) entity.RoomID
	FromString(string) entity.RoomPublicID
}

type MessageIDFactory interface {
	NewMessagePublicID() (entity.MessagePublicID, error)
	FromString(string) entity.MessagePublicID
	FromInt(int) entity.MessageID
}

type WsClientIDFactory interface {
	NewWsClientPublicID() (entity.WsClientPublicID, error)
	FromInt(int) entity.WsClientID
	FromString(string) entity.WsClientPublicID
}
