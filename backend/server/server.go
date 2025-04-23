package server

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/infrastructure/adapter_impl"
	"example.com/webrtc-practice/internal/infrastructure/factory_impl"
	sqlite3 "example.com/webrtc-practice/internal/infrastructure/repository_impl/sqlite"
	"example.com/webrtc-practice/internal/interface/handler"
	"example.com/webrtc-practice/internal/usecase"
	"example.com/webrtc-practice/routes"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServerStart(cfg *config.Config, db *sqlx.DB) {
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                              // すべてのオリジンを許可
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},       // 許可するHTTPメソッド
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization}, // 許可するHTTPヘッダー
	}))

	// ユーザーハンドラの初期化
	userRepository := sqlite3.NewUserRepositoryImpl(&sqlite3.NewUserRepositoryImplParams{
		DB: db,
	})
	hasher := adapter_impl.NewHasherAdapter(adapter_impl.NewHasherAddapterParams{
		Cost: cfg.HashCost,
	})
	tokenService := adapter_impl.NewTokenServiceAdapter(adapter_impl.NewTokenServiceAdapterParams{
		SecretKey:     cfg.SecretKey,
		ExpireMinutes: int(cfg.TokenExpiry),
	})
	userIDFactory := factory_impl.NewUserIDFactory()
	userUsecase := usecase.NewUserUseCase(usecase.NewUserUseCaseParams{
		UserRepo:      userRepository,
		Hasher:        hasher,
		TokenSvc:      tokenService,
		UserIDFactory: userIDFactory,
	})
	userHandler := handler.NewUserHandler(handler.NewUserHandlerParams{
		UserUseCase:   userUsecase,
		UserIDFactory: userIDFactory,
	})

	// ルーティングの設定
	routes.SetupRoutes(e, cfg, *userHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
