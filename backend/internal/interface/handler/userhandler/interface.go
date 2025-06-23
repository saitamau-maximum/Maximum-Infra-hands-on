package userhandler

import "github.com/labstack/echo/v4"

type UserHandlerInterface interface {
	// RegisterUser はユーザー登録を行う
	RegisterUser(c echo.Context) error

	// Login はユーザーのログインを行い、クッキーをセットする
	Login(c echo.Context) error

	// GetMe は現在のユーザー情報を取得する
	GetMe(c echo.Context) error

	// Logout はユーザーのログアウトを行い、クッキーをクリアする
	Logout(c echo.Context) error

	// SaveUserIcon はユーザーのアイコン画像を保存する
	SaveUserIcon(c echo.Context) error

	// GetUserIcon はユーザーのアイコン画像を取得する
	GetUserIcon(c echo.Context) error
}
