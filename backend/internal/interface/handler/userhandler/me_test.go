package userhandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. コンテキストからID取得できない
// 3. string キャストできない不正なID
// 4. GetUserByID 失敗
func TestGetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.NewUser(entity.UserParams{
		ID:         "mockUserID",
		Name:       "Test User",
		Email:      "test@mail.com",
		PasswdHash: "hashedPassword",
		CreatedAt:  time.Now(),
		UpdatedAt:  nil,
	})

	handler, mockDeps, e := userhandler.NewTestUserHandler(ctrl)

	t.Run("正常系", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "mockUserID")

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
		mockDeps.UserUseCase.EXPECT().GetUserByID(gomock.Any(), entity.UserID("mockUserID")).Return(user, nil)

		if err := handler.GetMe(c); err != nil {
			t.Errorf("GetMe failed: %v", err)
		}

		if rec.Code != 200 {
			t.Errorf("Expected status code 200, got %d", rec.Code)
		}
	})

	t.Run("コンテキストからID取得できない", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

		if err := handler.GetMe(c); err != nil {
			t.Errorf("GetMe failed: %v", err)
		}

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %d", rec.Code)
		}
	})

	t.Run("string キャストできない不正なID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", 12345) // 不正なID

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

		if err := handler.GetMe(c); err != nil {
			t.Errorf("GetMe failed: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 401, got %d", rec.Code)
		}
	})

	t.Run("GetUserByID 失敗", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "mockUserID")

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
		mockDeps.UserUseCase.EXPECT().GetUserByID(gomock.Any(), entity.UserID("mockUserID")).Return(nil, errors.New("internal error"))

		if err := handler.GetMe(c); err != nil {
			t.Errorf("GetMe failed: %v", err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %d", rec.Code)
		}
	})
}
