package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/usecase"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_usecase "example.com/infrahandson/test/mocks/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetMessageHistoryInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockUseCase := mock_usecase.NewMockMessageUseCaseInterface(ctrl)
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	handler := handler.NewMessageHandler(handler.NewMessageHandlerParams{
		MsgUseCase: mockUseCase,
		Logger:     mockLogger,
	})

	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes() // ロガーは何回呼ばれてもいい（呼ばれなくても怒らない）設定

	roomPublicID := "test-room"
	mockMessages := []*entity.Message{
		entity.NewMessage(entity.MessageParams{
			ID:       entity.MessageID(1),
			PublicID: entity.MessagePublicID("message1"),
			RoomID:   entity.RoomID(1),
			UserID:   entity.UserID(1),
			Content:  "Hello",
			SentAt:   time.Now(),
		}),
		entity.NewMessage(entity.MessageParams{
			ID:       entity.MessageID(2),
			PublicID: entity.MessagePublicID("message2"),
			RoomID:   entity.RoomID(1),
			UserID:   entity.UserID(2),
			Content:  "Hi",
			SentAt:   time.Now(),
		}),
	}
	mockResponse := usecase.GetMessageHistoryInRoomResponse{
		Messages:         mockMessages,
		NextBeforeSentAt: time.Now(),
		HasNext:          true,
	}

	mockUseCase.EXPECT().
		GetMessageHistoryInRoom(gomock.Any()).
		Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages/"+roomPublicID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_public_id")
	c.SetParamValues(roomPublicID)

	mockUseCase.EXPECT().FormatMessage(mockMessages[0]).Return(
		usecase.FormatMessageResponse{
			PublicID:     string(mockMessages[0].GetPublicID()),
			RoomPublicID: "test-room",
			UserPublicID: "user1",
			Content:      mockMessages[0].GetContent(),
			SentAt:       mockMessages[0].GetSentAt(),
		}, nil,
	)
	mockUseCase.EXPECT().FormatMessage(mockMessages[1]).Return(
		usecase.FormatMessageResponse{
			PublicID:     string(mockMessages[1].GetPublicID()),
			RoomPublicID: "test-room",
			UserPublicID: "user2",
			Content:      mockMessages[1].GetContent(),
			SentAt:       mockMessages[1].GetSentAt(),
		}, nil,
	)

	err := handler.GetMessageHistoryInRoom(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, strings.Contains(rec.Body.String(), `"has_next":true`))
}
