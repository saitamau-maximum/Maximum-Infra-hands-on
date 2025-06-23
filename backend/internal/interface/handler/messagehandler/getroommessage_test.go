package messagehandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler/messagehandler"
	"example.com/infrahandson/internal/usecase/messagecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. c.Param("room_id") が空文字列の場合
// 3. c.QueryParam("limit")が strconv.Atoi で変換できない場合
// 4. c.QueryParam("before_sent_at") が空でも undefined でもなく，なおかつ time.Parase に失敗する
// 5. MsgUseCase.GetMessageHistoryInRoom がエラーを返す場合
func TestGetRoomMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := messagehandler.NewTestMessageHandler(ctrl)
	mockDeps.Logger.EXPECT().Info(gomock.Any()).AnyTimes()
	mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// 1. 正常系
	t.Run("正常系", func(t *testing.T) {
		roomID := "room123"
		now := time.Now()
		mockDeps.MsgUseCase.EXPECT().
			GetMessageHistoryInRoom(gomock.Any(), gomock.Any()).
			Return(messagecase.GetMessageHistoryInRoomResponse{
				Messages: []*entity.Message{
					entity.NewMessage(entity.MessageParams{
						ID:      "msg1",
						RoomID:  entity.RoomID(roomID),
						UserID:  "user1",
						Content: "Hello, World!",
						SentAt:  now,
					}),
				},
				NextBeforeSentAt: now,
				HasNext:          false,
			}, nil)

		req := httptest.NewRequest("GET", "/rooms/"+roomID+"/messages?limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/messages")
		c.SetParamNames("room_id")
		c.SetParamValues(roomID)

		if err := handler.GetRoomMessage(c); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected status code 200, got %d", rec.Code)
		}
	})

	// 2. room_id が空文字列の場合
	t.Run("room_id is empty", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rooms//messages", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/messages")
		c.SetParamNames("room_id")
		c.SetParamValues("")

		err := handler.GetRoomMessage(c)
		if err == nil {
			t.Error("expected error, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusBadRequest {
			t.Errorf("expected 400 BadRequest, got %v", err)
		}
	})

	// 3. limit が整数変換できない場合
	t.Run("limit is not integer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rooms/room123/messages?limit=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/messages")
		c.SetParamNames("room_id")
		c.SetParamValues("room123")

		err := handler.GetRoomMessage(c)
		if err == nil {
			t.Error("expected error, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusBadRequest {
			t.Errorf("expected 400 BadRequest, got %v", err)
		}
	})

	// 4. before_sent_at が time.Parse に失敗する値の場合
	t.Run("before_sent_at is invalid format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rooms/room123/messages?before_sent_at=invalid-time", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/messages")
		c.SetParamNames("room_id")
		c.SetParamValues("room123")

		err := handler.GetRoomMessage(c)
		if err == nil {
			t.Error("expected error, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusBadRequest {
			t.Errorf("expected 400 BadRequest, got %v", err)
		}
	})

	// 5. UseCase がエラーを返す場合
	t.Run("usecase returns error", func(t *testing.T) {
		roomID := "room123"
		mockDeps.MsgUseCase.EXPECT().
			GetMessageHistoryInRoom(gomock.Any(), gomock.Any()).
			Return(messagecase.GetMessageHistoryInRoomResponse{}, assert.AnError)

		req := httptest.NewRequest("GET", "/rooms/"+roomID+"/messages", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/messages")
		c.SetParamNames("room_id")
		c.SetParamValues(roomID)

		err := handler.GetRoomMessage(c)
		if err == nil {
			t.Error("expected error, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 InternalServerError, got %v", err)
		}
	})
}
