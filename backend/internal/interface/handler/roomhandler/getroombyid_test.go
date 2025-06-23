package roomhandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"example.com/infrahandson/internal/usecase/roomcase"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

// 1. 正常系
// 2. パラメータ room_id が空
// 3. UseCase GetRoomByID がエラーを返す
func TestGetRoomByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler, mockDeps, e := roomhandler.NewTestRoomHandler(ctrl)

	mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	// 1. 正常系
	t.Run("正常系", func(t *testing.T) {
		roomID := "test-room-id"
		mockDeps.RoomUseCase.EXPECT().GetRoomByID(gomock.Any(), gomock.Any()).Return(
			roomcase.GetRoomByIDResponse{
				Room: entity.NewRoom(entity.RoomParams{
					ID:      entity.RoomID(roomID),
					Name:    "Test Room",
					Members: []entity.UserID{entity.UserID("member1"), entity.UserID("member2")},
				}),
			}, nil,
		)

		req := httptest.NewRequest("GET", "/rooms/"+roomID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id")
		c.SetParamNames("room_id")
		c.SetParamValues(roomID)

		if err := handler.GetRoomByID(c); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if rec.Code != 200 {
			t.Errorf("Expected status code 200, got %d", rec.Code)
		}
	})

	// 2. パラメータ room_id が空
	t.Run("room_id is empty", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rooms/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id")
		c.SetParamNames("room_id")
		c.SetParamValues("")

		mockDeps.Logger.EXPECT().Error("Room ID is missing").Times(1)

		err := handler.GetRoomByID(c)
		if err == nil {
			t.Error("Expected error for empty room_id, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 BadRequest, got %v", err)
		}
	})

	// 3. UseCase GetRoomByID がエラーを返す
	t.Run("usecase returns error", func(t *testing.T) {
		roomID := "test-room-id"
		mockDeps.RoomUseCase.EXPECT().GetRoomByID(gomock.Any(), gomock.Any()).Return(
			roomcase.GetRoomByIDResponse{}, echo.NewHTTPError(http.StatusNotFound, "room not found"),
		)
		mockDeps.Logger.EXPECT().Error("Failed to get room", gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/rooms/"+roomID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:room_id")
		c.SetParamNames("room_id")
		c.SetParamValues(roomID)

		err := handler.GetRoomByID(c)
		if err == nil {
			t.Error("Expected error from usecase, got nil")
		}
		he, ok := err.(*echo.HTTPError)
		if !ok || he.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500 InternalServerError, got %v", err)
		}
	})
}
