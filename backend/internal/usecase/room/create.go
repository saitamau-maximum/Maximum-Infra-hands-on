package room

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// CreateRoomRequest構造体: 部屋作成リクエストのデータ
type CreateRoomRequest struct {
	Name string `json:"name"` // 部屋名
}

// CreateRoomResponse構造体: 部屋作成レスポンスのデータ
type CreateRoomResponse struct {
	Room *entity.Room `json:"room"` // 作成した部屋
}

// CreateRoom: 新しい部屋を作成
func (r *RoomUseCase) CreateRoom(ctx context.Context, req CreateRoomRequest) (CreateRoomResponse, error) {
	id, err := r.roomIDFactory.NewRoomID()
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:      id,
		Name:    req.Name,
		Members: []entity.UserID{},
	})
	savedRoomID, err := r.roomRepo.SaveRoom(ctx, room)
	if err != nil {
		return CreateRoomResponse{nil}, err
	}
	res, err := r.roomRepo.GetRoomByID(ctx, savedRoomID)
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	return CreateRoomResponse{Room: res}, nil
}
