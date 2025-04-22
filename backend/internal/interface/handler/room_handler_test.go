package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/infrastructure/validator"
	"example.com/webrtc-practice/internal/interface/handler"
	"example.com/webrtc-practice/internal/usecase"
	mock_factory "example.com/webrtc-practice/mocks/interface/factory"
	mock_usecase "example.com/webrtc-practice/mocks/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomUseCase := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
	})

	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	req := httptest.NewRequest(http.MethodPost, "/rooms", strings.NewReader(`{"name":"testRoom"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "mockUserID")

	mockRoomUseCase.EXPECT().CreateRoom(gomock.Any()).Return(usecase.CreateRoomResponse{
		Room: entity.NewRoom(entity.RoomParams{
			ID:       1,
			PublicID: "mockRoomID",
			Name:     "testRoom",
			Members:  []entity.UserPublicID{"mockUserID"},
		}),
	}, nil)
	mockUserIDFactory.EXPECT().FromString("mockUserID").Return(entity.UserPublicID("mockUserID"))
	mockRoomUseCase.EXPECT().JoinRoom(gomock.Any()).Return(nil)

	if assert.NoError(t, handler.CreateRoom(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "mockRoomID")
	}
}

func TestJoinRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomUseCase := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/rooms/mockRoomID/join", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "mockUserID")
	c.SetParamNames("public_id")
	c.SetParamValues("mockRoomID")

	mockUserIDFactory.EXPECT().FromString("mockUserID").Return(entity.UserPublicID("mockUserID"))
	mockRoomIDFactory.EXPECT().FromString("mockRoomID").Return(entity.RoomPublicID("mockRoomID"))
	mockRoomUseCase.EXPECT().JoinRoom(gomock.Any()).Return(nil)

	if assert.NoError(t, handler.JoinRoom(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Joined room successfully")
	}
}

func TestLeaveRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomUseCase := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/rooms/mockRoomID/leave", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "mockUserID")
	c.SetParamNames("public_id")
	c.SetParamValues("mockRoomID")

	mockUserIDFactory.EXPECT().FromString("mockUserID").Return(entity.UserPublicID("mockUserID"))
	mockRoomIDFactory.EXPECT().FromString("mockRoomID").Return(entity.RoomPublicID("mockRoomID"))
	mockRoomUseCase.EXPECT().LeaveRoom(gomock.Any()).Return(nil)

	if assert.NoError(t, handler.LeaveRoom(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Left room successfully")
	}
}

func TestGetRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomUseCase := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:      mockRoomUseCase,
		UserIDFactory:    mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/rooms/mockRoomID", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("public_id")
	c.SetParamValues("mockRoomID")

	mockRoomIDFactory.EXPECT().FromString("mockRoomID").Return(entity.RoomPublicID("mockRoomID"))
	mockRoomUseCase.EXPECT().GetRoomByPublicID(gomock.Any()).Return(usecase.GetRoomByPublicIDResponse{
		Room: entity.NewRoom(entity.RoomParams{
			ID:       1,
			PublicID: "mockRoomID",
			Name:     "testRoom",
			Members:  []entity.UserPublicID{"user1", "user2"},
		}),
	}, nil)

	if assert.NoError(t, handler.GetRoom(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "mockRoomID")
		assert.Contains(t, rec.Body.String(), "testRoom")
	}
}

func TestGetRooms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomUseCase := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:      mockRoomUseCase,
		UserIDFactory:    mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/rooms", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRoomUseCase.EXPECT().GetAllRooms().Return([]*entity.Room{
		entity.NewRoom(entity.RoomParams{
			ID:       1,
			PublicID: "room1",
			Name:     "Room 1",
			Members:  []entity.UserPublicID{"user1", "user2"},
		}),
		entity.NewRoom(entity.RoomParams{
			ID:       2,
			PublicID: "room2",
			Name:     "Room 2",
			Members:  []entity.UserPublicID{"user3", "user4"},
		}),
	}, nil)

	if assert.NoError(t, handler.GetRooms(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "room1")
		assert.Contains(t, rec.Body.String(), "Room 1")
		assert.Contains(t, rec.Body.String(), "room2")
		assert.Contains(t, rec.Body.String(), "Room 2")
	}
}
