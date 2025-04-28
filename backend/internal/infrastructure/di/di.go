package di

import (
	"example.com/infrahandson/config"
	bcryptadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/hasherAdapterImpl/bcrypt"
	fmtloggerimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/loggerAdapterImpl/fmtLogger"
	tokenadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/tokenServiceAdapterImpl/JWT"
	gorillawebsocketupgraderImpl "example.com/infrahandson/internal/infrastructure/adapterImpl/upgraderAdapterImpl/gorillawebsocket"
	factoryimpl "example.com/infrahandson/internal/infrastructure/factoryImpl"
	sqlitemsgrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/sqlite"
	sqliteroomrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/sqlite"
	sqliteuserrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/sqlite"
	inmemorywsclientrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/websocketClientRepositoryImpl/InMemory"
	inmemorywsmanagerimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/websocketManagerImpl/InMemory"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/usecase"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	UserHandler *handler.UserHandler
	RoomHandler *handler.RoomHandler
	WsHandler   *handler.WebSocketHandler
}

func InitializeDependencies(cfg *config.Config, db *sqlx.DB) *Dependencies {
	// Loggerの設定
	logger := fmtloggerimpl.NewFmtLogger()

	// Factoryの初期化
	userIDFactory := factoryimpl.NewUserIDFactory()
	roomIDFactory := factoryimpl.NewRoomIDFactory()
	MsgIDFactory := factoryimpl.NewMessageIDFactory()
	clientDFactory := factoryimpl.NewWsClientIDFactory()
	upgrader := gorillawebsocketupgraderImpl.NewGorillaWebSocketUpgrader()
	wsConnFactory := factoryimpl.NewWebSocketConnectionFactoryImpl()

	// Repositoryの初期化
	userRepository := sqliteuserrepoimpl.NewUserRepositoryImpl(&sqliteuserrepoimpl.NewUserRepositoryImplParams{
		DB: db,
	})
	roomRepository := sqliteroomrepoimpl.NewRoomRepositoryImpl(&sqliteroomrepoimpl.NewRoomRepositoryImplParams{
		DB: db,
	})
	msgRepository := sqlitemsgrepoimpl.NewMessageRepositoryImpl(&sqlitemsgrepoimpl.NewMessageRepositoryImplParams{
		DB: db,
	})
	wsClientRepository := inmemorywsclientrepoimpl.NewInMemoryWebsocketClientRepository(inmemorywsclientrepoimpl.NewInMemoryWebsocketClientRepositoryParams{})

	// Serviceの初期化
	wsManager := inmemorywsmanagerimpl.NewInMemoryWebSocketManager()

	// AdapterとServiceの初期化
	hasher := bcryptadapterimpl.NewHasherAdapter(bcryptadapterimpl.NewHasherAddapterParams{
		Cost: cfg.HashCost,
	})
	tokenService := tokenadapterimpl.NewTokenServiceAdapter(tokenadapterimpl.NewTokenServiceAdapterParams{
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
	wsUseCase := usecase.NewWebsocketUseCase(usecase.NewWebsocketUseCaseParams{
		UserRepo:     userRepository,
		RoomRepo:     roomRepository,
		MsgRepo:      msgRepository,
		WsClientRepo: wsClientRepository,
		WebsocketManager: wsManager,
		MsgIDFactory:    MsgIDFactory,
		ClientIDFactory: clientDFactory,
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
	wsHandler := handler.NewWebSocketHandler(handler.NewWebSocketHandlerParams{
		WsUseCase:  wsUseCase,
		WsUpgrader: upgrader,
		WsConnFactory:wsConnFactory,
		UserIDFactory: userIDFactory,
		RoomIDFactory: roomIDFactory,
		Logger:        logger,
	})

	return &Dependencies{
		UserHandler: userHandler,
		RoomHandler: roomHandler,
		WsHandler:   wsHandler,
	}
}
