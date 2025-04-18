package handler

import (
	"net/http"

	"example.com/webrtc-practice/internal/interface/factory"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type WebsocketHandler struct {
	UseCase                    usecase.IWebsocketUseCaseInterface
	WebsocketConnectionFactory factory.WebsocketConnectionFactory
}

func NewWebsocketHandler(
	usecase usecase.IWebsocketUseCaseInterface,
	factory factory.WebsocketConnectionFactory,
) WebsocketHandler {
	h := WebsocketHandler{
		UseCase:                    usecase,
		WebsocketConnectionFactory: factory,
	}

	// WebSocketメッセージ処理のゴルーチンを起動
	go h.HandleMessages()

	return h
}

func (h *WebsocketHandler) Register(g *echo.Group) {
	g.GET("/", h.HandleWebsocket)
}

// WebSocket接続
func (h *WebsocketHandler) HandleWebsocket(c echo.Context) error {
	// リクエストをコネクションにアップグレード
	conn, err := h.WebsocketConnectionFactory.NewConnection(c.Response().Writer, c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to upgrade connection"})
	}
	defer conn.Close()

	err = h.UseCase.RegisterClient(conn)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "client already registered"})
	}

	h.UseCase.ListenForMessages(conn)

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

// メッセージ処理の呼び出し
func (h *WebsocketHandler) HandleMessages() {
	h.UseCase.ProcessMessage()
}
