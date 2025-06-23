package di

import (
	"example.com/infrahandson/config"
	mysqlgatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/mysql"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	"example.com/infrahandson/internal/interface/gateway"
	"example.com/infrahandson/internal/interface/handler/messagehandler"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"example.com/infrahandson/internal/interface/handler/websockethandler"
	"example.com/infrahandson/internal/usecase/messagecase"
	"example.com/infrahandson/internal/usecase/roomcase"
	"example.com/infrahandson/internal/usecase/usercase"
	"example.com/infrahandson/internal/usecase/websocketcase"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DB          *sqlx.DB
	Cache       *memcache.Client
	UserHandler userhandler.UserHandlerInterface
	RoomHandler roomhandler.RoomHandlerInterface
	WsHandler   websockethandler.WebSocketHandlerInterface
	MsgHandler  messagehandler.MessageHandlerInterface
}

func InitializeDependencies(cfg *config.Config) *Dependencies {

	// アダプターの初期化
	// 詳細は internal/infrastructure/di/adapter.go を参照
	adapters := InitializeAdapter(cfg)

	// Factoryの初期化
	// 詳細は internal/infrastructure/di/factory.go を参照
	factorys := InitializeFactory()

	// DBの初期化
	var initializer gateway.DBInitializer
	var dbType DBType
	// initializerの作成
	if cfg.MySQLDSN != nil {
		// MySQL用のDSNが設定されている場合、MySQL用のイニシャライザーを用意
		dbType = DBTypeMySQL
		initializer = mysqlgatewayimpl.NewMySQLInitializer(&mysqlgatewayimpl.NewMySQLInitializerParams{
			DSN:            cfg.MySQLDSN,
			MigrationsPath: "./internal/infrastructure/gatewayImpl/db/mysql/migrations",
		})
	} else {
		// SQLite用のイニシャライザーを用意
		dbType = DBTypeSQLite
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

	// Repositoryの初期化
	// 詳細は internal/infrastructure/di/repository.go を参照
	repositories := RepositoryInitialize(dbType, db)

	// Serviceの初期化
	// 詳細は internal/infrastructure/di/service.go を参照
	services, cacheClient := ServiceInitialize(cfg, repositories)

	// UseCaseの初期化
	userUseCase := usercase.NewUserUseCase(usercase.NewUserUseCaseParams{
		UserRepo:      repositories.UserRepository,
		Hasher:        adapters.HasherAdapter,
		TokenSvc:      adapters.TokenServiceAdapter,
		IconSvc:       services.IconStoreService,
		UserIDFactory: factorys.UserIDFactory,
	})
	roomUseCase := roomcase.NewRoomUseCase(roomcase.NewRoomUseCaseParams{
		RoomRepo:      repositories.RoomRepository,
		UserRepo:      repositories.UserRepository,
		RoomIDFactory: factorys.RoomIDFactory,
	})
	wsUseCase := websocketcase.NewWebsocketUseCase(websocketcase.NewWebsocketUseCaseParams{
		UserRepo:         repositories.UserRepository,
		RoomRepo:         repositories.RoomRepository,
		MsgRepo:          repositories.MessageRepository,
		MsgCache:         services.MessageCacheService,
		WsClientRepo:     repositories.WsClientRepository,
		WebsocketManager: services.WebsocketManager,
		MsgIDFactory:     factorys.MessageIDFactory,
		ClientIDFactory:  factorys.WsClientIDFactory,
	})
	msgUseCase := messagecase.NewMessageUseCase(messagecase.NewMessageUseCaseParams{
		MsgRepo:  repositories.MessageRepository,
		MsgCache: services.MessageCacheService,
		RoomRepo: repositories.RoomRepository,
		UserRepo: repositories.UserRepository,
	})

	// Handlerの初期化
	userHandler := userhandler.NewUserHandler(userhandler.NewUserHandlerParams{
		UserUseCase:   userUseCase,
		UserIDFactory: factorys.UserIDFactory,
		Logger:        adapters.LoggerAdapter,
	})
	roomHandler := roomhandler.NewRoomHandler(roomhandler.NewRoomHandlerParams{
		RoomUseCase:   roomUseCase,
		UserIDFactory: factorys.UserIDFactory,
		RoomIDFactory: factorys.RoomIDFactory,
		Logger:        adapters.LoggerAdapter,
	})
	wsHandler := websockethandler.NewWebSocketHandler(websockethandler.NewWebSocketHandlerParams{
		WsUseCase:     wsUseCase,
		WsUpgrader:    adapters.Upgrader,
		WsConnFactory: factorys.WsConnFactory,
		UserIDFactory: factorys.UserIDFactory,
		RoomIDFactory: factorys.RoomIDFactory,
		Logger:        adapters.LoggerAdapter,
	})
	msgHansler := messagehandler.NewMessageHandler(messagehandler.NewMessageHandlerParams{
		MsgUseCase: msgUseCase,
		Logger:     adapters.LoggerAdapter,
	})

	return &Dependencies{
		DB:          db,
		Cache:       cacheClient,
		UserHandler: userHandler,
		RoomHandler: roomHandler,
		WsHandler:   wsHandler,
		MsgHandler:  msgHansler,
	}
}
