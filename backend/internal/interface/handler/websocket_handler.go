package handler

import (
	"net/http"

	"example.com/webrtc-practice/internal/interface/adapter"
	"example.com/webrtc-practice/internal/interface/factory"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	WsUseCase        usecase.WebsocketUseCaseInterface
	WsUpgrader       adapter.WebSocketUpgraderAdapter
	WsConnFactory    factory.WebSocketConnectionFactory
	UserIDFactory    factory.UserIDFactory
	RoomPubIDFactory factory.RoomIDFactory
}

type WebSocketHandlerParams struct {
	WsUseCase        usecase.WebsocketUseCaseInterface
	WsUpgrader       adapter.WebSocketUpgraderAdapter
	WsConnFactory    factory.WebSocketConnectionFactory
	UserIDFactory    factory.UserIDFactory
	RoomPubIDFactory factory.RoomIDFactory
}

func (p *WebSocketHandlerParams) Validate() error {
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
	if p.RoomPubIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "RoomPubIDFactory is required")
	}
	return nil
}

func NewWebSocketHandler(params WebSocketHandlerParams) *WebSocketHandler {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &WebSocketHandler{
		WsUseCase:        params.WsUseCase,
		WsUpgrader:       params.WsUpgrader,
		WsConnFactory:    params.WsConnFactory,
		UserIDFactory:    params.UserIDFactory,
		RoomPubIDFactory: params.RoomPubIDFactory,
	}
}

func (h *WebSocketHandler) Register(g *echo.Group) {
	g.GET("/ws/:room_public_id", h.ConnectToChatRoom)
}

func (h *WebSocketHandler) ConnectToChatRoom(c echo.Context) error {
	userPublicIDStr, ok := c.Get("user_id").(string)
	if !ok || userPublicIDStr == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID is required")
	}
	userPublicID, err := h.UserIDFactory.FromString(userPublicIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid User ID")
	}

	roomPublicIDStr := c.Param("room_public_id")
	if roomPublicIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Room public ID is required")
	}
	roomPublicID, err := h.RoomPubIDFactory.FromString(roomPublicIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Room public ID")
	}

	connRaw, err := h.WsUpgrader.Upgrade(c.Response().Writer, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upgrade connection")
	}

	conn, err := h.WsConnFactory.CreateWebSocketConnection(connRaw)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create WebSocket connection")
	}

	if err := h.WsUseCase.ConnectUserToRoom(usecase.ConnectUserToRoomRequest{
		UserPublicID: userPublicID,
		RoomPublicID: roomPublicID,
		Conn:         conn,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect user to room")
	}

	go func() {
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				_ = h.WsUseCase.DisconnectUser(usecase.DisconnectUserRequest{
					UserID: userPublicID,
				})
				return
			}

			h.WsUseCase.SendMessage(usecase.SendMessageRequest{
				RoomPublicID: roomPublicID,
				Sender:       userPublicID,
				Content:      message.GetContent(),
			})
		}
	}()

	return nil
}
