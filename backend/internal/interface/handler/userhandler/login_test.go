package userhandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/infrahandson/internal/interface/handler/userhandler"
	"example.com/infrahandson/internal/usecase/usercase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. バインド失敗
// 3. バリデーション失敗（例: emailが空）
// 4. 認証失敗（例: パスワードが間違っている）
func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := userhandler.NewTestUserHandler(ctrl)

	token := "mockToken"
	var tokenRes usercase.AuthenticateUserResponse
	tokenRes.SetToken(token)
	tokenRes.SetExp(3600)

	// 1. 正常系
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"test@example.com","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
		mockDeps.UserUseCase.EXPECT().AuthenticateUser(context.Background(), gomock.Any()).Return(tokenRes, nil)

		if assert.NoError(t, handler.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Login successful")
			cookies := rec.Result().Cookies()
			found := false
			for _, cookie := range cookies {
				if cookie.Name == "token" && cookie.Value == token {
					found = true
				}
			}
			assert.True(t, found, "token cookie should be set")
		}
	})

	// 2. バインド失敗
	t.Run("bind error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email": "test@example.com", "password": invalid}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request")
	})

	// 3. バリデーション失敗（例: emailが空）
	t.Run("validation error", func(t *testing.T) {
		body := map[string]string{
			"email":    "",
			"password": "password123",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Contains(t, rec.Body.String(), "Validation failed")
	})

	// 4. 認証失敗（例: パスワードが間違っている）
	t.Run("authentication failed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"test@example.com","password":"wrongpass"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
		mockDeps.UserUseCase.EXPECT().AuthenticateUser(context.Background(), gomock.Any()).Return(usercase.AuthenticateUserResponse{}, errors.New("invalid credentials"))

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Authentication failed")
	})
}
