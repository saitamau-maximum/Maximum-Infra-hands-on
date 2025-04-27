package routes

import (
	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/interface/handler"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	cfg *config.Config,
	AuthMiddleware echo.MiddlewareFunc,
	userHandler handler.UserHandler,
	roomHandler handler.RoomHandler,
) {
	userGroup := e.Group("/api/user")
	RegisterUserRoutes(userGroup, &userHandler, AuthMiddleware)
	roomGroup := e.Group("/api/room")
	roomHandler.Register(roomGroup)
}

// routes/user_routes.go
func RegisterUserRoutes(g *echo.Group, h *handler.UserHandler, authMiddleware echo.MiddlewareFunc) {
	g.POST("/register", h.RegisterUser)
	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout, authMiddleware)
	g.GET("/me", h.GetMe, authMiddleware)
}
