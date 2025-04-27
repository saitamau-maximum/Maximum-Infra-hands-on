package routes

import (
	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/interface/handler"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	cfg *config.Config,
	userHandler handler.UserHandler,
	roomHandler handler.RoomHandler,
) {
	userGroup := e.Group("/api/user")
	userHandler.Register(userGroup)
	roomGroup := e.Group("/api/room")
	roomHandler.Register(roomGroup)
}
