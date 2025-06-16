package userhandler

import (
	"net/http"

	"example.com/infrahandson/internal/usecase/usercase"
	"github.com/labstack/echo/v4"
)

// LoginRequest represents the structure of the login request payload.
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login: クッキーをセットしてログイン処理
func (h *UserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "Validation failed"})
	}

	authReq := usercase.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	authRes, err := h.UserUseCase.AuthenticateUser(ctx, authReq)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authentication failed"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    authRes.GetToken(),
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   authRes.GetExp(),
	})

	return c.JSON(http.StatusOK, echo.Map{"message": "Login successful"})
}
