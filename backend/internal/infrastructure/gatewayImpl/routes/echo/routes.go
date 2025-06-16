package routes

import (
	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	cfg *config.Config,
	AuthMiddleware echo.MiddlewareFunc,
	userHandler userhandler.UserHandlerInterface,
	roomHandler handler.RoomHandler,
	wsHandler handler.WebSocketHandler,
	msgHansler handler.MessageHandler,
) {
	userGroup := e.Group("/api/user")
	RegisterUserRoutes(userGroup, userHandler, AuthMiddleware)
	roomGroup := e.Group("/api/room", AuthMiddleware)
	RegisterRoomRoutes(roomGroup, &roomHandler)
	wsGroup := e.Group("/api/ws", AuthMiddleware)
	RegisterWsRoutes(wsGroup, &wsHandler)
	msgGroup := e.Group("/api/message", AuthMiddleware)
	RegisterMsgRoutes(msgGroup, &msgHansler)
}

// RegisterUserRoutes はユーザー関連のルートを登録する
func RegisterUserRoutes(g *echo.Group, h userhandler.UserHandlerInterface, authMiddleware echo.MiddlewareFunc) {
	g.POST("/register", h.RegisterUser)
	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout, authMiddleware)
	g.POST("/icon", h.SaveUserIcon, authMiddleware)
	g.GET("/me", h.GetMe, authMiddleware)
	g.GET("/icon/:user_id", h.GetUserIcon)
}

func RegisterRoomRoutes(g *echo.Group, h *handler.RoomHandler) {
	g.POST("", h.CreateRoom)
	g.POST("/:room_id/join", h.JoinRoom)
	g.POST("/:room_id/leave", h.LeaveRoom)
	g.GET("/:room_id", h.GetRoom)
	g.GET("", h.GetRooms)
}

func RegisterWsRoutes(g *echo.Group, h *handler.WebSocketHandler) {
	g.GET("/:room_id", h.ConnectToChatRoom)
}

func RegisterMsgRoutes(g *echo.Group, h *handler.MessageHandler) {
	g.GET("/:room_id", h.GetRoomMessage)
}
