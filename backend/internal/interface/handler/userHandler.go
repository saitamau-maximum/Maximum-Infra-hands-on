package handler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUseCase   usecase.UserUseCaseInterface
	UserIDFactory factory.UserIDFactory
	Logger        adapter.LoggerAdapter
}

type NewUserHandlerParams struct {
	UserUseCase   usecase.UserUseCaseInterface
	UserIDFactory factory.UserIDFactory
	Logger        adapter.LoggerAdapter
}

func (p *NewUserHandlerParams) Validate() error {
	if p.UserUseCase == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserUseCase is required")
	}
	if p.UserIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserIDFactory is required")
	}
	if p.Logger == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Logger is required")
	}
	return nil
}

func NewUserHandler(params NewUserHandlerParams) *UserHandler {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &UserHandler{
		UserUseCase:   params.UserUseCase,
		UserIDFactory: params.UserIDFactory,
		Logger:        params.Logger,
	}
}

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

	signUpReq := usecase.SignUpRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := h.UserUseCase.SignUp(ctx, signUpReq)
	if err != nil {
		h.Logger.Error("SignUp error: ", err)
		return c.JSON(500, echo.Map{"error": "Internal server error"})
	}

	authReq := usecase.AuthenticateUserRequest{
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

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "Validation failed"})
	}

	authReq := usecase.AuthenticateUserRequest{
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

type GetMeResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()
	uidRaw := c.Get("user_id")
	if uidRaw == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	userID, ok := uidRaw.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid user ID"})
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
