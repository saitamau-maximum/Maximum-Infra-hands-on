package di

import (
	"example.com/webrtc-practice/config"
	adapterimpl "example.com/webrtc-practice/internal/infrastructure/adapter_impl"
	fmtloggerimpl "example.com/webrtc-practice/internal/infrastructure/adapter_impl/logger_adapter_impl/fmt_logger"
	factoryimpl "example.com/webrtc-practice/internal/infrastructure/factory_impl"
	sqliteroomrepoimpl "example.com/webrtc-practice/internal/infrastructure/repository_impl/room_repository_impl/sqlite"
	sqliteuserrepoimpl "example.com/webrtc-practice/internal/infrastructure/repository_impl/user_repository_impl/sqlite"
	"example.com/webrtc-practice/internal/interface/handler"
	"example.com/webrtc-practice/internal/usecase"

	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	UserHandler *handler.UserHandler
	RoomHandler *handler.RoomHandler
}

func InitializeDependencies(cfg *config.Config, db *sqlx.DB) *Dependencies {
	// Loggerの設定
	logger := fmtloggerimpl.NewFmtLogger()

	// IDFactoryの初期化
	userIDFactory := factoryimpl.NewUserIDFactory()
	roomIDFactory := factoryimpl.NewRoomIDFactory()

	// Repositoryの初期化
	userRepository := sqliteuserrepoimpl.NewUserRepositoryImpl(&sqliteuserrepoimpl.NewUserRepositoryImplParams{
		DB: db,
	})
	roomRepository := sqliteroomrepoimpl.NewRoomRepositoryImpl(&sqliteroomrepoimpl.NewRoomRepositoryImplParams{
		DB: db,
	})

	// AdapterとServiceの初期化
	hasher := adapterimpl.NewHasherAdapter(adapterimpl.NewHasherAddapterParams{
		Cost: cfg.HashCost,
	})
	tokenService := adapterimpl.NewTokenServiceAdapter(adapterimpl.NewTokenServiceAdapterParams{
		SecretKey:     cfg.SecretKey,
		ExpireMinutes: int(cfg.TokenExpiry),
	})

	// UseCaseの初期化
	userUseCase := usecase.NewUserUseCase(usecase.NewUserUseCaseParams{
		UserRepo:      userRepository,
		Hasher:        hasher,
		TokenSvc:      tokenService,
		UserIDFactory: userIDFactory,
	})
	roomUseCase := usecase.NewRoomUseCase(usecase.NewRoomUseCaseParams{
		RoomRepo:      roomRepository,
		UserRepo:      userRepository,
		RoomIDFactory: roomIDFactory,
	})

	// Handlerの初期化
	userHandler := handler.NewUserHandler(handler.NewUserHandlerParams{
		UserUseCase:   userUseCase,
		UserIDFactory: userIDFactory,
		Logger:        logger,
	})
	roomHandler := handler.NewRoomHandler(handler.NewRoomHandlerParams{
		RoomUseCase:   roomUseCase,
		UserIDFactory: userIDFactory,
		RoomIDFactory: roomIDFactory,
		Logger:        logger,
	})

	return &Dependencies{
		UserHandler: userHandler,
		RoomHandler: roomHandler,
	}
}
