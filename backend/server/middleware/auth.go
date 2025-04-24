package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TokenService interface {
	ValidateToken(token string) (string, error) // トークンからUserID取得
}

func AuthMiddleware(tokenService TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return next(c)
			}

			userID, err := tokenService.ValidateToken(cookie.Value)
			if err != nil {
				// 検証失敗ならUnauthorized
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			// Contextにuser_idを保存
			c.Set("user_id", userID)
			return next(c)
		}
	}
}
