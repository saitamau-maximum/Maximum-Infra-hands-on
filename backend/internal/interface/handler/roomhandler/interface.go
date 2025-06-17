package roomhandler

import "github.com/labstack/echo/v4"

type RoomHandlerInterface interface {
	// CreateRoom は新しいルームを作成するハンドラーです。
	// 名前からルームIDを生成し、ルームを作成します。
	CreateRoom(c echo.Context) error

	// JoinRoom は指定されたルームに参加するハンドラーです。
	// ルームIDをパラメータから取得し、参加処理を行います。
	JoinRoom(c echo.Context) error

	// LeaveRoom は指定されたルームから退出するハンドラーです。
	// ルームIDをパラメータから取得し、退出処理を行います。
	LeaveRoom(c echo.Context) error

	// GetRoom は指定されたルームの情報を取得するハンドラーです。
	// ルームIDをパラメータから取得し、ルームの詳細を返します。
	GetRoomByID(c echo.Context) error

	// GetRooms は部屋の情報を一括取得するハンドラーです。
	GetRooms(c echo.Context) error
}
