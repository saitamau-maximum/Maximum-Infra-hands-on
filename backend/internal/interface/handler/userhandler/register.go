package userhandler

import (
	"net/http"

	"example.com/infrahandson/internal/usecase/usercase"
	"github.com/labstack/echo/v4"
)

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()
	var req RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		h.Logger.Error("Failed to bind request: ", err)
		return c.JSON(400, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		h.Logger.Error("Validation failed: ", err)
		return c.JSON(400, echo.Map{"error": "Validation failed"})
	}

	signUpReq := usercase.SignUpRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := h.UserUseCase.SignUp(ctx, signUpReq)
	if err != nil {
		h.Logger.Error("SignUp error: ", err)
		return c.JSON(500, echo.Map{"error": "Internal server error"})
	}

	authReq := usercase.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// ログインまで済ませてしまう
	authRes, err := h.UserUseCase.AuthenticateUser(ctx, authReq)
	if err != nil {
		h.Logger.Error("Authentication error: ", err)
		return c.JSON(500, echo.Map{"error": err.Error()})
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
