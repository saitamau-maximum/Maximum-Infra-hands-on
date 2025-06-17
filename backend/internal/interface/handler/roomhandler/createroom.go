package roomhandler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/roomcase"
	"github.com/labstack/echo/v4"
)

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

// CreateRoom は新しいルームを作成するハンドラーです。
// 名前からルームIDを生成し、ルームを作成します。
// ユーザーは自動的にそのルームに参加します。
func (h *RoomHandler) CreateRoom(c echo.Context) error {
	ctx := c.Request().Context()
	var req CreateRoomRequest

	if err := c.Bind(&req); err != nil {
		h.Logger.Error("Failed to bind request", err, req)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		h.Logger.Error("Validation failed", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation failed"})
	}

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		h.Logger.Error("User ID is missing or invalid")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is missing or invalid"})
	}

	// 部屋作成
	createRoomRes, err := h.RoomUseCase.CreateRoom(ctx, roomcase.CreateRoomRequest{Name: req.Name})
	if err != nil {
		h.Logger.Error("Failed to create room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create room"})
	}
	room := createRoomRes.Room

	if err = h.RoomUseCase.JoinRoom(ctx, roomcase.JoinRoomRequest{
		RoomID: room.GetID(),
		UserID: entity.UserID(userID),
	}); err != nil {
		h.Logger.Error("Failed to add user to room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to add user to room"})
	}

	// NOTE: WebSocketの接続 (ConnectUserToRoom) はこのタイミングでは行わない。
	// - 部屋の作成および論理参加（JoinRoom）のみを行う
	// - 実際のWebSocket接続はフロントエンド側で、部屋作成完了後に `/ws` へ接続する形で行う

	h.Logger.Info("Room created successfully", map[string]any{
		"roomID": room.GetID(),
	})

	return c.JSON(http.StatusOK, echo.Map{
		"roomID": room.GetID(),
	})
}
