package server

import (
	"example.com/infrahandson/config"
	tokenadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/tokenServiceAdapterImpl/JWT"
	"example.com/infrahandson/internal/infrastructure/di"
	middleware "example.com/infrahandson/internal/infrastructure/gatewayImpl/middleware/echo"
	routes "example.com/infrahandson/internal/infrastructure/gatewayImpl/routes/echo"
	"example.com/infrahandson/internal/infrastructure/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ServerStart(cfg *config.Config) (*echo.Echo, *sqlx.DB) {
	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	// 依存関係の初期化
	dependencies := di.InitializeDependencies(cfg)

	// ミドルウェアの設定
	tokenService := tokenadapterimpl.NewTokenServiceAdapter(tokenadapterimpl.NewTokenServiceAdapterParams{
		SecretKey:     cfg.SecretKey,
		ExpireMinutes: int(cfg.TokenExpiry),
	})
	e.Use(middleware.CORS(cfg.CORSOrigin))

	// ルーティングの設定
	routes.SetupRoutes(
		e,
		cfg,
		middleware.AuthMiddleware(tokenService),
		*dependencies.UserHandler,
		*dependencies.RoomHandler,
		*dependencies.WsHandler,
		*dependencies.MsgHandler,
	)

	return e, dependencies.DB
}
