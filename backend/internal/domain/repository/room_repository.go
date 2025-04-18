package repository

import "example.com/webrtc-practice/internal/domain/entity"

type RoomRepository interface {
	SaveRoom(*entity.Room) (entity.RoomID, error)
	GetRoomByID(id entity.RoomID) (*entity.Room, error)
	GetAllRooms() ([]*entity.Room, error)
	GetUsersInRoom(entity.RoomID) ([]*entity.User, error)
	AddMemberToRoom(entity.RoomID, entity.UserID) error
	RemoveMemberFromRoom(entity.RoomID, entity.UserID) error
	GetRoomByNameLike(name string) ([]*entity.Room, error)
	UpdateRoomName(entity.RoomID, string) error
	DeleteRoom(entity.RoomID) error
	GetRoomIDByPublicID(id entity.RoomPublicID) (entity.RoomID, error)
	GetPublicIDByRoomID(id entity.RoomID) (entity.RoomPublicID, error)
}
