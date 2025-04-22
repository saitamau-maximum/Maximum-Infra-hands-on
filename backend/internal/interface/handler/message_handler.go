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
}

type GetMessageHistoryInRoomRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_public_id"`
	Limit        int                 `json:"limit"`
	BeforeSentAt time.Time           `json:"before_sent_at"`
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
	var req GetMessageHistoryInRoomRequest
	roomPublicIDStr := c.Param("room_public_id")
	if roomPublicIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "room_public_id is required")
	}
	req.RoomPublicID = entity.RoomPublicID(roomPublicIDStr)

	// クエリ: limit（任意、デフォルト 10）
	req.Limit = 10
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limitNum, err := strconv.Atoi(limitStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer")
		}
		req.Limit = limitNum
	}

	// クエリ: before_sent_at（任意）
	beforeSentAtStr := c.QueryParam("before_sent_at")
	if beforeSentAtStr != "" {
		var err error
		req.BeforeSentAt, err = time.Parse(time.RFC3339, beforeSentAtStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "before_sent_at must be in RFC3339 format")
		}
	} else {
		req.BeforeSentAt = time.Now()
	}

	// Usecase呼び出し
	res, err := h.MsgUseCase.GetMessageHistoryInRoom(usecase.GetMessageHistoryInRoomRequest{
		RoomPublicID: entity.RoomPublicID(roomPublicIDStr),
		Limit:        req.Limit,
		BeforeSentAt: req.BeforeSentAt,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get message history")
	}

	// メッセージ整形
	messages := make([]MessageResponse, len(res.Messages))
	for i, msg := range res.Messages {
		formatedMsg, err := h.MsgUseCase.FormatMessage(msg)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to format message")
		}
		messages[i] = MessageResponse{
			ID:      formatedMsg.PublicID,
			UserID:  formatedMsg.UserPublicID,
			SentAt:  formatedMsg.SentAt.Format(time.RFC3339),
			Content: formatedMsg.Content,
		}
	}

	// レスポンス構築
	return c.JSON(http.StatusOK, GetMessageHistoryInRoomResponse{
		Messages:         messages,
		NextBeforeSentAt: res.NextBeforeSentAt.Format(time.RFC3339),
		HasNext:          res.HasNext,
	})
}
