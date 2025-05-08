package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}

	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		req := usecase.CreateRoomRequest{
			Name: "Test Room",
		}
		roomID := entity.RoomID("public_room_1")

		mockRoomIDFactory.EXPECT().NewRoomID().Return(roomID, nil)
		mockRoomRepo.EXPECT().SaveRoom(context.Background(), gomock.Any()).Return(roomID, nil)
		mockRoomRepo.EXPECT().GetRoomByID(context.Background(), roomID).Return(entity.NewRoom(entity.RoomParams{
			ID:      roomID,
			Name:    req.Name,
			Members: []entity.UserID{},
		}), nil)

		resp, err := roomUseCase.CreateRoom(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp.Room)
	})

	t.Run("RoomID生成失敗時", func(t *testing.T) {
		req := usecase.CreateRoomRequest{
			Name: "Test Room",
		}
		expectedErr := errors.New("failed to generate room ID")
		roomID := entity.RoomID("public_room_1")

		mockRoomIDFactory.EXPECT().NewRoomID().Return(roomID, expectedErr)

		resp, err := roomUseCase.CreateRoom(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp.Room)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("RoomID生成失敗時", func(t *testing.T) {
		req := usecase.CreateRoomRequest{
			Name: "Test Room",
		}
		publicID := entity.RoomID("public_room_1")
		expectedErr := errors.New("failed to generate room public ID")

		mockRoomIDFactory.EXPECT().NewRoomID().Return(publicID, expectedErr)

		resp, err := roomUseCase.CreateRoom(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp.Room)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetRoomByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}

	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")
		room := entity.NewRoom(entity.RoomParams{
			ID:      roomID,
			Name:    "Test Room",
			Members: []entity.UserID{"user_1", "user_2"},
		})

		mockRoomRepo.EXPECT().GetRoomByID(context.Background(), roomID).Return(room, nil)

		resp, err := roomUseCase.GetRoomByID(context.Background() ,usecase.GetRoomByIDRequest{ID: roomID})

		assert.NoError(t, err)
		assert.NotNil(t, resp.Room)
		assert.Equal(t, room, resp.Room)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomID("nonexistent_room")

		mockRoomRepo.EXPECT().GetRoomByID(context.Background(), publicID).Return(nil, errors.New("room not found"))
		resp, err := roomUseCase.GetRoomByID(context.Background(), usecase.GetRoomByIDRequest{ID: publicID})

		assert.Error(t, err)
		assert.Nil(t, resp.Room)
		assert.Equal(t, "room not found", err.Error())
	})
}

func TestGetAllRooms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		rooms := []*entity.Room{
			entity.NewRoom(entity.RoomParams{ID: "room1", Name: "Room1"}),
			entity.NewRoom(entity.RoomParams{ID: "room2", Name: "Room2"}),
		}

		mockRoomRepo.EXPECT().GetAllRooms(context.Background()).Return(rooms, nil)

		result, err := roomUseCase.GetAllRooms(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, rooms, result)
	})

	t.Run("エラー発生時", func(t *testing.T) {
		expectedErr := errors.New("failed to fetch rooms")

		mockRoomRepo.EXPECT().GetAllRooms(context.Background()).Return(nil, expectedErr)

		result, err := roomUseCase.GetAllRooms(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetUsersInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := usecase.NewRoomUseCase(params)
	// == Data ==
	roomID := entity.RoomID("public_room_1")
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

		mockRoomRepo.EXPECT().GetUsersInRoom(context.Background(), roomID).Return(users, nil)

		resp, err := roomUseCase.GetUsersInRoom(context.Background(), usecase.GetUsersInRoomRequest{ID: roomID})

		assert.NoError(t, err)
		assert.Equal(t, users, resp.Users)
	})

	t.Run("部屋が見つからない場合", func(t *testing.T) {
		publicID := entity.RoomID("nonexistent_room")

		mockRoomRepo.EXPECT().GetUsersInRoom(context.Background(), publicID).Return(nil, errors.New("room not found"))
		resp, err := roomUseCase.GetUsersInRoom(context.Background(), usecase.GetUsersInRoomRequest{ID: publicID})

		assert.Error(t, err)
		assert.Nil(t, resp.Users)
		assert.Equal(t, "room not found", err.Error())
	})
}

func TestJoinRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("room_1")
		userID := entity.UserID("test_user")

		mockRoomRepo.EXPECT().AddMemberToRoom(context.Background(), roomID, userID).Return(nil)

		err := roomUseCase.JoinRoom(context.Background(), usecase.JoinRoomRequest{
			RoomID: roomID,
			UserID: userID,
		})

		assert.NoError(t, err)
	})
}

func TestLeaveRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("room_1")
		userID := entity.UserID("test_user")

		mockRoomRepo.EXPECT().RemoveMemberFromRoom(context.Background(), roomID, userID).Return(nil)

		err := roomUseCase.LeaveRoom(context.Background(), usecase.LeaveRoomRequest{
			RoomID: roomID,
			UserID: userID,
		})

		assert.NoError(t, err)
	})
}

func TestUpdateRoomName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")
		newName := "Updated Room Name"

		mockRoomRepo.EXPECT().UpdateRoomName(context.Background(), roomID, newName).Return(nil)

		err := roomUseCase.UpdateRoomName(context.Background(), usecase.UpdateRoomNameRequest{RoomID: roomID, NewName: newName})

		assert.NoError(t, err)
	})
}

func TestDeleteRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := usecase.NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	roomUseCase := usecase.NewRoomUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")

		mockRoomRepo.EXPECT().DeleteRoom(context.Background(), roomID).Return(nil)

		err := roomUseCase.DeleteRoom(context.Background(), usecase.DeleteRoomRequest{RoomID: roomID})

		assert.NoError(t, err)
	})
}
