package handler_test

import (
	"sync"
	"testing"

	"net/http"
	"net/http/httptest"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/interface/handler"
	"example.com/webrtc-practice/internal/usecase"
	mock_service "example.com/webrtc-practice/mocks/domain/service"
	mock_adapter "example.com/webrtc-practice/mocks/interface/adapter"
	mock_factory "example.com/webrtc-practice/mocks/interface/factory"
	mock_usecase "example.com/webrtc-practice/mocks/usecase"
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
	mockRoomPubIDFactory := mock_factory.NewMockRoomIDFactory(ctrl)

	WsHandler := handler.NewWebSocketHandler(handler.WebSocketHandlerParams{
		WsUseCase:        mockWsUseCase,
		WsUpgrader:       mockWsUpGrader,
		WsConnFactory:    mockWsConnFactory,
		UserIDFactory:    mockUserIDFactory,
		RoomPubIDFactory: mockRoomPubIDFactory,
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
			ID:      1,
			PublicID: entity.MessagePublicID("test-message"),
			UserID:  1,
			RoomID:  1,
			Content: "testcontent",
			SentAt:  time.Now(),
		})



		mockWsUpGrader.EXPECT().Upgrade(gomock.Any(), gomock.Any()).Return(mockConnRaw, nil)
		mockWsConnFactory.EXPECT().CreateWebSocketConnection(mockConnRaw).Return(mockConn, nil)
		mockUserIDFactory.EXPECT().FromString("test-user").Return(entity.UserPublicID("test-user"), nil)
		mockRoomPubIDFactory.EXPECT().FromString("test-room").Return(entity.RoomPublicID("test-room"), nil)
		mockWsUseCase.EXPECT().ConnectUserToRoom(gomock.Any()).Return(nil)
		time.Sleep(100 * time.Millisecond) // goroutine内の処理を待つためのスリープ
		// ここからゴルーチン内の処理
		mockConn.EXPECT().ReadMessage().Return(1, testMessage, nil)
		mockWsUseCase.EXPECT().SendMessage(gomock.Any()).Return(nil)

		// 少し待つ
		time.Sleep(100 * time.Millisecond)
		mockConn.EXPECT().ReadMessage().Return(0, nil, assert.AnError)
		mockUserIDFactory.EXPECT().FromString("test-user").Return(entity.UserPublicID("test-user"), nil)
		mockWsUseCase.EXPECT().DisconnectUser(gomock.Any()).DoAndReturn(func(req usecase.DisconnectUserRequest) error {
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

		err := WsHandler.ConnectToChatRoom(c)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	})
}
