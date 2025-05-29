package room

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// UpdateRoomNameRequest構造体: 部屋名を更新するリクエスト
type UpdateRoomNameRequest struct {
	RoomID  entity.RoomID `json:"room_id"`
	NewName string        `json:"new_name"` // 新しい部屋名
}

// UpdateRoomName: 部屋名を更新
func (r *RoomUseCase) UpdateRoomName(ctx context.Context, req UpdateRoomNameRequest) error {
	err := r.roomRepo.UpdateRoomName(ctx, req.RoomID, req.NewName)
	if err != nil {
		return err
	}
	return nil
}
