package middleware

import (
	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

func CORS(CORSOrigin string) echo.MiddlewareFunc {
	return emw.CORSWithConfig(emw.CORSConfig{
		AllowOrigins: []string{CORSOrigin},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	})
}
