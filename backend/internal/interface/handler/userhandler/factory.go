package userhandler

import (
	"errors"

	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase/usercase"
)

type NewUserHandlerParams struct {
	UserUseCase   usercase.UserUseCaseInterface
	UserIDFactory factory.UserIDFactory
	Logger        adapter.LoggerAdapter
}

func (p *NewUserHandlerParams) Validate() error {
	if p.UserUseCase == nil {
		return errors.New("userUseCase is required")
	}
	if p.UserIDFactory == nil {
		return errors.New("userIDFactory is required")
	}
	if p.Logger == nil {
		return errors.New("logger is required")
	}
	return nil
}

func NewUserHandler(params NewUserHandlerParams) UserHandlerInterface {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &UserHandler{
		UserUseCase:   params.UserUseCase,
		UserIDFactory: params.UserIDFactory,
		Logger:        params.Logger,
	}
}
