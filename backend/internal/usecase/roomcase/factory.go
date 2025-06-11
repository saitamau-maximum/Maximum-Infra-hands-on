package roomcase

import (
	"errors"

	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/interface/factory"
)

type NewRoomUseCaseParams struct {
	RoomRepo      repository.RoomRepository
	UserRepo      repository.UserRepository
	RoomIDFactory factory.RoomIDFactory
}

func (p NewRoomUseCaseParams) Validate() error {
	if p.RoomRepo == nil {
		return errors.New("RoomRepo is required")
	}
	if p.UserRepo == nil {
		return errors.New("UserRepo is required")
	}
	if p.RoomIDFactory == nil {
		return errors.New("RoomIDFactory is required")
	}
	return nil
}

// NewRoomUseCase: RoomUseCaseのインスタンスを生成
func NewRoomUseCase(p NewRoomUseCaseParams) *RoomUseCase {
	if err := p.Validate(); err != nil {
		panic(err)
	}

	return &RoomUseCase{
		roomRepo:      p.RoomRepo,
		userRepo:      p.UserRepo,
		roomIDFactory: p.RoomIDFactory,
	}
}
