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
	wsHandler handler.WebSocketHandler,
	msgHansler handler.MessageHandler,
) {
	userGroup := e.Group("/api/user")
	RegisterUserRoutes(userGroup, &userHandler, AuthMiddleware)
	roomGroup := e.Group("/api/room", AuthMiddleware)
	RegisterRoomRoutes(roomGroup, &roomHandler)
	wsGroup := e.Group("/api/ws", AuthMiddleware)
	RegisterWsRoutes(wsGroup, &wsHandler)
	msgGroup := e.Group("/api/message", AuthMiddleware)
	RegisterMsgRoutes(msgGroup, &msgHansler)
}

func RegisterUserRoutes(g *echo.Group, h *handler.UserHandler, authMiddleware echo.MiddlewareFunc) {
	g.POST("/register", h.RegisterUser)
	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout, authMiddleware)
	g.POST("/icon", h.SaveUserIcon, authMiddleware)
	g.GET("/me", h.GetMe, authMiddleware)
	g.GET("/icon/:user_id", h.GetUserIcon)
}

func RegisterRoomRoutes(g *echo.Group, h *handler.RoomHandler) {
	g.POST("", h.CreateRoom)
	g.POST("/:room_public_id/join", h.JoinRoom)
	g.POST("/:room_public_id/leave", h.LeaveRoom)
	g.GET("/:room_public_id", h.GetRoom)
	g.GET("", h.GetRooms)
}

func RegisterWsRoutes(g *echo.Group, h *handler.WebSocketHandler) {
	g.GET("/:room_id", h.ConnectToChatRoom)
}

func RegisterMsgRoutes(g *echo.Group, h *handler.MessageHandler) {
	g.GET("/:room_id", h.GetMessageHistoryInRoom)
}
