package messagecase

import (
	"errors"

	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
)

type NewMessageUseCaseParams struct {
	MsgRepo  repository.MessageRepository
	MsgCache service.MessageCacheService
	RoomRepo repository.RoomRepository
	UserRepo repository.UserRepository
}

func (p *NewMessageUseCaseParams) Validate() error {
	if p.MsgRepo == nil {
		return errors.New("MsgRepo is required")
	}
	if p.MsgCache == nil {
		return errors.New("MsgCache is required")
	}
	if p.RoomRepo == nil {
		return errors.New("RoomRepo is required")
	}
	if p.UserRepo == nil {
		return errors.New("UserRepo is required")
	}
	return nil
}

func NewMessageUseCase(params NewMessageUseCaseParams) *MessageUseCase {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &MessageUseCase{
		msgRepo:  params.MsgRepo,
		msgCache: params.MsgCache,
		roomRepo: params.RoomRepo,
		userRepo: params.UserRepo,
	}
}
