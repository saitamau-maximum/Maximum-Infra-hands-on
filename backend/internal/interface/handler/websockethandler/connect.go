package websockethandler

import (
	"context"
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/websocketcase"
	"github.com/labstack/echo/v4"
)

func (h *WebSocketHandler) ConnectToChatRoom(c echo.Context) error {
	ctx := c.Request().Context()
	h.Logger.Info("ConnectToChatRoom called")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		h.Logger.Warn("User ID is missing or invalid")
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID is required")
	}

	roomID := c.Param("room_id")
	if roomID == "" {
		h.Logger.Warn("Room ID is missing")
		return echo.NewHTTPError(http.StatusBadRequest, "Room ID is required")
	}

	h.Logger.Info("Upgrading WebSocket connection", "room_id", roomID, "user_id", userID)
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

	if err := h.WsUseCase.ConnectUserToRoom(ctx, websocketcase.ConnectUserToRoomRequest{
		UserID: entity.UserID(userID),
		RoomID: entity.RoomID(roomID),
		Conn:   conn,
	}); err != nil {
		h.Logger.Error("Failed to connect user to room", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect user to room")
	}

	h.Logger.Info("User connected to room", "room_id", roomID, "user_id", userID)

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
				_ = h.WsUseCase.DisconnectUser(wsCtx, websocketcase.DisconnectUserRequest{
					UserID: entity.UserID(userID),
				})
				return
			}

			h.Logger.Info("Message received", "room_public_id", roomID, "user_id", userID)
			err = h.WsUseCase.SendMessage(wsCtx, websocketcase.SendMessageRequest{
				RoomID:  entity.RoomID(roomID),
				Sender:  entity.UserID(userID),
				Content: message.GetContent(),
			})
			if err != nil {
				h.Logger.Error("connection closed", err)
				break
			}
		}
	}()

	return nil
}
