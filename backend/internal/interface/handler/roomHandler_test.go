package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/infrastructure/validator"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/usecase"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	mock_usecase "example.com/infrahandson/test/mocks/usecase"

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
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:        mockLogger,
	})

	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	req := httptest.NewRequest(http.MethodPost, "/rooms", strings.NewReader(`{"name":"testRoom"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "mockUserID")
	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes() // ロガーは何回呼ばれてもいい（呼ばれなくても怒らない）設定

	mockRoomUseCase.EXPECT().CreateRoom(context.Background(), gomock.Any()).Return(usecase.CreateRoomResponse{
		Room: entity.NewRoom(entity.RoomParams{
			ID:      "mockRoomID",
			Name:    "testRoom",
			Members: []entity.UserID{"mockUserID"},
		}),
	}, nil)
	mockRoomUseCase.EXPECT().JoinRoom(context.Background(), gomock.Any()).Return(nil)

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
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:        mockLogger,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/rooms/mockRoomID/join", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "mockUserID")
	c.SetParamNames("public_id")
	c.SetParamValues("mockRoomID")
	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes() // ロガーは何回呼ばれてもいい（呼ばれなくても怒らない）設定

	mockRoomUseCase.EXPECT().JoinRoom(context.Background(), gomock.Any()).Return(nil)

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
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:        mockLogger,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/rooms/mockRoomID/leave", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "mockUserID")
	c.SetParamNames("public_id")
	c.SetParamValues("mockRoomID")

	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes() // ロガーは何回呼ばれてもいい（呼ばれなくても怒らない）設定

	mockRoomUseCase.EXPECT().LeaveRoom(context.Background(), gomock.Any()).Return(nil)

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
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:        mockLogger,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/rooms/mockRoomID", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("public_id")
	c.SetParamValues("mockRoomID")

	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes() // ロガーは何回呼ばれてもいい（呼ばれなくても怒らない）設定

	mockRoomUseCase.EXPECT().GetRoomByID(context.Background(), gomock.Any()).Return(usecase.GetRoomByIDResponse{
		Room: entity.NewRoom(entity.RoomParams{
			ID:      "mockRoomID",
			Name:    "testRoom",
			Members: []entity.UserID{"user1", "user2"},
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
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	handler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   mockRoomUseCase,
		UserIDFactory: mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:        mockLogger,
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/rooms", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLogger.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes() // ロガーは何回呼ばれてもいい（呼ばれなくても怒らない）設定

	mockRoomUseCase.EXPECT().GetAllRooms(context.Background()).Return([]*entity.Room{
		entity.NewRoom(entity.RoomParams{
			ID:      "room1",
			Name:    "Room 1",
			Members: []entity.UserID{"user1", "user2"},
		}),
		entity.NewRoom(entity.RoomParams{
			ID:      "room2",
			Name:    "Room 2",
			Members: []entity.UserID{"user3", "user4"},
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
