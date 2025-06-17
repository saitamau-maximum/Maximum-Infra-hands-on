package roomhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetRoomsResponse struct {
	ID   string `json:"room_id"`
	Name string `json:"name"`
}

func (h *RoomHandler) GetRooms(c echo.Context) error {
	ctx := c.Request().Context()
	rooms, err := h.RoomUseCase.GetAllRooms(ctx)
	if err != nil {
		h.Logger.Error("Failed to get rooms", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get rooms"})
	}

	res := []GetRoomsResponse{}
	for _, room := range rooms {
		res = append(res, GetRoomsResponse{
			ID:   string(room.GetID()),
			Name: room.GetName(),
		})
	}

	h.Logger.Info("Got rooms successfully")

	return c.JSON(http.StatusOK, res)
}
