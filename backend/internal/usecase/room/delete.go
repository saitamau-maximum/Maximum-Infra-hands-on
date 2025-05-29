package room

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// DeleteRoomRequest構造体: 部屋を削除するリクエスト
type DeleteRoomRequest struct {
	RoomID entity.RoomID `json:"room_id"` // 部屋の公開ID
}

// DeleteRoom: 部屋を削除
func (r *RoomUseCase) DeleteRoom(ctx context.Context, req DeleteRoomRequest) error {
	err := r.roomRepo.DeleteRoom(ctx, req.RoomID)
	if err != nil {
		return err
	}
	return nil
}
