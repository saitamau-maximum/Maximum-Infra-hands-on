package routes

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/interface/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	cfg *config.Config,
	userHandler handler.UserHandler,
	websocketHandler handler.WebSocketHandler,
	messageHandler handler.MessageHandler,
	roomHandler handler.RoomHandler,
) {
	userGroup := e.Group("/api/user")
	userHandler.Register(userGroup)
	// Websocket
	websocketGroup := e.Group("/ws")
	websocketHandler.Register(websocketGroup)
	// Message
	messageGroup := e.Group("/api/message")
	messageHandler.Register(messageGroup)
	// Room
	roomGroup := e.Group("/api/room")
	roomHandler.Register(roomGroup)
}
