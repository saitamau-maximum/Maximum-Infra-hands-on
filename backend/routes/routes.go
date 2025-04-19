package routes

import (
	"example.com/webrtc-practice/config"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo, 
	cfg *config.Config, 
	// userHandler handler.UserHandler,
	// websocketHandler handler.WebsocketHandler,
) {

	// userGroup := e.Group("/api/user")
	// userHandler.Register(userGroup)

	// // Websocket
	// websocketGroup := e.Group("/ws")
	// websocketHandler.Register(websocketGroup)
}
