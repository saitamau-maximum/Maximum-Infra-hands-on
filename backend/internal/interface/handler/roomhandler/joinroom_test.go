package roomhandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"go.uber.org/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// 1. 正常系
// 2. ユーザーIDが存在しない場合
// 3. user_id が型アサーションに失敗した場合
// 4. room_id がリクエストパラメータに含まれていない場合
// 5. 部屋に参加できなかった場合
func TestJoinRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := roomhandler.NewTestRoomHandler(ctrl)

	// 1. 正常系
	t.Run("正常系", func(t *testing.T) {
		mockDeps.RoomUseCase.EXPECT().JoinRoom(gomock.Any(), gomock.Any()).Return(nil)
		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest(http.MethodPost, "/rooms/room123/join", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/join")
		c.SetParamNames("room_id")
		c.SetParamValues("room123")
		c.Set("user_id", "user123")

		err := handler.JoinRoom(c)
		assert.NoError(t, err)
		assert.Equal(t, 200, rec.Code)
	})

	// 2. ユーザーIDが存在しない場合
	t.Run("ユーザーIDが存在しない場合", func(t *testing.T) {
		mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest(http.MethodPost, "/rooms/room123/join", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/join")
		c.SetParamNames("room_id")
		c.SetParamValues("room123")

		err := handler.JoinRoom(c)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	// 3. user_id が型アサーションに失敗した場合
	t.Run("user_id が型アサーションに失敗した場合", func(t *testing.T) {
		mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest(http.MethodPost, "/rooms/room123/join", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/join")
		c.SetParamNames("room_id")
		c.SetParamValues("room123")
		c.Set("user_id", 123) // stringでない

		err := handler.JoinRoom(c)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	// 4. room_id がリクエストパラメータに含まれていない場合
	t.Run("room_id が空の場合", func(t *testing.T) {
		mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest(http.MethodPost, "/rooms//join", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/join")
		c.SetParamNames("room_id")
		c.SetParamValues("")
		c.Set("user_id", "user123")

		err := handler.JoinRoom(c)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	// 5. 部屋に参加できなかった場合
	t.Run("部屋に参加できなかった場合", func(t *testing.T) {
		mockDeps.RoomUseCase.EXPECT().JoinRoom(gomock.Any(), gomock.Any()).Return(errors.New("failed to join"))
		mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest(http.MethodPost, "/rooms/room123/join", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id/join")
		c.SetParamNames("room_id")
		c.SetParamValues("room123")
		c.Set("user_id", "user123")

		err := handler.JoinRoom(c)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
	})
}
