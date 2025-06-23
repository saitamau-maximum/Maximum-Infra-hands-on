package roomcase

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// JoinRoomRequest構造体: 部屋に参加するリクエスト
type JoinRoomRequest struct {
	RoomID entity.RoomID `json:"room_id"` // 部屋の公開ID
	UserID entity.UserID `json:"user_id"` // 参加するユーザー
}

// JoinRoom: 部屋にユーザーを参加させる
func (r *RoomUseCase) JoinRoom(ctx context.Context, req JoinRoomRequest) error {
	err := r.roomRepo.AddMemberToRoom(ctx, req.RoomID, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

// LeaveRoomRequest構造体: 部屋から退出するリクエスト
type LeaveRoomRequest struct {
	RoomID entity.RoomID `json:"room_id"` // 部屋の公開ID
	UserID entity.UserID `json:"user_id"` // 退出するユーザーID
}

// LeaveRoom: 部屋からユーザーを退出させる
func (r *RoomUseCase) LeaveRoom(ctx context.Context, req LeaveRoomRequest) error {
	err := r.roomRepo.RemoveMemberFromRoom(ctx, req.RoomID, req.UserID)
	if err != nil {
		return err
	}

	return nil
}
