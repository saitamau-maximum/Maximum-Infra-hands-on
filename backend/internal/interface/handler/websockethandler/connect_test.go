package websockethandler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/handler/websockethandler"
	"example.com/infrahandson/internal/usecase/websocketcase"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestConnectToChatRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockDeps, e := websockethandler.NewTestWebsocketHandler(ctrl)

	t.Run("Successful connection and message handling", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ws/test-room", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "test-user")
		c.SetParamNames("room_id")
		c.SetParamValues("test-room")

		mockConnRaw := mock_adapter.NewMockConnAdapter(ctrl)
		mockConn := mock_service.NewMockWebSocketConnection(ctrl)

		var wg sync.WaitGroup
		wg.Add(1)

		testMessage := entity.NewMessage(entity.MessageParams{
			ID:      entity.MessageID("test-message"),
			UserID:  "test-user",
			RoomID:  "test-room",
			Content: "testcontent",
			SentAt:  time.Now(),
		})

		mockDeps.Logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

		mockDeps.WsUpgrader.EXPECT().Upgrade(gomock.Any(), gomock.Any()).Return(mockConnRaw, nil)
		mockDeps.WsConnFactory.EXPECT().CreateWebSocketConnection(mockConnRaw).Return(mockConn, nil)
		mockDeps.WsUseCase.EXPECT().ConnectUserToRoom(gomock.Any(), gomock.Any()).Return(nil)

		gomock.InOrder(
			mockConn.EXPECT().ReadMessage().Return(testMessage, nil),
			mockDeps.WsUseCase.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil),
			mockConn.EXPECT().ReadMessage().Return(nil, assert.AnError),
			mockDeps.Logger.EXPECT().Warn(gomock.Any(), gomock.Any()),
			mockDeps.WsUseCase.EXPECT().DisconnectUser(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ websocketcase.DisconnectUserRequest) error {
				wg.Done()
				return nil
			}),
			mockConn.EXPECT().Close().Return(nil),
		)

		go func() {
			err := handler.ConnectToChatRoom(c)
			assert.NoError(t, err)
		}()

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// OK
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

		mockDeps.Logger.EXPECT().Warn(gomock.Any(), gomock.Any())
		err := handler.ConnectToChatRoom(c)
		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
	})

	t.Run("Missing room public ID", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ws/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "test-user")

		mockDeps.Logger.EXPECT().Warn(gomock.Any(), gomock.Any())

		err := handler.ConnectToChatRoom(c)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	})
}
