// room usecaseのinterface
package room

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

type RoomUseCaseInterface interface {
	// CreateRoom は新しい部屋を作成する(create.go)
	CreateRoom(ctx context.Context, req CreateRoomRequest) (CreateRoomResponse, error)

	// GetRoomByID は公開IDを使用して部屋を取得する(get.go)
	GetRoomByID(ctx context.Context, params GetRoomByIDRequest) (GetRoomByIDResponse, error)

	// GetAllRooms は全ての部屋を取得する(get.go)
	GetAllRooms(ctx context.Context) ([]*entity.Room, error)

	// GetUsersInRoom は部屋内のユーザーを取得する(get.go)
	GetUsersInRoom(ctx context.Context, req GetUsersInRoomRequest) (GetUsersInRoomResponse, error)

	// JoinRoom は部屋にユーザーを参加させる(membership.go)
	JoinRoom(ctx context.Context, req JoinRoomRequest) error
	
	// LeaveRoom は部屋からユーザーを退出させる(membership.go)
	LeaveRoom(ctx context.Context, req LeaveRoomRequest) error
}
