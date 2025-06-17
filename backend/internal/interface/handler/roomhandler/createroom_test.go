package roomhandler_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"example.com/infrahandson/internal/usecase/roomcase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. Bind 失敗
// 3. バリデーション失敗
// 4. コンテキストにユーザーIDがない
// 5. ルーム作成失敗
// 6. ルーム参加失敗
func TestCreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := roomhandler.NewTestRoomHandler(ctrl)
	mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	t.Run("正常系", func(t *testing.T) {
		mockDeps.RoomUseCase.EXPECT().CreateRoom(gomock.Any(), gomock.Any()).Return(
			roomcase.CreateRoomResponse{
				Room: entity.NewRoom(entity.RoomParams{
					ID:      "room123",
					Name:    "Test Room",
					Members: []entity.UserID{},
				}),
			}, nil,
		)
		mockDeps.RoomUseCase.EXPECT().JoinRoom(gomock.Any(), gomock.Any()).Return(nil)

		req := httptest.NewRequest("POST", "/rooms", strings.NewReader(`{"name":"Test Room"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user123")

		if err := handler.CreateRoom(c); err != nil {
			t.Errorf("CreateRoom failed: %v", err)
		}
		if rec.Code != 200 {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Bind失敗", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/rooms", strings.NewReader(`{invalid json`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user123")

		err := handler.CreateRoom(c)
		if err == nil {
			t.Error("expected bind error, got nil")
		}
	})

	t.Run("バリデーション失敗", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/rooms", strings.NewReader(`{"name":""}`)) // 空の名前
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user123")

		err := handler.CreateRoom(c)
		if err == nil {
			t.Error("expected validation error, got nil")
		}
	})

	t.Run("ユーザーIDなし", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/rooms", strings.NewReader(`{"name":"Test Room"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateRoom(c)
		if err == nil {
			t.Error("expected error due to missing user_id, got nil")
		}
	})

	t.Run("CreateRoom失敗", func(t *testing.T) {
		mockDeps.RoomUseCase.EXPECT().CreateRoom(gomock.Any(), gomock.Any()).Return(
			roomcase.CreateRoomResponse{},
			assert.AnError, // モックのエラー
		)

		req := httptest.NewRequest("POST", "/rooms", strings.NewReader(`{"name":"Test Room"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user123")

		err := handler.CreateRoom(c)
		if err == nil {
			t.Error("expected CreateRoom error, got nil")
		}
	})

	t.Run("JoinRoom失敗", func(t *testing.T) {
		mockDeps.RoomUseCase.EXPECT().CreateRoom(gomock.Any(), gomock.Any()).Return(
			roomcase.CreateRoomResponse{
				Room: entity.NewRoom(entity.RoomParams{
					ID:      "room123",
					Name:    "Test Room",
					Members: []entity.UserID{},
				}),
			}, nil,
		)
		mockDeps.RoomUseCase.EXPECT().JoinRoom(gomock.Any(), gomock.Any()).Return(assert.AnError)

		req := httptest.NewRequest("POST", "/rooms", strings.NewReader(`{"name":"Test Room"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user123")

		err := handler.CreateRoom(c)
		if err == nil {
			t.Error("expected JoinRoom error, got nil")
		}
	})
}
