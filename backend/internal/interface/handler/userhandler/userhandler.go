package userhandler

import (
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase/usercase"
)

type UserHandler struct {
	UserUseCase   usercase.UserUseCaseInterface
	UserIDFactory factory.UserIDFactory
	Logger        adapter.LoggerAdapter
}
