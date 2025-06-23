package roomhandler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/roomcase"
	"github.com/labstack/echo/v4"
)

func (h *RoomHandler) JoinRoom(c echo.Context) error {
	ctx := c.Request().Context()
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		// user_id が存在しない、もしくは型アサーションに失敗した場合
		h.Logger.Error("User ID is missing or invalid")
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID is missing or invalid")
	}

	roomID := c.Param("room_id")
	if roomID == "" {
		// room_id がリクエストパラメータに含まれていない場合
		h.Logger.Error("Room ID is missing")
		return echo.NewHTTPError(http.StatusBadRequest, "Room ID is missing")
	}

	// 部屋に参加
	err := h.RoomUseCase.JoinRoom(ctx, roomcase.JoinRoomRequest{
		RoomID: entity.RoomID(roomID),
		UserID: entity.UserID(userID),
	})
	if err != nil {
		h.Logger.Error("Failed to join room", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to join room")
	}

	h.Logger.Info("Joined room successfully", map[string]any{
		"roomID": roomID,
		"userID": userID,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Joined room successfully",
	})
}
