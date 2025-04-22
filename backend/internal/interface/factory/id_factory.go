package factory

import "example.com/webrtc-practice/internal/domain/entity"

type UserIDFactory interface {
	NewUserPublicID() (entity.UserPublicID, error)
	FromInt(int) (entity.UserID, error)
	FromString(string) (entity.UserPublicID, error)
}

type RoomIDFactory interface {
	NewRoomPublicID() (entity.RoomPublicID, error)
	FromInt(int) (entity.RoomID, error)
	FromString(string) (entity.RoomPublicID, error)
}

type MessageIDFactory interface {
	NewMessagePublicID() (entity.MessagePublicID, error)
	FromString(string) (entity.MessagePublicID, error)
	FromInt(int) (entity.MessagePublicID, error)
}

type WsClientIDFactory interface {
	NewWsClientPublicID() (entity.WsClientPublicID, error)
	FromInt(int) (entity.WsClientPublicID, error)
	FromString(string) (entity.WsClientPublicID, error)
}
