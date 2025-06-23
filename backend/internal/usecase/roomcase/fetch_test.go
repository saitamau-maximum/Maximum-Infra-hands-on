package roomcase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/roomcase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. GetRoomByID（Repo）エラー
// 3. GetRoomByID（Repo）の返答がnil

func TestGetRoomByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomUseCase, mockDeps := roomcase.NewTestRoomUseCase(ctrl)

	t.Run("1.正常系", func(t *testing.T) {
		roomID := entity.RoomID("public_room_1")
		room := entity.NewRoom(entity.RoomParams{
			ID:      roomID,
			Name:    "Test Room",
			Members: []entity.UserID{"user_1", "user_2"},
		})

		mockDeps.RoomRepo.EXPECT().GetRoomByID(context.Background(), roomID).Return(room, nil)

		resp, err := roomUseCase.GetRoomByID(context.Background(), roomcase.GetRoomByIDRequest{ID: roomID})

		assert.NoError(t, err)
		assert.NotNil(t, resp.Room)
		assert.Equal(t, room, resp.Room)
	})

	t.Run("2.GetRoomByID（Repo）エラー", func(t *testing.T) {
		publicID := entity.RoomID("nonexistent_room")

		mockDeps.RoomRepo.EXPECT().GetRoomByID(context.Background(), publicID).Return(nil, assert.AnError)
		resp, err := roomUseCase.GetRoomByID(context.Background(), roomcase.GetRoomByIDRequest{ID: publicID})

		assert.Error(t, err)
		assert.Nil(t, resp.Room)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("3.GetRoomByID（Repo）の返答がnil", func(t *testing.T) {
		publicID := entity.RoomID("nonexistent_room")

		mockDeps.RoomRepo.EXPECT().GetRoomByID(context.Background(), publicID).Return(nil, nil)
		resp, err := roomUseCase.GetRoomByID(context.Background(), roomcase.GetRoomByIDRequest{ID: publicID})

		assert.Error(t, err)
		assert.Nil(t, resp.Room)
		assert.Equal(t, "room not found", err.Error())
	})
}

// 1. 正常系
// 2. GetAllRooms（Repo）エラー

func TestGetAllRooms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomUseCase, mockDeps := roomcase.NewTestRoomUseCase(ctrl)

	t.Run("1.正常系", func(t *testing.T) {
		rooms := []*entity.Room{
			entity.NewRoom(entity.RoomParams{ID: "room1", Name: "Room1"}),
			entity.NewRoom(entity.RoomParams{ID: "room2", Name: "Room2"}),
		}

		mockDeps.RoomRepo.EXPECT().GetAllRooms(context.Background()).Return(rooms, nil)

		result, err := roomUseCase.GetAllRooms(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, rooms, result)
	})

	t.Run("2.GetAllRooms（Repo）エラー", func(t *testing.T) {
		expectedErr := errors.New("failed to fetch rooms")

		mockDeps.RoomRepo.EXPECT().GetAllRooms(context.Background()).Return(nil, expectedErr)

		result, err := roomUseCase.GetAllRooms(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

// 1. 正常系
// 2. GetUsersInRoom（Repo）エラー

func TestGetUsersInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roomUseCase, mockDeps := roomcase.NewTestRoomUseCase(ctrl)
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

		mockDeps.RoomRepo.EXPECT().GetUsersInRoom(context.Background(), roomID).Return(users, nil)

		resp, err := roomUseCase.GetUsersInRoom(context.Background(), roomcase.GetUsersInRoomRequest{ID: roomID})

		assert.NoError(t, err)
		assert.Equal(t, users, resp.Users)
	})

	t.Run("2.GetUsersInRoom（Repo）エラー", func(t *testing.T) {
		publicID := entity.RoomID("nonexistent_room")

		mockDeps.RoomRepo.EXPECT().GetUsersInRoom(context.Background(), publicID).Return(nil, assert.AnError)
		resp, err := roomUseCase.GetUsersInRoom(context.Background(), roomcase.GetUsersInRoomRequest{ID: publicID})

		assert.Error(t, err)
		assert.Nil(t, resp.Users)
		assert.Equal(t, assert.AnError, err)
	})
}
