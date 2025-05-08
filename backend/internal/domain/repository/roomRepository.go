package repository

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

type RoomRepository interface {
	SaveRoom(context.Context, *entity.Room) (entity.RoomID, error)
	GetRoomByID(ctx context.Context, id entity.RoomID) (*entity.Room, error)
	GetAllRooms(context.Context) ([]*entity.Room, error)
	GetUsersInRoom(context.Context, entity.RoomID) ([]*entity.User, error)
	AddMemberToRoom(context.Context, entity.RoomID, entity.UserID) error
	RemoveMemberFromRoom(context.Context, entity.RoomID, entity.UserID) error
	GetRoomByNameLike(ctx context.Context, name string) ([]*entity.Room, error)
	UpdateRoomName(context.Context, entity.RoomID, string) error
	DeleteRoom(context.Context, entity.RoomID) error
}
