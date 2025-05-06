package di

import (
	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/domain/repository"
	bcryptadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/hasherAdapterImpl/bcrypt"
	fmtloggerimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/loggerAdapterImpl/fmtLogger"
	tokenadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/tokenServiceAdapterImpl/JWT"
	gorillawebsocketupgraderImpl "example.com/infrahandson/internal/infrastructure/adapterImpl/upgraderAdapterImpl/gorillawebsocket"
	factoryimpl "example.com/infrahandson/internal/infrastructure/factoryImpl"
	mysqlgatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/mysql"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	mysqlmsgrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/mysql"
	sqlitemsgrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/sqlite"
	mysqlroomrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/mysql"
	sqliteroomrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/sqlite"
	mysqluserrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/mysql"
	sqliteuserrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/sqlite"
	inmemorywsclientrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/websocketClientRepositoryImpl/InMemory"
	inmemorywsmanagerimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/websocketManagerImpl/InMemory"
	"example.com/infrahandson/internal/interface/gateway"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/usecase"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DB          *sqlx.DB
	UserHandler *handler.UserHandler
	RoomHandler *handler.RoomHandler
	WsHandler   *handler.WebSocketHandler
	MsgHandler  *handler.MessageHandler
}

func InitializeDependencies(cfg *config.Config) *Dependencies {
	// Loggerの設定
	logger := fmtloggerimpl.NewFmtLogger()

	// DBの初期化
	var initializer gateway.DBInitializer
	// initializerの作成
	if cfg.MySQLDSN != nil {
		// MySQL用のDSNが設定されている場合、MySQL用の初期化処理を行う
		initializer = mysqlgatewayimpl.NewMySQLInitializer(&mysqlgatewayimpl.NewMySQLInitializerParams{
			DSN:            cfg.MySQLDSN,
			MigrationsPath: "./internal/infrastructure/gatewayImpl/db/mysql/migrations",
		})
	} else {
		// SQLite用の初期化処理を行う
		initializer = sqlitegatewayimpl.NewSQLiteInitializer(&sqlitegatewayimpl.NewSQLiteInitializerParams{
			Path:           cfg.DBPath,
			MigrationsPath: "./internal/infrastructure/gatewayImpl/db/sqlite/migrations",
		})
	}
	db, err := initializer.Init()
	if err != nil {
		panic("failed to initialize database: " + err.Error())
	}

	// スキーマの初期化
	if err := initializer.InitSchema(db); err != nil {
		panic("failed to initialize schema: " + err.Error())
	}

	// Factoryの初期化
	userIDFactory := factoryimpl.NewUserIDFactory()
	roomIDFactory := factoryimpl.NewRoomIDFactory()
	MsgIDFactory := factoryimpl.NewMessageIDFactory()
	clientDFactory := factoryimpl.NewWsClientIDFactory()
	upgrader := gorillawebsocketupgraderImpl.NewGorillaWebSocketUpgrader()
	wsConnFactory := factoryimpl.NewWebSocketConnectionFactoryImpl()

	var userRepository repository.UserRepository
	var roomRepository repository.RoomRepository
	var msgRepository repository.MessageRepository

	// Repositoryの初期化
	if cfg.MySQLDSN != nil {
		userRepository = mysqluserrepoimpl.NewUserRepositoryImpl(&mysqluserrepoimpl.NewUserRepositoryImplParams{DB: db})
		roomRepository = mysqlroomrepoimpl.NewRoomRepositoryImpl(&mysqlroomrepoimpl.NewRoomRepositoryImplParams{DB: db})
		msgRepository = mysqlmsgrepoimpl.NewMessageRepositoryImpl(&mysqlmsgrepoimpl.NewMessageRepositoryImplParams{DB: db})
	} else {
		userRepository = sqliteuserrepoimpl.NewUserRepositoryImpl(&sqliteuserrepoimpl.NewUserRepositoryImplParams{DB: db})
		roomRepository = sqliteroomrepoimpl.NewRoomRepositoryImpl(&sqliteroomrepoimpl.NewRoomRepositoryImplParams{DB: db})
		msgRepository = sqlitemsgrepoimpl.NewMessageRepositoryImpl(&sqlitemsgrepoimpl.NewMessageRepositoryImplParams{DB: db})
	}

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
		UserRepo:         userRepository,
		RoomRepo:         roomRepository,
		MsgRepo:          msgRepository,
		WsClientRepo:     wsClientRepository,
		WebsocketManager: wsManager,
		MsgIDFactory:     MsgIDFactory,
		ClientIDFactory:  clientDFactory,
	})
	msgUseCase := usecase.NewMessageUseCase(usecase.NewMessageUseCaseParams{
		MsgRepo:  msgRepository,
		RoomRepo: roomRepository,
		UserRepo: userRepository,
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
		WsUseCase:     wsUseCase,
		WsUpgrader:    upgrader,
		WsConnFactory: wsConnFactory,
		UserIDFactory: userIDFactory,
		RoomIDFactory: roomIDFactory,
		Logger:        logger,
	})
	msgHansler := handler.NewMessageHandler(handler.NewMessageHandlerParams{
		MsgUseCase: msgUseCase,
		Logger:     logger,
	})

	return &Dependencies{
		DB:          db,
		UserHandler: userHandler,
		RoomHandler: roomHandler,
		WsHandler:   wsHandler,
		MsgHandler:  msgHansler,
	}
}
