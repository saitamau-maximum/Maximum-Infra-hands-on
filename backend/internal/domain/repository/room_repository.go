package repository

import "example.com/webrtc-practice/internal/domain/entity"

type RoomRepository interface {
	CreateRoom(name string) (entity.RoomID, error)
	GetRoomByID(id entity.RoomID) (entity.Room, error)
	GetAllRooms() ([]entity.Room, error)
	AddMemberToRoom(entity.RoomID, entity.UserID) error
	GetUsersInRoom(entity.RoomID) ([]entity.UserID, error)
	RemoveMemberFromRoom(entity.RoomID, entity.UserID) error
	DeleteRoom(entity.RoomID) error
	UpdateRoomName(entity.RoomID, string) error
	GetRoomByName(name string) (entity.Room, error)
}
