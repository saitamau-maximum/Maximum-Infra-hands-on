package handler

import (
	"net/http"

	"example.com/webrtc-practice/internal/interface/factory"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUseCase   usecase.UserUseCaseInterface
	UserIDFactory factory.UserIDFactory
}

type NewUserHandlerParams struct {
	UserUseCase   usecase.UserUseCaseInterface
	UserIDFactory factory.UserIDFactory
}

func (p *NewUserHandlerParams) Validate() error {
	if p.UserUseCase == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserUseCase is required")
	}
	if p.UserIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserIDFactory is required")
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
	}
}

func RegisterUserRoutes(e *echo.Echo, handler *UserHandler) {
	e.POST("/register", handler.RegisterUser)
	e.POST("/login", handler.Login)
	e.GET("/me", handler.GetMe)
}

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	var req RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(400, echo.Map{"error": "Validation failed"})
	}

	signUpReq := usecase.SignUpRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := h.UserUseCase.SignUp(signUpReq)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "Internal server error"})
	}

	authReq := usecase.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// ログインまで済ませてしまう
	authRes, err := h.UserUseCase.AuthenticateUser(authReq)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "Internal server error"})
	}

	return c.JSON(200, echo.Map{"token": authRes.GetToken()})
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Login(c echo.Context) error {
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

	authRes, err := h.UserUseCase.AuthenticateUser(authReq)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authentication failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": authRes.GetToken(),
	})
}

type GetMeResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) GetMe(c echo.Context) error {
	uidRaw := c.Get("user_id")
	if uidRaw == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	userIDStr, ok := uidRaw.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid user ID"})
	}

	userID := h.UserIDFactory.FromString(userIDStr)

	user, err := h.UserUseCase.GetUserByID(userID)
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
