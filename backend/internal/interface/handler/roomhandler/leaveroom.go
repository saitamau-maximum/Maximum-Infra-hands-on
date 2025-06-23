package roomhandler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/roomcase"
	"github.com/labstack/echo/v4"
)

func (h *RoomHandler) LeaveRoom(c echo.Context) error {
	ctx := c.Request().Context()
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		h.Logger.Error("User ID is missing or invalid")
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is missing or invalid")
	}

	roomID := c.Param("room_id")
	if roomID == "" {
		h.Logger.Error("Room ID is missing")
		return echo.NewHTTPError(http.StatusBadRequest, "Room ID is missing")
	}

	// 部屋から退出
	if err := h.RoomUseCase.LeaveRoom(ctx, roomcase.LeaveRoomRequest{
		RoomID: entity.RoomID(roomID),
		UserID: entity.UserID(userID),
	}); err != nil {
		h.Logger.Error("Failed to leave room", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to leave room")
	}

	h.Logger.Info("Left room successfully", map[string]any{
		"roomID": roomID,
		"userID": userID,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Left room successfully",
	})
}
