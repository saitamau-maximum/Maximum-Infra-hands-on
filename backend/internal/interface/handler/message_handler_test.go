package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/interface/handler"
	"example.com/webrtc-practice/internal/usecase"
	mock_usecase "example.com/webrtc-practice/mocks/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetMessageHistoryInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockUseCase := mock_usecase.NewMockMessageUseCaseInterface(ctrl)
	handler := handler.NewMessageHandler(handler.MessageHandlerParams{MsgUseCase: mockUseCase})

	roomPublicID := "test-room"
	mockMessages := []*entity.Message{
		entity.NewMessage(entity.MessageParams{
			ID:      "1",
			RoomID:  1,
			UserID:  "user1",
			Content: "Hello",
			SentAt:  time.Now(),
		}),
		entity.NewMessage(entity.MessageParams{
			ID:      "2",
			RoomID:  1,
			UserID:  "user2",
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
		GetMessageHistoryInRoom(gomock.Any()).
		Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages/"+roomPublicID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_public_id")
	c.SetParamValues(roomPublicID)

	err := handler.GetMessageHistoryInRoom(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, strings.Contains(rec.Body.String(), `"has_next":true`))
}

func TestGetNextMessageHistoryInRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockUseCase := mock_usecase.NewMockMessageUseCaseInterface(ctrl)
	handler := handler.NewMessageHandler(handler.MessageHandlerParams{MsgUseCase: mockUseCase})

	roomPublicID := "test-room"
	beforeSentAt := time.Now().Add(-time.Hour).UTC() // Ensure UTC for RFC3339 compliance
	mockMessages := []*entity.Message{
		entity.NewMessage(entity.MessageParams{
			ID:      "3",
			RoomID:  1,
			UserID:  "user3",
			Content: "Good morning",
			SentAt:  time.Now(),
		}),
		entity.NewMessage(entity.MessageParams{
			ID:      "4",
			RoomID:  1,
			UserID:  "user4",
			Content: "Good evening",
			SentAt:  time.Now(),
		}),
	}
	mockResponse := usecase.GetMessageHistoryInRoomResponse{
		Messages:         mockMessages,
		NextBeforeSentAt: time.Now(),
		HasNext:          false,
	}

	mockUseCase.EXPECT().
		GetMessageHistoryInRoom(gomock.Any()).
		Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages/"+roomPublicID+"/next?before_sent_at="+beforeSentAt.Format(time.RFC3339), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("room_public_id")
	c.SetParamValues(roomPublicID)

	err := handler.GetNextMessageHistoryInRoom(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, strings.Contains(rec.Body.String(), `"has_next":false`))
}
