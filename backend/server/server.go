package server

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/infrastructure/adapter_impl"
	"example.com/webrtc-practice/internal/infrastructure/factory_impl"
	sqlite3 "example.com/webrtc-practice/internal/infrastructure/repository_impl/sqlite"
	"example.com/webrtc-practice/internal/infrastructure/validator"
	"example.com/webrtc-practice/internal/interface/handler"
	"example.com/webrtc-practice/internal/usecase"
	"example.com/webrtc-practice/routes"
	"example.com/webrtc-practice/server/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ServerStart(cfg *config.Config, db *sqlx.DB) {
	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	tokenService := adapter_impl.NewTokenServiceAdapter(adapter_impl.NewTokenServiceAdapterParams{
		SecretKey:     cfg.SecretKey,
		ExpireMinutes: int(cfg.TokenExpiry),
	})

	// ミドルウェアの設定
	e.Use(middleware.CORS())
	e.Use(middleware.AuthMiddleware(tokenService))

	// ユーザーハンドラの初期化
	userRepository := sqlite3.NewUserRepositoryImpl(&sqlite3.NewUserRepositoryImplParams{
		DB: db,
	})
	hasher := adapter_impl.NewHasherAdapter(adapter_impl.NewHasherAddapterParams{
		Cost: cfg.HashCost,
	})

	userIDFactory := factory_impl.NewUserIDFactory()
	userUseCase := usecase.NewUserUseCase(usecase.NewUserUseCaseParams{
		UserRepo:      userRepository,
		Hasher:        hasher,
		TokenSvc:      tokenService,
		UserIDFactory: userIDFactory,
	})
	userHandler := handler.NewUserHandler(handler.NewUserHandlerParams{
		UserUseCase:   userUseCase,
		UserIDFactory: userIDFactory,
	})

	roomIDFactory := factory_impl.NewRoomIDFactory()
	// roomRepositoryimpl実装
	roomRepository := sqlite3.NewRoomRepositoryImpl(&sqlite3.NewRoomRepositoryImplParams{
		DB: db,
	})
	roomUseCase := usecase.NewRoomUseCase(usecase.NewRoomUseCaseParams{
		RoomRepo: roomRepository,
		UserRepo: userRepository,
		RoomIDFactory: roomIDFactory,
	})
	roomHandler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase: roomUseCase,
		UserIDFactory: userIDFactory,
		RoomIDFactory: roomIDFactory,
	})

	// ルーティングの設定
	routes.SetupRoutes(e, cfg, *userHandler, *roomHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
