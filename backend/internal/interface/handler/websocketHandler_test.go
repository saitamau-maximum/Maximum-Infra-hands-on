package handler_test

import (
	"context"
	"sync"
	"testing"

	"net/http"
	"net/http/httptest"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/usecase"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	mock_usecase "example.com/infrahandson/test/mocks/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestConnectToChatRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWsUseCase := mock_usecase.NewMockWebsocketUseCaseInterface(ctrl)
	mockWsUpGrader := mock_adapter.NewMockWebSocketUpgraderAdapter(ctrl)
	mockWsConnFactory := mock_factory.NewMockWebSocketConnectionFactory(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockRoomIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)

	WsHandler := handler.NewWebSocketHandler(handler.NewWebSocketHandlerParams{
		WsUseCase:        mockWsUseCase,
		WsUpgrader:       mockWsUpGrader,
		WsConnFactory:    mockWsConnFactory,
		UserIDFactory:    mockUserIDFactory,
		RoomIDFactory: mockRoomIDFactory,
		Logger:           mockLogger,
	})

	t.Run("Successful connection and message handling", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ws/test-room", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "test-user")
		c.SetParamNames("room_public_id")
		c.SetParamValues("test-room")

		mockConnRaw := mock_adapter.NewMockConnAdapter(ctrl)
		mockConn := mock_service.NewMockWebSocketConnection(ctrl)

		// goroutine終了のためのWaitGroup
		var wg sync.WaitGroup
		wg.Add(1)

		// メッセージのモック
		testMessage := entity.NewMessage(entity.MessageParams{
			ID: entity.MessageID("test-message"),
			UserID:   "test-user",
			RoomID:   "test-room",
			Content:  "testcontent",
			SentAt:   time.Now(),
		})

		mockLogger.EXPECT().
			Info(gomock.Any(), gomock.Any()).
			AnyTimes() // ロガーは何回呼ばれてもいい（呼ばれなくても怒らない）設定

		mockWsUpGrader.EXPECT().Upgrade(gomock.Any(), gomock.Any()).Return(mockConnRaw, nil)
		mockWsConnFactory.EXPECT().CreateWebSocketConnection(mockConnRaw).Return(mockConn, nil)
		mockWsUseCase.EXPECT().ConnectUserToRoom(context.Background(), gomock.Any()).Return(nil)
		time.Sleep(100 * time.Millisecond) // goroutine内の処理を待つためのスリープ
		// ここからゴルーチン内の処理
		mockConn.EXPECT().ReadMessage().Return(1, testMessage, nil)
		mockWsUseCase.EXPECT().SendMessage(context.Background(), gomock.Any()).Return(nil)

		// 少し待つ
		time.Sleep(100 * time.Millisecond)
		mockConn.EXPECT().ReadMessage().Return(0, nil, assert.AnError)
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any())
		mockWsUseCase.EXPECT().DisconnectUser(context.Background(), gomock.Any()).DoAndReturn(func(req usecase.DisconnectUserRequest) error {
			wg.Done() // goroutine終了の合図
			return nil
		})
		mockConn.EXPECT().Close().Return(nil)

		go func() {
			err := WsHandler.ConnectToChatRoom(c)
			assert.NoError(t, err)
		}()

		// goroutine終了まで待機
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// 正常終了
		case <-time.After(1 * time.Second):
			t.Fatal("Test timeout: goroutine did not finish")
		}

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Missing user ID", func(t *testing.T) {
		e := echo.New()
	
		req := httptest.NewRequest(http.MethodGet, "/ws/test-room", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any())
		err := WsHandler.ConnectToChatRoom(c)
		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
	})

	t.Run("Missing room public ID", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ws/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "test-user")

		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any())

		err := WsHandler.ConnectToChatRoom(c)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	})
}
