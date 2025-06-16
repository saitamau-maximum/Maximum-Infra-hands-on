package userhandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. user_id がコンテキストに存在しない
// 3. user_id が string 型でない
// 4. ファイル取得に失敗
// 5. UseCase.SaveUserIcon がエラーを返す
func TestSaveUserIcon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := userhandler.NewTestUserHandler(ctrl)

	userID := entity.UserID("mockUserID")
	iconContent := []byte("dummy image data")
	contentTypeVal := "image/png"
	body, contentType, _ := userhandler.NewMultipartRequestWithFile("icon", "icon.png", iconContent, contentTypeVal)

	t.Run("正常系: アイコン保存成功", func(t *testing.T) {
		mockDeps.UserUseCase.EXPECT().
			SaveUserIcon(gomock.Any(), gomock.Any(), userID).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/user/icon", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "mockUserID")

		handler.SaveUserIcon(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Icon saved successfully")
	})

	t.Run("異常系: user_id がコンテキストに存在しない", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user/icon", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// user_id を設定しない

		handler.SaveUserIcon(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Unauthorized")
	})

	t.Run("異常系: user_id が string 型でない", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user/icon", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", 12345)

		handler.SaveUserIcon(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid user ID")
	})

	t.Run("異常系: ファイル取得に失敗", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user/icon", strings.NewReader("invalid body"))
		req.Header.Set("Content-Type", "text/plain")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "mockUserID")

		handler.SaveUserIcon(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "error")
	})

	// io.Reader は1回しか読みだすことができない。正常系で消費したので作り直す
	body, contentType, _ = userhandler.NewMultipartRequestWithFile("icon", "icon.png", iconContent, contentTypeVal)

	t.Run("異常系: UseCase.SaveUserIcon がエラーを返す", func(t *testing.T) {
		mockDeps.UserUseCase.EXPECT().
			SaveUserIcon(gomock.Any(), gomock.Any(), userID).
			Return(errors.New("save failed"))

		req := httptest.NewRequest(http.MethodPost, "/user/icon", body)
		req.Header.Set("Content-Type", contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "mockUserID")

		handler.SaveUserIcon(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "error/saveIcon")
	})
}

// 1. 正常系
// 2. user_id が空
// 3. アイコンURL取得でエラー（アイコンが存在しない）
func TestGetUserIcon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := userhandler.NewTestUserHandler(ctrl)
	mockDeps.Logger.EXPECT().Info(gomock.Any()).AnyTimes() // 情報ログ出力は任意なので、全ての呼び出しに対して期待値を設定

	t.Run("正常系: アイコンURLをリダイレクトで返す", func(t *testing.T) {
		userID := "user-123"
		iconURL := "https://example.com/icons/user-123.webp"

		mockDeps.UserUseCase.EXPECT().
			GetUserIconPath(gomock.Any(), entity.UserID(userID)).
			Return(iconURL, nil)

		req := httptest.NewRequest(http.MethodGet, "/user/icon/"+userID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/icon/:user_id")
		c.SetParamNames("user_id")
		c.SetParamValues(userID)

		err := handler.GetUserIcon(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, iconURL, rec.Header().Get("Location"))
	})

	t.Run("異常系: user_id が空", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/icon/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/icon/:user_id")
		c.SetParamNames("user_id")
		c.SetParamValues("")

		err := handler.GetUserIcon(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "user_id is required")
	})

	t.Run("異常系: アイコンURL取得でエラー（アイコンが存在しない）", func(t *testing.T) {
		userID := "user-456"

		mockDeps.UserUseCase.EXPECT().
			GetUserIconPath(gomock.Any(), entity.UserID(userID)).
			Return("", errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/user/icon/"+userID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/icon/:user_id")
		c.SetParamNames("user_id")
		c.SetParamValues(userID)

		err := handler.GetUserIcon(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Icon not found")
	})
}
