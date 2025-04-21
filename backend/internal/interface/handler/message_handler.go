package handler

import (
	"net/http"
	"strconv"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	MsgUseCase usecase.MessageUseCaseInterface
}

type MessageHandlerParams struct {
	MsgUseCase usecase.MessageUseCaseInterface
}

func (p *MessageHandlerParams) Validate() error {
	if p.MsgUseCase == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "MessageUseCase is required")
	}
	return nil
}

func NewMessageHandler(params MessageHandlerParams) *MessageHandler {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &MessageHandler{
		MsgUseCase: params.MsgUseCase,
	}
}

func (h *MessageHandler) Register(g *echo.Group) {
	g.GET("/messages/:room_public_id", h.GetMessageHistoryInRoom)
	g.GET("/messages/:room_public_id/next", h.GetNextMessageHistoryInRoom)
}

type GetMessageHistoryInRoomResponse struct {
	Messages         []MessageResponse `json:"messages"`
	NextBeforeSentAt string            `json:"next_before_sent_at"`
	HasNext          bool              `json:"has_next"`
}
type MessageResponse struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	SentAt  string `json:"sent_at"`
	Content string `json:"content"`
}

func (h *MessageHandler) GetMessageHistoryInRoom(c echo.Context) error {
	roomPublicIDStr := c.Param("room_public_id")
	if roomPublicIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "room_public_id is required")
	}

	limit := 10 // デフォルトの取得件数
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limitNum, err := strconv.Atoi(limitStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer")
		}
		limit = limitNum
	}

	beforeSentAtStr := c.QueryParam("before_sent_at")
	var beforeSentAt time.Time
	if beforeSentAtStr != "" {
		var err error
		beforeSentAt, err = time.Parse(time.RFC3339, beforeSentAtStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "before_sent_at must be in RFC3339 format")
		}
	}

	res, err := h.MsgUseCase.GetMessageHistoryInRoom(usecase.GetMessageHistoryInRoomRequest{
		RoomPublicID: entity.RoomPublicID(roomPublicIDStr),
		Limit:        limit,
		BeforeSentAt: beforeSentAt,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get message history")
	}

	messages := make([]MessageResponse, len(res.Messages))
	for i, msg := range res.Messages {
		messages[i] = MessageResponse{
			ID:      string(msg.GetID()),
			UserID:  string(msg.GetUserID()),
			SentAt:  msg.GetSentAt().Format(time.RFC3339),
			Content: msg.GetContent(),
		}
	}

	return c.JSON(http.StatusOK, GetMessageHistoryInRoomResponse{
		Messages:         messages,
		NextBeforeSentAt: res.NextBeforeSentAt.Format(time.RFC3339),
		HasNext:          res.HasNext,
	})
}

func (h *MessageHandler) GetNextMessageHistoryInRoom(c echo.Context) error {
	roomPublicIDStr := c.Param("room_public_id")
	if roomPublicIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "room_public_id is required")
	}

	beforeSentAtStr := c.QueryParam("before_sent_at")
	if beforeSentAtStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "before_sent_at is required")
	}

	beforeSentAt, err := time.Parse(time.RFC3339, beforeSentAtStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "before_sent_at must be in RFC3339 format")
	}

	res, err := h.MsgUseCase.GetMessageHistoryInRoom(usecase.GetMessageHistoryInRoomRequest{
		RoomPublicID: entity.RoomPublicID(roomPublicIDStr),
		Limit:        30,
		BeforeSentAt: beforeSentAt,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get message history")
	}

	messages := make([]MessageResponse, len(res.Messages))
	for i, msg := range res.Messages {
		messages[i] = MessageResponse{
			ID:      string(msg.GetID()),
			UserID:  string(msg.GetUserID()),
			SentAt:  msg.GetSentAt().Format(time.RFC3339),
			Content: msg.GetContent(),
		}
	}

	return c.JSON(http.StatusOK, GetMessageHistoryInRoomResponse{
		Messages:         messages,
		NextBeforeSentAt: res.NextBeforeSentAt.Format(time.RFC3339),
		HasNext:          res.HasNext,
	})
}
