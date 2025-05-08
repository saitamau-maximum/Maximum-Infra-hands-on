package handler_test

import (
	"context"
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

	roomID := "test-room"
	mockMessages := []*entity.Message{
		entity.NewMessage(entity.MessageParams{
			ID:      entity.MessageID("message1"),
			RoomID:  entity.RoomID("test-room"),
			UserID:  entity.UserID("user1"),
			Content: "Hello",
			SentAt:  time.Now(),
		}),
		entity.NewMessage(entity.MessageParams{
			ID:      entity.MessageID("message2"),
			RoomID:  entity.RoomID("test-room"),
			UserID:  entity.UserID("user2"),
			Content: "Hi",
			SentAt:  time.Now(),
		}),
	}
	mockResponse := usecase.GetMessageHistoryInRoomResponse{
		Messages:         mockMessages,
		NextBeforeSentAt: time.Now(),
		HasNext:          true,
	}

	mockUseCase.EXPECT().
		GetMessageHistoryInRoom(context.Background(), gomock.Any()).
		Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages/"+roomID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_public_id")
	c.SetParamValues(roomID)

	err := handler.GetMessageHistoryInRoom(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, strings.Contains(rec.Body.String(), `"has_next":true`))
}
