package userhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *UserHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusOK, echo.Map{"message": "Logout successful"})
}
