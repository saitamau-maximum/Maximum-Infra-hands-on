package middleware

import (
	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

func CORS() echo.MiddlewareFunc {
	return emw.CORSWithConfig(emw.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},// TODO: 環境変数化
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	})
}
