package roomhandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. UseCase GetAllRooms がエラーを返す
func TestGetRomms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler, mockDeps, e := roomhandler.NewTestRoomHandler(ctrl)

	mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	// 1. 正常系
	t.Run("正常系", func(t *testing.T) {
		mockDeps.RoomUseCase.EXPECT().GetAllRooms(gomock.Any()).Return(
			[]*entity.Room{
				entity.NewRoom(entity.RoomParams{
					ID:   "room1",
					Name: "Test Room 1",
				}),
				entity.NewRoom(entity.RoomParams{
					ID:   "room2",
					Name: "Test Room 2",
				}),
			}, nil,
		)

		req := httptest.NewRequest("GET", "/rooms", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms")

		if err := handler.GetRooms(c); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// 2. UseCase GetAllRooms がエラーを返す
	t.Run("UseCase GetAllRooms がエラーを返す", func(t *testing.T) {
		mockDeps.RoomUseCase.EXPECT().GetAllRooms(gomock.Any()).Return(nil, assert.AnError)
		mockDeps.Logger.EXPECT().Error("Failed to get rooms", gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/rooms", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms")

		err := handler.GetRooms(c)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		httpErr, ok := err.(*echo.HTTPError)
		if !ok {
			t.Fatalf("Expected *echo.HTTPError, got %T", err)
		}

		assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
		assert.Equal(t, "Failed to get rooms", httpErr.Message)
	})

}
