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

// 1.正常系のテスト
// 2.AddMemberToRoom（Repo）のエラー
func TestJoinRoom(t *testing.T) {
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
		roomID := entity.RoomID("room_1")
		userID := entity.UserID("test_user")

		mockRoomRepo.EXPECT().AddMemberToRoom(context.Background(), roomID, userID).Return(nil)

		err := roomUseCase.JoinRoom(context.Background(), roomUC.JoinRoomRequest{
			RoomID: roomID,
			UserID: userID,
		})

		assert.NoError(t, err)
	})

	t.Run("AddMemberToRoomのエラー", func(t *testing.T) {
		roomID := entity.RoomID("room_1")
		userID := entity.UserID("test_user")
		expectedErr := assert.AnError

		mockRoomRepo.EXPECT().AddMemberToRoom(context.Background(), roomID, userID).Return(expectedErr)

		err := roomUseCase.JoinRoom(context.Background(), roomUC.JoinRoomRequest{
			RoomID: roomID,
			UserID: userID,
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// 1.正常系のテスト
// 2.RemoveMemberFromRoom（Repo）のエラー
func TestLeaveRoom(t *testing.T) {
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
		roomID := entity.RoomID("room_1")
		userID := entity.UserID("test_user")

		mockRoomRepo.EXPECT().RemoveMemberFromRoom(context.Background(), roomID, userID).Return(nil)

		err := roomUseCase.LeaveRoom(context.Background(), roomUC.LeaveRoomRequest{
			RoomID: roomID,
			UserID: userID,
		})

		assert.NoError(t, err)
	})

	t.Run("RemoveMemberFromRoomのエラー", func(t *testing.T) {
		roomID := entity.RoomID("room_1")
		userID := entity.UserID("test_user")
		expectedErr := assert.AnError

		mockRoomRepo.EXPECT().RemoveMemberFromRoom(context.Background(), roomID, userID).Return(expectedErr)

		err := roomUseCase.LeaveRoom(context.Background(), roomUC.LeaveRoomRequest{
			RoomID: roomID,
			UserID: userID,
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
