package server

import (
	"example.com/infrahandson/config"
	tokenadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/tokenServiceAdapterImpl/JWT"
	"example.com/infrahandson/internal/infrastructure/di"
	middleware "example.com/infrahandson/internal/infrastructure/gatewayImpl/middleware/echo"
	routes "example.com/infrahandson/internal/infrastructure/gatewayImpl/routes/echo"
	"example.com/infrahandson/internal/infrastructure/validator"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// ServerStart initializes the Echo server, sets up dependencies, middleware, and routes,
// and returns the Echo instance along with the database connection.
//
// Parameters:
// - cfg: A pointer to the configuration object containing server settings.
//
// Returns:
// - *echo.Echo: The initialized Echo server instance.
// - *sqlx.DB: The database connection initialized with the provided configuration.
func ServerStart(cfg *config.Config) (*echo.Echo, *sqlx.DB, *memcache.Client) {
	e := echo.New()
	e.Static("/images/icons", "./images/icons")
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
		dependencies.Handler,
	)

	return e, dependencies.DB, dependencies.Cache
}
