package userhandler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"github.com/labstack/echo/v4"
)

type GetMeResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetMe: Fetches the current user's information based on the user ID stored in the context.
func (h *UserHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()
	uidRaw := c.Get("user_id")
	if uidRaw == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	userID, ok := uidRaw.(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	user, err := h.UserUseCase.GetUserByID(ctx, entity.UserID(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	res := GetMeResponse{
		ID:    string(user.GetID()),
		Name:  user.GetName(),
		Email: user.GetEmail(),
	}

	return c.JSON(http.StatusOK, res)
}
