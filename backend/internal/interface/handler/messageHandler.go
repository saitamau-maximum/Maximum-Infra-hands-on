package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/usecase"

	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	MsgUseCase usecase.MessageUseCaseInterface
	Logger     adapter.LoggerAdapter
}

type NewMessageHandlerParams struct {
	MsgUseCase usecase.MessageUseCaseInterface
	Logger     adapter.LoggerAdapter
}

func (p *NewMessageHandlerParams) Validate() error {
	if p.MsgUseCase == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "MessageUseCase is required")
	}
	if p.Logger == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Logger is required")
	}
	return nil
}

func NewMessageHandler(params NewMessageHandlerParams) *MessageHandler {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &MessageHandler{
		MsgUseCase: params.MsgUseCase,
		Logger:     params.Logger,
	}
}

type GetMessageHistoryInRoomRequest struct {
	RoomID       entity.RoomID `json:"room_id"`
	Limit        int           `json:"limit"`
	BeforeSentAt time.Time     `json:"before_sent_at"`
}

type GetMessageHistoryInRoomResponse struct {
	Messages         []MessageResponse `json:"messages"`
	NextBeforeSentAt string           `json:"next_before_sent_at"`
	HasNext          bool             `json:"has_next"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	SentAt    time.Time `json:"sent_at"`
}

func (h *MessageHandler) GetMessageHistoryInRoom(c echo.Context) error {
	ctx := c.Request().Context()
	h.Logger.Info("GetMessageHistoryInRoom called")
	var req GetMessageHistoryInRoomRequest
	roomIDStr := c.Param("room_public_id")
	if roomIDStr == "" {
		h.Logger.Error("room_public_id is required")
		return echo.NewHTTPError(http.StatusBadRequest, "room_public_id is required")
	}
	req.RoomID = entity.RoomID(roomIDStr)

	// クエリ: limit（任意、デフォルト 10）
	req.Limit = 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limitNum, err := strconv.Atoi(limitStr)
		if err != nil {
			h.Logger.Error("limit must be an integer")
			return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer")
		}
		req.Limit = limitNum
	}

	// クエリ: before_sent_at（任意）
	beforeSentAtStr := c.QueryParam("before_sent_at")

	// Check for "undefined" as a workaround for cases where the frontend or external system
	// sends the string "undefined" instead of leaving the parameter empty.
	if beforeSentAtStr != "" && beforeSentAtStr != "undefined" {
		fixedStr := strings.Replace(beforeSentAtStr, " ", "+", 1)
		var err error
		req.BeforeSentAt, err = time.Parse(time.RFC3339, fixedStr)
		if err != nil {
			h.Logger.Error("before_sent_at must be in RFC3339 format")
			return echo.NewHTTPError(http.StatusBadRequest, "before_sent_at must be in RFC3339 format")
		}
	} else {
		req.BeforeSentAt = time.Now()
	}

	// Usecase呼び出し
	res, err := h.MsgUseCase.GetMessageHistoryInRoom(ctx, usecase.GetMessageHistoryInRoomRequest{
		RoomID:       req.RoomID,
		Limit:        req.Limit,
		BeforeSentAt: req.BeforeSentAt,
	})
	if err != nil {
		h.Logger.Error("failed to get message history: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get message history")
	}
	
	// レスポンス構築
	messages := make([]MessageResponse, len(res.Messages))
	for i, msg := range res.Messages {
		messages[i] = MessageResponse{
			ID:      string(msg.GetID()),
			RoomID:  string(msg.GetRoomID()),
			UserID:  string(msg.GetUserID()),
			Content: msg.GetContent(),
			SentAt:  msg.GetSentAt(),
		}
	}
	return c.JSON(http.StatusOK, GetMessageHistoryInRoomResponse{
		Messages:         messages,
		NextBeforeSentAt: res.NextBeforeSentAt.Format(time.RFC3339),
		HasNext:          res.HasNext,
	})
}
