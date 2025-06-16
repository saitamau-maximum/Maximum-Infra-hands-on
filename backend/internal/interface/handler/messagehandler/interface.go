package messagehandler

import "github.com/labstack/echo/v4"

type MessageHandlerInterface interface {
	// GetMessageHistoryInRoom は指定されたルームのメッセージ履歴を取得する
	GetRoomMessage(c echo.Context) error
}
