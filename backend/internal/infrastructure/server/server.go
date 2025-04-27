package server

import (
	"example.com/infrahandson/config"
	adapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl"
	"example.com/infrahandson/internal/infrastructure/di"
	middleware "example.com/infrahandson/internal/infrastructure/gatewayImpl/middleware/echo"
	routes "example.com/infrahandson/internal/infrastructure/gatewayImpl/routes/echo"
	"example.com/infrahandson/internal/infrastructure/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ServerStart(cfg *config.Config, db *sqlx.DB) {
	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	// 依存関係の初期化
	dependencies := di.InitializeDependencies(cfg, db)

	// ミドルウェアの設定
	tokenService := adapterimpl.NewTokenServiceAdapter(adapterimpl.NewTokenServiceAdapterParams{
		SecretKey:     cfg.SecretKey,
		ExpireMinutes: int(cfg.TokenExpiry),
	})
	e.Use(middleware.CORS())
	e.Use(middleware.AuthMiddleware(tokenService))

	// ルーティングの設定
	routes.SetupRoutes(e, cfg, *dependencies.UserHandler, *dependencies.RoomHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
