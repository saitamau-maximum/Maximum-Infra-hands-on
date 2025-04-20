package usecase_test

import (
	"errors"
	"testing"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/usecase"
	mock_repository "example.com/webrtc-practice/mocks/domain/repository"
	mock_factory "example.com/webrtc-practice/mocks/interface/factory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	mockRoomPublicIDFactory := mock_factory.NewMockRoomPublicIDFactory(ctrl)

	params := usecase.NewRoomUseCaseParams{
		RoomRepo:            mockRoomRepo,
		RoomIDFactory:       mockRoomIDFactory,
		RoomPublicIDFactory: mockRoomPublicIDFactory,
	}

	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		req := usecase.CreateRoomRequest{
			Name:        "Test Room",
		}
		roomID := entity.RoomID(1)
		roomPublicID := entity.RoomPublicID("public_room_1")

		mockRoomIDFactory.EXPECT().NewRoomID().Return(roomID, nil)
		mockRoomPublicIDFactory.EXPECT().NewRoomPublicID().Return(roomPublicID, nil)
		mockRoomRepo.EXPECT().SaveRoom(gomock.Any()).Return(roomID, nil)

		resp, err := roomUseCase.CreateRoom(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp.RoomID)
		assert.Equal(t, roomID, *resp.RoomID)
	})

	t.Run("RoomID生成失敗時", func(t *testing.T) {
		req := usecase.CreateRoomRequest{
			Name:        "Test Room",
		}
		expectedErr := errors.New("failed to generate room ID")

		mockRoomIDFactory.EXPECT().NewRoomID().Return(entity.RoomID(1), expectedErr)

		resp, err := roomUseCase.CreateRoom(req)

		assert.Error(t, err)
		assert.Nil(t, resp.RoomID)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("RoomPublicID生成失敗時", func(t *testing.T) {
		req := usecase.CreateRoomRequest{
			Name:        "Test Room",
		}
		roomID := entity.RoomID(1)
		expectedErr := errors.New("failed to generate room public ID")

		mockRoomIDFactory.EXPECT().NewRoomID().Return(roomID, nil)
		mockRoomPublicIDFactory.EXPECT().NewRoomPublicID().Return(entity.RoomPublicID(""), expectedErr)

		resp, err := roomUseCase.CreateRoom(req)

		assert.Error(t, err)
		assert.Nil(t, resp.RoomID)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetRoomByPublicID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	mockRoomPublicIDFactory := mock_factory.NewMockRoomPublicIDFactory(ctrl)

	params := usecase.NewRoomUseCaseParams{
		RoomRepo:            mockRoomRepo,
		RoomIDFactory:       mockRoomIDFactory,
		RoomPublicIDFactory: mockRoomPublicIDFactory,
	}

	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		publicID := entity.RoomPublicID("public_room_1")
		roomID := entity.RoomID(1)
		room := entity.NewRoom(entity.RoomParams{
			ID:       roomID,
			PublicID: publicID,
			Name:     "Test Room",
			Members:  []entity.UserID{"user_1", "user_2"},
		})

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(roomID, nil)
		mockRoomRepo.EXPECT().GetRoomByID(roomID).Return(room, nil)

		resp, err := roomUseCase.GetRoomByPublicID(usecase.GetRoomByPublicIDParams{PublicID: publicID})

		assert.NoError(t, err)
		assert.NotNil(t, resp.Room)
		assert.Equal(t, room, resp.Room)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomPublicID("nonexistent_room")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(entity.RoomID(0), nil)

		resp, err := roomUseCase.GetRoomByPublicID(usecase.GetRoomByPublicIDParams{PublicID: publicID})

		assert.Error(t, err)
		assert.Nil(t, resp.Room)
		assert.Equal(t, "room not found", err.Error())
	})
}

func TestGetAllRooms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo: mockRoomRepo,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		rooms := []*entity.Room{
			entity.NewRoom(entity.RoomParams{ID: 1, Name: "Room1"}),
			entity.NewRoom(entity.RoomParams{ID: 2, Name: "Room2"}),
		}

		mockRoomRepo.EXPECT().GetAllRooms().Return(rooms, nil)

		result, err := roomUseCase.GetAllRooms()

		assert.NoError(t, err)
		assert.Equal(t, rooms, result)
	})

	t.Run("エラー発生時", func(t *testing.T) {
		expectedErr := errors.New("failed to fetch rooms")

		mockRoomRepo.EXPECT().GetAllRooms().Return(nil, expectedErr)

		result, err := roomUseCase.GetAllRooms()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetUsersInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo: mockRoomRepo,
	}
	roomUseCase := usecase.NewRoomUseCase(params)
	// == Data ==
	publicID := entity.RoomPublicID("public_room_1")
	roomID := entity.RoomID(1)
	users := []*entity.User{
		entity.NewUser(entity.UserParams{
			ID:         "user_1",
			Name:       "User 1",
			Email:      "test@test",
			PasswdHash: "hash",
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}),
		entity.NewUser(entity.UserParams{
			ID:         "user_2",
			Name:       "User 2",
			Email:      "test@test",
			PasswdHash: "hash",
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}),
	}

	t.Run("正常系", func(t *testing.T) {

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(roomID, nil)
		mockRoomRepo.EXPECT().GetUsersInRoom(roomID).Return(users, nil)

		resp, err := roomUseCase.GetUsersInRoom(usecase.GetUsersInRoomRequest{PublicID: publicID})

		assert.NoError(t, err)
		assert.Equal(t, users, resp.Users)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomPublicID("nonexistent_room")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(entity.RoomID(0), nil)

		resp, err := roomUseCase.GetUsersInRoom(usecase.GetUsersInRoomRequest{PublicID: publicID})

		assert.Error(t, err)
		assert.Nil(t, resp.Users)
		assert.Equal(t, "room not found", err.Error())
	})
}

func TestJoinRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo: mockRoomRepo,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		publicID := entity.RoomPublicID("public_room_1")
		roomID := entity.RoomID(1)
		user := entity.NewUser(entity.UserParams{
			ID:         "user_1",
			Name:       "User 1",
			Email:      "test@test",
			PasswdHash: "hash",
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		})

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(roomID, nil)
		mockRoomRepo.EXPECT().AddMemberToRoom(roomID, user.GetID()).Return(nil)

		err := roomUseCase.JoinRoom(usecase.JoinRoomRequest{RoomPublicID: publicID, User: user})

		assert.NoError(t, err)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomPublicID("nonexistent_room")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(entity.RoomID(0), nil)

		err := roomUseCase.JoinRoom(usecase.JoinRoomRequest{RoomPublicID: publicID, User: entity.NewUser(
			entity.UserParams{
				ID:         "user_1",
				Name:       "User 1",
				Email:      "test@test",
				PasswdHash: "hash",
				CreatedAt:  time.Now(),
				UpdatedAt:  nil,
			},
		)})

		assert.Error(t, err)
		assert.Equal(t, "room not found", err.Error())
	})
}

func TestLeaveRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo: mockRoomRepo,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		publicID := entity.RoomPublicID("public_room_1")
		roomID := entity.RoomID(1)
		userID := entity.UserID("user_1")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(roomID, nil)
		mockRoomRepo.EXPECT().RemoveMemberFromRoom(roomID, userID).Return(nil)

		err := roomUseCase.LeaveRoom(usecase.LeaveRoomRequest{RoomPublicID: publicID, UserID: userID})

		assert.NoError(t, err)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomPublicID("nonexistent_room")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(entity.RoomID(0), nil)

		err := roomUseCase.LeaveRoom(usecase.LeaveRoomRequest{RoomPublicID: publicID, UserID: "user_1"})

		assert.Error(t, err)
		assert.Equal(t, "room not found", err.Error())
	})
}

func TestUpdateRoomName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo: mockRoomRepo,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		publicID := entity.RoomPublicID("public_room_1")
		roomID := entity.RoomID(1)
		newName := "Updated Room Name"

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(roomID, nil)
		mockRoomRepo.EXPECT().UpdateRoomName(roomID, newName).Return(nil)

		err := roomUseCase.UpdateRoomName(usecase.UpdateRoomNameRequest{RoomPublicID: publicID, NewName: newName})

		assert.NoError(t, err)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomPublicID("nonexistent_room")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(entity.RoomID(0), nil)

		err := roomUseCase.UpdateRoomName(usecase.UpdateRoomNameRequest{RoomPublicID: publicID, NewName: "New Name"})

		assert.Error(t, err)
		assert.Equal(t, "room not found", err.Error())
	})
}

func TestDeleteRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo: mockRoomRepo,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		publicID := entity.RoomPublicID("public_room_1")
		roomID := entity.RoomID(1)

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(roomID, nil)
		mockRoomRepo.EXPECT().DeleteRoom(roomID).Return(nil)

		err := roomUseCase.DeleteRoom(usecase.DeleteRoomRequest{RoomPublicID: publicID})

		assert.NoError(t, err)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomPublicID("nonexistent_room")

		mockRoomRepo.EXPECT().GetRoomIDByPublicID(publicID).Return(entity.RoomID(0), nil)

		err := roomUseCase.DeleteRoom(usecase.DeleteRoomRequest{RoomPublicID: publicID})

		assert.Error(t, err)
		assert.Equal(t, "room not found", err.Error())
	})
}
