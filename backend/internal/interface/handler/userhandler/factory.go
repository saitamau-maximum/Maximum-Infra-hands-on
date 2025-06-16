package userhandler

import (
	"net/http"

	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase/usercase"
	"github.com/labstack/echo/v4"
)

type NewUserHandlerParams struct {
	UserUseCase   usercase.UserUseCaseInterface
	UserIDFactory factory.UserIDFactory
	Logger        adapter.LoggerAdapter
}

func (p *NewUserHandlerParams) Validate() error {
	if p.UserUseCase == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserUseCase is required")
	}
	if p.UserIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserIDFactory is required")
	}
	if p.Logger == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Logger is required")
	}
	return nil
}

func NewUserHandler(params NewUserHandlerParams) *UserHandler {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &UserHandler{
		UserUseCase:   params.UserUseCase,
		UserIDFactory: params.UserIDFactory,
		Logger:        params.Logger,
	}
}
