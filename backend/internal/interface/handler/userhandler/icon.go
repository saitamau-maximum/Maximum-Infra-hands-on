package userhandler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"github.com/labstack/echo/v4"
)

// SaveUserIcon: Saves the user's icon uploaded via a form file.
func (h *UserHandler) SaveUserIcon(c echo.Context) error {
	// UserID をコンテキストから取得して，userID 型までもっていく
	ctx := c.Request().Context()
	uidRaw := c.Get("user_id")
	if uidRaw == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	userIDStr, ok := uidRaw.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid user ID"})
	}
	userID := entity.UserID(userIDStr)

	// FromFile でリクエストからファイルを取得
	file, err := c.FormFile("icon")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// UseCase でアイコン保存
	err = h.UserUseCase.SaveUserIcon(ctx, file, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error/saveIcon": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Icon saved successfully"})
}

// GetUserIcon: Retrieves the user's icon URL based on the user ID from the request parameters.
func (h *UserHandler) GetUserIcon(c echo.Context) error {
	ctx := c.Request().Context()
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id is required"})
	}
	userID := entity.UserID(userIDStr)

	// Fetch the icon URL for the user
	iconURL, err := h.UserUseCase.GetUserIconPath(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Icon not found"})
	}

	h.Logger.Info(iconURL)

	// Redirect to the icon URL
	return c.Redirect(http.StatusFound, iconURL)
}
