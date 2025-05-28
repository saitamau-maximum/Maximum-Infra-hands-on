// Roomエンティティの永続化のためのインターフェース
// infrastructure/repositoryImplで具体実装
package repository

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// RoomRepository defines the interface for managing chat rooms and their members.
type RoomRepository interface {
	// SaveRoom persists a new room and returns its ID.
	SaveRoom(ctx context.Context, room *entity.Room) (entity.RoomID, error)

	// GetRoomByID retrieves a room by its unique ID.
	GetRoomByID(ctx context.Context, id entity.RoomID) (*entity.Room, error)

	// GetAllRooms returns a list of all rooms.
	GetAllRooms(ctx context.Context) ([]*entity.Room, error)

	// GetUsersInRoom retrieves all users who are members of the specified room.
	GetUsersInRoom(ctx context.Context, roomID entity.RoomID) ([]*entity.User, error)

	// AddMemberToRoom adds a user to the specified room.
	AddMemberToRoom(ctx context.Context, roomID entity.RoomID, userID entity.UserID) error

	// RemoveMemberFromRoom removes a user from the specified room.
	RemoveMemberFromRoom(ctx context.Context, roomID entity.RoomID, userID entity.UserID) error

	// GetRoomByNameLike performs a partial match search for room names.
	GetRoomByNameLike(ctx context.Context, name string) ([]*entity.Room, error)

	// UpdateRoomName updates the name of the specified room.
	UpdateRoomName(ctx context.Context, roomID entity.RoomID, name string) error

	// DeleteRoom deletes the specified room and its associations.
	DeleteRoom(ctx context.Context, roomID entity.RoomID) error
}

