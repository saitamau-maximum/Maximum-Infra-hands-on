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

// 1.正常系
// 2.バインド失敗
// 3.バリデーション失敗（例: nameが空）
// 4.サインアップUseCaseのエラー
// 5.認証UseCaseのエラー

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := userhandler.NewTestUserHandler(ctrl)

	token := "mockToken"
	var tokenRes usercase.AuthenticateUserResponse
	tokenRes.SetToken(token)

	// 1. 正常系
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"name":"test","email":"test@example.com","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
		mockDeps.UserUseCase.EXPECT().SignUp(context.Background(), gomock.Any()).Return(usercase.SignUpResponse{}, nil)
		mockDeps.UserUseCase.EXPECT().AuthenticateUser(context.Background(), gomock.Any()).Return(tokenRes, nil)

		if assert.NoError(t, handler.RegisterUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Login successful")
		}
	})

	// 2. バインド失敗
	t.Run("bind error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email": "test@example.com", "password": invalid}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

		err := handler.RegisterUser(c)
		assert.NoError(t, err) // Echoはエラーを返す代わりにHTTPレスポンスを返すので、エラー自体はnil
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})


	// 3. バリデーション失敗（例: nameが空）
	t.Run("validation error", func(t *testing.T) {
		body := map[string]string{
			"name":     "",
			"email":    "test@example.com",
			"password": "password123",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

		err := handler.RegisterUser(c)
		assert.NoError(t, err) // Echoはエラーを返す代わりにHTTPレスポンスを返すので、エラー自体はnil
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// 4. サインアップUseCaseのエラー
	t.Run("signup usecase error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"name":"test","email":"test@example.com","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
		mockDeps.UserUseCase.EXPECT().SignUp(context.Background(), gomock.Any()).Return(usercase.SignUpResponse{}, errors.New("signup error"))
		
		err := handler.RegisterUser(c)
		assert.NoError(t, err) // Echoはエラーを返す代わりにHTTPレスポンスを返すので、エラー自体はnil
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	// 5. 認証UseCaseのエラー
	t.Run("authenticate usecase error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"name":"test","email":"test@example.com","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
		mockDeps.UserUseCase.EXPECT().SignUp(context.Background(), gomock.Any()).Return(usercase.SignUpResponse{}, nil)
		mockDeps.UserUseCase.EXPECT().AuthenticateUser(context.Background(), gomock.Any()).Return(usercase.AuthenticateUserResponse{}, errors.New("auth error"))

		err := handler.RegisterUser(c)
		assert.NoError(t, err) // Echoはエラーを返す代わりにHTTPレスポンスを返すので、エラー自体はnil
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
