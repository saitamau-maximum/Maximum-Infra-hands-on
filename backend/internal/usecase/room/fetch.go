package room

import (
	"context"
	"errors"

	"example.com/infrahandson/internal/domain/entity"
)

// GetRoomByIDParams構造体: 公開IDで部屋を取得するためのパラメータ
type GetRoomByIDRequest struct {
	ID entity.RoomID `json:"id"`
}

// GetRoomByIDResponse構造体: 公開IDで部屋を取得した結果
type GetRoomByIDResponse struct {
	Room *entity.Room `json:"room"` // 取得した部屋
}

// GetRoomByID: 公開IDを使用して部屋を取得
func (r *RoomUseCase) GetRoomByID(ctx context.Context, req GetRoomByIDRequest) (GetRoomByIDResponse, error) {
	room, err := r.roomRepo.GetRoomByID(ctx, req.ID)
	if err != nil {
		return GetRoomByIDResponse{}, err
	}
	if room == nil {
		return GetRoomByIDResponse{}, errors.New("room not found")
	}
	return GetRoomByIDResponse{Room: room}, nil
}

// GetAllRooms: 全ての部屋を取得
func (r *RoomUseCase) GetAllRooms(ctx context.Context) ([]*entity.Room, error) {
	rooms, err := r.roomRepo.GetAllRooms(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetUsersInRoomRequest構造体: 部屋内のユーザーを取得するリクエスト
type GetUsersInRoomRequest struct {
	ID entity.RoomID `json:"id"` // 公開ID
}

// GetUsersInRoomResponse構造体: 部屋内のユーザー取得結果
type GetUsersInRoomResponse struct {
	Users []*entity.User `json:"users"` // 部屋内のユーザーリスト
}

// GetUsersInRoom: 部屋内のユーザーを取得
func (r *RoomUseCase) GetUsersInRoom(ctx context.Context, req GetUsersInRoomRequest) (GetUsersInRoomResponse, error) {
	users, err := r.roomRepo.GetUsersInRoom(ctx, req.ID)
	if err != nil {
		return GetUsersInRoomResponse{}, err
	}

	return GetUsersInRoomResponse{Users: users}, nil
}