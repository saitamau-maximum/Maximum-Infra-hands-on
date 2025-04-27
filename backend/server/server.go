package server

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/infrastructure/adapter_impl"
	fmtloggerimpl "example.com/webrtc-practice/internal/infrastructure/adapter_impl/logger_adapter_impl/fmt_logger"
	"example.com/webrtc-practice/internal/infrastructure/factory_impl"
	sqliteroomrepoimpl "example.com/webrtc-practice/internal/infrastructure/repository_impl/room_repository_impl/sqlite"
	sqliteuserrepoimpl "example.com/webrtc-practice/internal/infrastructure/repository_impl/user_repository_impl/sqlite"
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

	tokenService := adapterimpl.NewTokenServiceAdapter(adapterimpl.NewTokenServiceAdapterParams{
		SecretKey:     cfg.SecretKey,
		ExpireMinutes: int(cfg.TokenExpiry),
	})

	// ミドルウェアの設定
	e.Use(middleware.CORS())
	e.Use(middleware.AuthMiddleware(tokenService))

	// Loggerの設定
	logger := fmtloggerimpl.NewFmtLogger()

	// ユーザーハンドラの初期化
	userRepository := sqliteuserrepoimpl.NewUserRepositoryImpl(&sqliteuserrepoimpl.NewUserRepositoryImplParams{
		DB: db,
	})
	hasher := adapterimpl.NewHasherAdapter(adapterimpl.NewHasherAddapterParams{
		Cost: cfg.HashCost,
	})

	userIDFactory := factoryimpl.NewUserIDFactory()
	userUseCase := usecase.NewUserUseCase(usecase.NewUserUseCaseParams{
		UserRepo:      userRepository,
		Hasher:        hasher,
		TokenSvc:      tokenService,
		UserIDFactory: userIDFactory,
	})
	userHandler := handler.NewUserHandler(handler.NewUserHandlerParams{
		UserUseCase:   userUseCase,
		UserIDFactory: userIDFactory,
		Logger:        logger,
	})

	roomIDFactory := factoryimpl.NewRoomIDFactory()
	// roomRepositoryimpl実装
	roomRepository := sqliteroomrepoimpl.NewRoomRepositoryImpl(&sqliteroomrepoimpl.NewRoomRepositoryImplParams{
		DB: db,
	})
	roomUseCase := usecase.NewRoomUseCase(usecase.NewRoomUseCaseParams{
		RoomRepo:      roomRepository,
		UserRepo:      userRepository,
		RoomIDFactory: roomIDFactory,
	})
	roomHandler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   roomUseCase,
		UserIDFactory: userIDFactory,
		RoomIDFactory: roomIDFactory,
		Logger:        logger,
	})

	// ルーティングの設定
	routes.SetupRoutes(e, cfg, *userHandler, *roomHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
