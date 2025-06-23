package roomcase_test

import (
	"context"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/roomcase"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateRoomName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := roomcase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := roomcase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")
		newName := "Updated Room Name"

		mockRoomRepo.EXPECT().UpdateRoomName(context.Background(), roomID, newName).Return(nil)

		err := roomUseCase.UpdateRoomName(context.Background(), roomcase.UpdateRoomNameRequest{RoomID: roomID, NewName: newName})

		assert.NoError(t, err)
	})

	t.Run("UpdateRoomNameのエラー", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")
		newName := "Updated Room Name"
		expectedErr := assert.AnError

		mockRoomRepo.EXPECT().UpdateRoomName(context.Background(), roomID, newName).Return(expectedErr)

		err := roomUseCase.UpdateRoomName(context.Background(), roomcase.UpdateRoomNameRequest{RoomID: roomID, NewName: newName})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
