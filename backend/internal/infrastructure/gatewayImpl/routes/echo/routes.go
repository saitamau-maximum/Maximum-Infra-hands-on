package routes

import (
	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/interface/handler/messagehandler"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"example.com/infrahandson/internal/interface/handler/websockethandler"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	cfg *config.Config,
	AuthMiddleware echo.MiddlewareFunc,
	userHandler userhandler.UserHandlerInterface,
	roomHandler roomhandler.RoomHandlerInterface,
	wsHandler websockethandler.WebSocketHandlerInterface,
	msgHansler messagehandler.MessageHandlerInterface,
) {
	userGroup := e.Group("/api/user")
	RegisterUserRoutes(userGroup, userHandler, AuthMiddleware)
	roomGroup := e.Group("/api/room", AuthMiddleware)
	RegisterRoomRoutes(roomGroup, roomHandler)
	wsGroup := e.Group("/api/ws", AuthMiddleware)
	RegisterWsRoutes(wsGroup, wsHandler)
	msgGroup := e.Group("/api/message", AuthMiddleware)
	RegisterMsgRoutes(msgGroup, msgHansler)
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

func RegisterRoomRoutes(g *echo.Group, h roomhandler.RoomHandlerInterface) {
	g.POST("", h.CreateRoom)
	g.POST("/:room_id/join", h.JoinRoom)
	g.POST("/:room_id/leave", h.LeaveRoom)
	g.GET("/:room_id", h.GetRoomByID)
	g.GET("", h.GetRooms)
}

func RegisterWsRoutes(g *echo.Group, h websockethandler.WebSocketHandlerInterface) {
	g.GET("/:room_id", h.ConnectToChatRoom)
}

func RegisterMsgRoutes(g *echo.Group, h messagehandler.MessageHandlerInterface) {
	g.GET("/:room_id", h.GetRoomMessage)
}
