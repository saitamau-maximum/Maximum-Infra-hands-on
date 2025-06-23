package di

import (
	"fmt"

	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/infrastructure/gatewayImpl/cache"
	mysqlgatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/mysql"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	"example.com/infrahandson/internal/infrastructure/gatewayImpl/s3client"
	s3iconstoreimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/iconStoreServiceImpl/S3"
	localiconstoreimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/iconStoreServiceImpl/local"
	inmemorymsgcacheimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/messageCacheImpl/Inmemory"
	memcachedmsgcacheimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/messageCacheImpl/memcached"
	inmemorywsmanagerimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/websocketManagerImpl/InMemory"
	"example.com/infrahandson/internal/interface/gateway"
	"example.com/infrahandson/internal/interface/handler/messagehandler"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"example.com/infrahandson/internal/interface/handler/websockethandler"
	"example.com/infrahandson/internal/usecase/messagecase"
	"example.com/infrahandson/internal/usecase/roomcase"
	"example.com/infrahandson/internal/usecase/usercase"
	"example.com/infrahandson/internal/usecase/websocketcase"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	adapters := InitializeAdapter(cfg)

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
	repositories := RepoistoryInitialize(dbType, db)

	// Serviceの初期化
	wsManager := inmemorywsmanagerimpl.NewInMemoryWebSocketManager()
	var msgCache service.MessageCacheService
	var cacheClient *memcache.Client
	if cfg.MemcachedAddr != nil {
		// Memcachedの初期化
		fmt.Println("Memcached address:", *cfg.MemcachedAddr)
		cacheInitializer := cache.NewCacheInitializer(&cache.NewCacheInitializerParams{Cfg: cfg})
		cacheClient, err = cacheInitializer.Init()
		if err != nil {
			panic("failed to initialize memcached: " + err.Error())
		}
		fmt.Println("Memcached client initialized successfully")
		msgCache = memcachedmsgcacheimpl.NewMessageCacheService(&memcachedmsgcacheimpl.NewMessageCacheServiceParams{
			MsgRepo: repositories.MessageRepository,
			Client:  cacheClient,
		})
	} else {
		msgCache = inmemorymsgcacheimpl.NewMessageCacheService(&inmemorymsgcacheimpl.NewMessageCacheServiceParams{MsgRepo: repositories.MessageRepository})
	}

	var iconSvc service.IconStoreService
	isS3, Type, errs := cfg.IsS3()
	if isS3 {
		var StorageClient *s3.Client
		if Type == "aws_s3" {
			// TODO: AWS S3のクライアントを初期化
		} else if Type == "minio" {
			StorageClient = s3client.NewMinIOClient(s3client.NewMinIOClientParams{
				Endpoint:  *cfg.IconStoreEndpoint,
				AccessKey: *cfg.IconStoreAccessKey,
				SecretKey: *cfg.IconStoreSecretKey,
			})
		} else {
			panic("invalid storage type")
		}

		iconSvc = s3iconstoreimpl.NewS3IconStoreImpl(s3iconstoreimpl.NewS3IconStoreImplParams{
			BaseURL: *cfg.IconStoreBaseURL,
			Client:  StorageClient,
			Bucket:  *cfg.IconStoreBucket,
			Prefix:  *cfg.IconStorePrefix,
		})
	} else {
		// S3関連のすべてがnilではない場合は、不足している旨を表示
		if len(errs) != 5 {
			fmt.Println("S3の初期化に必要な情報が不足しています。")
			for _, err := range errs {
				fmt.Println(err)
			}
		}
		iconSvc = localiconstoreimpl.NewLocalIconStoreImpl(&localiconstoreimpl.NewLocalIconStoreImplParams{
			DirPath: cfg.LocalIconDir,
		})
	}

	// UseCaseの初期化
	userUseCase := usercase.NewUserUseCase(usercase.NewUserUseCaseParams{
		UserRepo:      repositories.UserRepository,
		Hasher:        adapters.HasherAdapter,
		TokenSvc:      adapters.TokenServiceAdapter,
		IconSvc:       iconSvc,
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
		MsgCache:         msgCache,
		WsClientRepo:     repositories.WsClientRepository,
		WebsocketManager: wsManager,
		MsgIDFactory:     factorys.MessageIDFactory,
		ClientIDFactory:  factorys.WsClientIDFactory,
	})
	msgUseCase := messagecase.NewMessageUseCase(messagecase.NewMessageUseCaseParams{
		MsgRepo:  repositories.MessageRepository,
		MsgCache: msgCache,
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
