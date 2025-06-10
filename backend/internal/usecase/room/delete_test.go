package room_test

import (
	"context"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	roomUC "example.com/infrahandson/internal/usecase/room"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. DeleteRoom失敗

func TestDeleteRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomUseCase, mockDeps := roomUC.NewTestRoomUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")

		mockDeps.RoomRepo.EXPECT().DeleteRoom(context.Background(), roomID).Return(nil)

		err := roomUseCase.DeleteRoom(context.Background(), roomUC.DeleteRoomRequest{RoomID: roomID})

		assert.NoError(t, err)
	})

	t.Run("DeleteRoom失敗", func(t *testing.T) {
		roomID := entity.RoomID("nonexistent_room")

		mockDeps.RoomRepo.EXPECT().DeleteRoom(context.Background(), roomID).Return(assert.AnError)

		err := roomUseCase.DeleteRoom(context.Background(), roomUC.DeleteRoomRequest{RoomID: roomID})

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}
