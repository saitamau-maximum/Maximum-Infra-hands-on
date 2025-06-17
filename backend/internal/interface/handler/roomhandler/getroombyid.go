package roomhandler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/roomcase"
	"github.com/labstack/echo/v4"
)

type GetRoomResponse struct {
	ID      string     `json:"room_id"`
	Name    string     `json:"name"`
	Members []MemberID `json:"members"`
}
type MemberID struct {
	ID string `json:"id"`
}

func (h *RoomHandler) GetRoomByID(c echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("room_id")
	if roomID == "" {
		h.Logger.Error("Room ID is missing")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Room ID is missing"})
	}

	GetRoomRes, err := h.RoomUseCase.GetRoomByID(ctx, roomcase.GetRoomByIDRequest{
		ID: entity.RoomID(roomID),
	})
	if err != nil {
		h.Logger.Error("Failed to get room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get room"})
	}

	room := GetRoomRes.Room

	res := GetRoomResponse{
		ID:      string(room.GetID()),
		Name:    room.GetName(),
		Members: []MemberID{},
	}

	for _, memberID := range room.GetMembers() {
		res.Members = append(res.Members, MemberID{
			ID: string(memberID),
		})
	}

	h.Logger.Info("Got room successfully", map[string]any{
		"roomID": roomID,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"room": res,
	})
}
