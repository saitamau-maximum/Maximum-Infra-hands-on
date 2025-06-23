package usercase

import (
	"errors"

	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
)

type NewUserUseCaseParams struct {
	UserRepo      repository.UserRepository
	Hasher        adapter.HasherAdapter
	TokenSvc      adapter.TokenServiceAdapter
	IconSvc       service.IconStoreService
	UserIDFactory factory.UserIDFactory
}

func (p NewUserUseCaseParams) Validate() error {
	if p.UserRepo == nil {
		return errors.New("userRepo is required")
	}
	if p.Hasher == nil {
		return errors.New("hasher is required")
	}
	if p.TokenSvc == nil {
		return errors.New("tokenSvc is required")
	}
	if p.IconSvc == nil {
		return errors.New("iconSvc is required")
	}
	if p.UserIDFactory == nil {
		return errors.New("userIDFactory is required")
	}
	return nil
}

func NewUserUseCase(p NewUserUseCaseParams) UserUseCaseInterface {
	if err := p.Validate(); err != nil {
		panic(err)
	}

	return &UserUseCase{
		userRepo:      p.UserRepo,
		hasher:        p.Hasher,
		tokenSvc:      p.TokenSvc,
		iconSvc:       p.IconSvc,
		userIDFactory: p.UserIDFactory,
	}
}
