package handler

import (
	"net/http"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase"

	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	RoomUseCase   usecase.RoomUseCaseInterface
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
}

type NewRoomHandlerParams struct {
	RoomUseCase   usecase.RoomUseCaseInterface
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
}

func (p *NewRoomHandlerParams) Validate() error {
	if p.RoomUseCase == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "RoomUseCase is required")
	}
	if p.UserIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "UserIDFactory is required")
	}
	if p.RoomIDFactory == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "RoomPubIDFactory is required")
	}
	if p.Logger == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Logger is required")
	}
	return nil
}

func NewRoomHandler(params NewRoomHandlerParams) *RoomHandler {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &RoomHandler{
		RoomUseCase:   params.RoomUseCase,
		UserIDFactory: params.UserIDFactory,
		RoomIDFactory: params.RoomIDFactory,
		Logger:        params.Logger,
	}
}

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *RoomHandler) CreateRoom(c echo.Context) error {
	ctx := c.Request().Context()
	var req CreateRoomRequest

	if err := c.Bind(&req); err != nil {
		h.Logger.Error("Failed to bind request", err, req)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		h.Logger.Error("Validation failed", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation failed"})
	}

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		h.Logger.Error("User ID is missing or invalid")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is missing or invalid"})
	}

	// 部屋作成
	createRoomRes, err := h.RoomUseCase.CreateRoom(ctx, usecase.CreateRoomRequest{Name: req.Name})
	if err != nil {
		h.Logger.Error("Failed to create room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create room"})
	}
	room := createRoomRes.Room

	if err = h.RoomUseCase.JoinRoom(ctx, usecase.JoinRoomRequest{
		RoomID: room.GetID(),
		UserID: entity.UserID(userID),
	}); err != nil {
		h.Logger.Error("Failed to add user to room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to add user to room"})
	}

	// NOTE: WebSocketの接続 (ConnectUserToRoom) はこのタイミングでは行わない。
	// - 部屋の作成および論理参加（JoinRoom）のみを行う
	// - 実際のWebSocket接続はフロントエンド側で、部屋作成完了後に `/ws` へ接続する形で行う

	h.Logger.Info("Room created successfully", map[string]any{
		"roomID": room.GetID(),
	})

	return c.JSON(http.StatusOK, echo.Map{
		"roomID": room.GetID(),
	})
}

func (h *RoomHandler) JoinRoom(c echo.Context) error {
	ctx := c.Request().Context()
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		// user_id が存在しない、もしくは型アサーションに失敗した場合
		h.Logger.Error("User ID is missing or invalid")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is missing or invalid"})
	}

	roomID := c.Param("public_id")
	if roomID == "" {
		// public_id がリクエストパラメータに含まれていない場合
		h.Logger.Error("Room public ID is missing")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Room public ID is missing"})
	}

	// 部屋に参加
	err := h.RoomUseCase.JoinRoom(ctx, usecase.JoinRoomRequest{
		RoomID: entity.RoomID(roomID),
		UserID: entity.UserID(userID),
	})
	if err != nil {
		h.Logger.Error("Failed to join room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to join room"})
	}

	h.Logger.Info("Joined room successfully", map[string]any{
		"roomID": roomID,
		"userID": userID,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Joined room successfully",
	})
}

func (h *RoomHandler) LeaveRoom(c echo.Context) error {
	ctx := c.Request().Context()
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		h.Logger.Error("User ID is missing or invalid")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is missing or invalid"})
	}

	roomID := c.Param("public_id")
	if roomID == "" {
		h.Logger.Error("Room public ID is missing")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Room public ID is missing"})
	}

	// 部屋から退出
	if err := h.RoomUseCase.LeaveRoom(ctx, usecase.LeaveRoomRequest{
		RoomID: entity.RoomID(roomID),
		UserID: entity.UserID(userID),
	}); err != nil {
		h.Logger.Error("Failed to leave room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to leave room"})
	}

	h.Logger.Info("Left room successfully", map[string]any{
		"roomID": roomID,
		"userID": userID,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Left room successfully",
	})
}

type GetRoomResponse struct {
	PubID   string     `json:"public_id"`
	Name    string     `json:"name"`
	Members []MemberID `json:"members"`
}
type MemberID struct {
	ID string `json:"id"`
}

func (h *RoomHandler) GetRoom(c echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("public_id")
	if roomID == "" {
		h.Logger.Error("Room public ID is missing")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Room public ID is missing"})
	}

	GetRoomRes, err := h.RoomUseCase.GetRoomByID(ctx, usecase.GetRoomByIDRequest{
		ID: entity.RoomID(roomID),
	})
	if err != nil {
		h.Logger.Error("Failed to get room", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get room"})
	}

	room := GetRoomRes.Room

	res := GetRoomResponse{
		PubID:   string(room.GetID()),
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

type GetRoomsResponse struct {
	PubID string `json:"public_id"`
	Name  string `json:"name"`
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
			PubID: string(room.GetID()),
			Name:  room.GetName(),
		})
	}

	h.Logger.Info("Got rooms successfully")

	return c.JSON(http.StatusOK, res)
}
