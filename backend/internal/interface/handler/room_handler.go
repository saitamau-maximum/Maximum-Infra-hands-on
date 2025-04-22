package handler

import (
	"net/http"

	"example.com/webrtc-practice/internal/interface/factory"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	RoomUseCase   usecase.RoomUseCaseInterface
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
}

type NewRoomHandlerParams struct {
	RoomUseCase   usecase.RoomUseCaseInterface
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
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
	}
}

func (h *RoomHandler) Register(g *echo.Group) {
	g.POST("/rooms", h.CreateRoom)
	g.POST("/rooms/:room_public_id/join", h.JoinRoom)
	g.POST("/rooms/:room_public_id/leave", h.LeaveRoom)
	g.GET("/rooms/:room_public_id", h.GetRoom)
	g.GET("/rooms", h.GetRooms)
}

type CreateRoomRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *RoomHandler) CreateRoom(c echo.Context) error {
	var req CreateRoomRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation failed"})
	}

	// 部屋作成
	createRoomRes, err := h.RoomUseCase.CreateRoom(usecase.CreateRoomRequest{Name: req.Name})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create room"})
	}
	room := createRoomRes.Room

	// 部屋に作成者を追加
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is missing or invalid"})
	}
	userPublicID, err := h.UserIDFactory.FromString(userIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get user by public ID"})
	}

	if err = h.RoomUseCase.JoinRoom(usecase.JoinRoomRequest{
		RoomPublicID: room.GetPubID(),
		UserPublicID: userPublicID,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to add user to room"})
	}

	// NOTE: WebSocketの接続 (ConnectUserToRoom) はこのタイミングでは行わない。
	// - 部屋の作成および論理参加（JoinRoom）のみを行う
	// - 実際のWebSocket接続はフロントエンド側で、部屋作成完了後に `/ws` へ接続する形で行う

	return c.JSON(http.StatusOK, echo.Map{
		"roomPublicID": room.GetPubID(),
	})
}

func (h *RoomHandler) JoinRoom(c echo.Context) error {
	userPublicIDStr, ok := c.Get("user_id").(string)
	if !ok || userPublicIDStr == "" {
		// user_id が存在しない、もしくは型アサーションに失敗した場合
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is missing or invalid"})
	}
	userPublicID, err := h.UserIDFactory.FromString(userPublicIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get user by public ID"})
	}

	roomPublicIDStr := c.Param("public_id")
	if roomPublicIDStr == "" {
		// public_id がリクエストパラメータに含まれていない場合
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Room public ID is missing"})
	}
	roomPublicID, err := h.RoomIDFactory.FromString(roomPublicIDStr)
	if err != nil {
		// public_id が不正な形式の場合
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid room public ID"})
	}

	// 部屋に参加
	err = h.RoomUseCase.JoinRoom(usecase.JoinRoomRequest{
		RoomPublicID: roomPublicID,
		UserPublicID: userPublicID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to join room"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Joined room successfully",
	})
}

func (h *RoomHandler) LeaveRoom(c echo.Context) error {
	userPublicIDStr, ok := c.Get("user_id").(string)
	if !ok || userPublicIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is missing or invalid"})
	}
	userPublicID, err := h.UserIDFactory.FromString(userPublicIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get user by public ID"})
	}

	roomPublicIDStr := c.Param("public_id")
	if roomPublicIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Room public ID is missing"})
	}
	roomPublicID, err := h.RoomIDFactory.FromString(roomPublicIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid room public ID"})
	}

	// 部屋から退出
	if err := h.RoomUseCase.LeaveRoom(usecase.LeaveRoomRequest{
		RoomPublicID: roomPublicID,
		UserPublicID: userPublicID,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to leave room"})
	}

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
	roomPublicIDStr := c.Param("public_id")
	if roomPublicIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Room public ID is missing"})
	}
	roomPublicID, err := h.RoomIDFactory.FromString(roomPublicIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid room public ID"})
	}

	GetRoomRes, err := h.RoomUseCase.GetRoomByPublicID(usecase.GetRoomByPublicIDParams{
		PublicID: roomPublicID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get room"})
	}

	room := GetRoomRes.Room

	res := GetRoomResponse{
		PubID:   string(room.GetPubID()),
		Name:    room.GetName(),
		Members: []MemberID{},
	}

	for _, memberID := range room.GetMembers() {
		res.Members = append(res.Members, MemberID{
			ID: string(memberID),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"room": res,
	})
}

type GetRoomsResponse struct {
	PubID string `json:"public_id"`
	Name  string `json:"name"`
}

type RoomInfo struct {
	PubID string `json:"public_id"`
	Name  string `json:"name"`
}

func (h *RoomHandler) GetRooms(c echo.Context) error {
	rooms, err := h.RoomUseCase.GetAllRooms()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get rooms"})
	}

	res := []GetRoomsResponse{}
	for _, room := range rooms {
		res = append(res, GetRoomsResponse{
			PubID: string(room.GetPubID()),
			Name:  room.GetName(),
		})
	}

	return c.JSON(http.StatusOK, res)
}
