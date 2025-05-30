package room_test

import (
	"context"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	roomUC "example.com/infrahandson/internal/usecase/room"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. DeleteRoom失敗

func TestDeleteRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := roomUC.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := roomUC.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")

		mockRoomRepo.EXPECT().DeleteRoom(context.Background(), roomID).Return(nil)

		err := roomUseCase.DeleteRoom(context.Background(), roomUC.DeleteRoomRequest{RoomID: roomID})

		assert.NoError(t, err)
	})

	t.Run("DeleteRoom失敗", func(t *testing.T) {
		roomID := entity.RoomID("nonexistent_room")

		mockRoomRepo.EXPECT().DeleteRoom(context.Background(), roomID).Return(assert.AnError)

		err := roomUseCase.DeleteRoom(context.Background(), roomUC.DeleteRoomRequest{RoomID: roomID})

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}
