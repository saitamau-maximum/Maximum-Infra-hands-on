package roomcase

import (
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	"go.uber.org/mock/gomock"
)

type mockDeps struct {
	RoomRepo      *mock_repository.MockRoomRepository
	UserRepo      *mock_repository.MockUserRepository
	RoomIDFactory *mock_factory.MockRoomIDFactory
}

func NewTestRoomUseCase(
	ctrl *gomock.Controller,
) (RoomUseCaseInterface, mockDeps) {
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockRoomRepo := mock_repository.NewMockRoomRepository(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	params := NewRoomUseCaseParams{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
	useCase := NewRoomUseCase(params)

	return useCase, mockDeps{
		RoomRepo:      mockRoomRepo,
		UserRepo:      mockUserRepo,
		RoomIDFactory: mockRoomIDFactory,
	}
}
