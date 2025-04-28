package handler

import (
	"net/http"

	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	WsUseCase        usecase.WebsocketUseCaseInterface
	WsUpgrader       adapter.WebSocketUpgraderAdapter
	WsConnFactory    factory.WebSocketConnectionFactory
	UserIDFactory    factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger           adapter.LoggerAdapter
}

type NewWebSocketHandlerParams struct {
	WsUseCase        usecase.WebsocketUseCaseInterface
	WsUpgrader       adapter.WebSocketUpgraderAdapter
	WsConnFactory    factory.WebSocketConnectionFactory
	UserIDFactory    factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger           adapter.LoggerAdapter
}

func (p *NewWebSocketHandlerParams) Validate() error {
	if p.WsUseCase == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "WebsocketUseCase is required")
	}
	if p.WsUpgrader == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "WebsocketUpgrader is required")
	}
	if p.WsConnFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "WsConnFactory is required")
	}
	if p.UserIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserIDFactory is required")
	}
	if p.RoomIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "RoomIDFactory is required")
	}
	if p.Logger == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Logger is required")
	}
	return nil
}

func NewWebSocketHandler(params NewWebSocketHandlerParams) *WebSocketHandler {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &WebSocketHandler{
		WsUseCase:        params.WsUseCase,
		WsUpgrader:       params.WsUpgrader,
		WsConnFactory:    params.WsConnFactory,
		UserIDFactory:    params.UserIDFactory,
		RoomIDFactory: params.RoomIDFactory,
		Logger:           params.Logger,
	}
}

func (h *WebSocketHandler) ConnectToChatRoom(c echo.Context) error {
	h.Logger.Info("ConnectToChatRoom called")

	userPublicIDStr, ok := c.Get("user_id").(string)
	if !ok || userPublicIDStr == "" {
		h.Logger.Warn("User ID is missing or invalid")
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID is required")
	}
	userPublicID, err := h.UserIDFactory.FromString(userPublicIDStr)
	if err != nil {
		h.Logger.Warn("Invalid User ID", "user_id", userPublicIDStr)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid User ID")
	}

	roomPublicIDStr := c.Param("room_public_id")
	if roomPublicIDStr == "" {
		h.Logger.Warn("Room public ID is missing")
		return echo.NewHTTPError(http.StatusBadRequest, "Room public ID is required")
	}
	roomPublicID, err := h.RoomIDFactory.FromString(roomPublicIDStr)
	if err != nil {
		h.Logger.Warn("Invalid Room public ID", "room_public_id", roomPublicIDStr)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Room public ID")
	}

	h.Logger.Info("Upgrading WebSocket connection", "room_public_id", roomPublicIDStr, "user_id", userPublicIDStr)
	connRaw, err := h.WsUpgrader.Upgrade(c.Response().Writer, c.Request())
	if err != nil {
		h.Logger.Error("Failed to upgrade connection", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upgrade connection")
	}

	conn, err := h.WsConnFactory.CreateWebSocketConnection(connRaw)
	if err != nil {
		h.Logger.Error("Failed to create WebSocket connection", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create WebSocket connection")
	}

	if err := h.WsUseCase.ConnectUserToRoom(usecase.ConnectUserToRoomRequest{
		UserPublicID: userPublicID,
		RoomPublicID: roomPublicID,
		Conn:         conn,
	}); err != nil {
		h.Logger.Error("Failed to connect user to room", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect user to room")
	}

	h.Logger.Info("User connected to room", "room_public_id", roomPublicIDStr, "user_id", userPublicIDStr)
	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				h.Logger.Warn("Connection closed or error reading message", "error", err)
				_ = h.WsUseCase.DisconnectUser(usecase.DisconnectUserRequest{
					UserID: userPublicID,
				})
				return
			}

			h.Logger.Info("Message received", "room_public_id", roomPublicIDStr, "user_id", userPublicIDStr)
			h.WsUseCase.SendMessage(usecase.SendMessageRequest{
				RoomPublicID: roomPublicID,
				Sender:       userPublicID,
				Content:      message.GetContent(),
			})
		}
	}()

	return nil
}
