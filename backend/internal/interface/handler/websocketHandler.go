package handler

import (
	"context"
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	WsUseCase     usecase.WebsocketUseCaseInterface
	WsUpgrader    adapter.WebSocketUpgraderAdapter
	WsConnFactory factory.WebSocketConnectionFactory
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
}

type NewWebSocketHandlerParams struct {
	WsUseCase     usecase.WebsocketUseCaseInterface
	WsUpgrader    adapter.WebSocketUpgraderAdapter
	WsConnFactory factory.WebSocketConnectionFactory
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
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
		WsUseCase:     params.WsUseCase,
		WsUpgrader:    params.WsUpgrader,
		WsConnFactory: params.WsConnFactory,
		UserIDFactory: params.UserIDFactory,
		RoomIDFactory: params.RoomIDFactory,
		Logger:        params.Logger,
	}
}

func (h *WebSocketHandler) ConnectToChatRoom(c echo.Context) error {
	ctx := c.Request().Context()
	h.Logger.Info("ConnectToChatRoom called")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		h.Logger.Warn("User ID is missing or invalid")
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID is required")
	}

	roomID := c.Param("room_public_id")
	if roomID == "" {
		h.Logger.Warn("Room public ID is missing")
		return echo.NewHTTPError(http.StatusBadRequest, "Room public ID is required")
	}

	h.Logger.Info("Upgrading WebSocket connection", "room_public_id", roomID, "user_id", userID)
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

	if err := h.WsUseCase.ConnectUserToRoom(ctx, usecase.ConnectUserToRoomRequest{
		UserID: entity.UserID(userID),
		RoomID: entity.RoomID(roomID),
		Conn:   conn,
	}); err != nil {
		h.Logger.Error("Failed to connect user to room", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect user to room")
	}

	h.Logger.Info("User connected to room", "room_public_id", roomID, "user_id", userID)

	go func() {
		h.Logger.Info("Starting message loop", "room_public_id", roomID, "user_id", userID)
		// 新しいキャンセラブルな context を作成
		wsCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var userID = userID
		var roomID = roomID
		defer conn.Close()
		for {
			message, err := conn.ReadMessage()
			if err != nil {
				h.Logger.Warn("Connection closed or error reading message", "error", err)
				_ = h.WsUseCase.DisconnectUser(wsCtx, usecase.DisconnectUserRequest{
					UserID: entity.UserID(userID),
				})
				return
			}

			h.Logger.Info("Message received", "room_public_id", roomID, "user_id", userID)
			h.WsUseCase.SendMessage(wsCtx, usecase.SendMessageRequest{
				RoomID:  entity.RoomID(roomID),
				Sender:  entity.UserID(userID),
				Content: message.GetContent(),
			})
		}
	}()

	return nil
}
