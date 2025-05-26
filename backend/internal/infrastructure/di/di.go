package di

import (
	"fmt"

	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	bcryptadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/hasherAdapterImpl/bcrypt"
	fmtloggerimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/loggerAdapterImpl/fmtLogger"
	tokenadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/tokenServiceAdapterImpl/JWT"
	gorillawebsocketupgraderImpl "example.com/infrahandson/internal/infrastructure/adapterImpl/upgraderAdapterImpl/gorillawebsocket"
	factoryimpl "example.com/infrahandson/internal/infrastructure/factoryImpl"
	"example.com/infrahandson/internal/infrastructure/gatewayImpl/cache"
	mysqlgatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/mysql"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	"example.com/infrahandson/internal/infrastructure/gatewayImpl/s3client"
	mysqlmsgrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/mysql"
	sqlitemsgrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/sqlite"
	mysqlroomrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/mysql"
	sqliteroomrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/sqlite"
	mysqluserrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/mysql"
	sqliteuserrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/sqlite"
	inmemorywsclientrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/websocketClientRepositoryImpl/InMemory"
	s3iconstoreimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/iconStoreServiceImpl/S3"
	localiconstoreimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/iconStoreServiceImpl/local"
	inmemorymsgcacheimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/messageCacheImpl/Inmemory"
	memcachedmsgcacheimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/messageCacheImpl/memcached"
	inmemorywsmanagerimpl "example.com/infrahandson/internal/infrastructure/serviceImpl/websocketManagerImpl/InMemory"
	"example.com/infrahandson/internal/interface/gateway"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/usecase"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DB          *sqlx.DB
	Cache       *memcache.Client
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
			MsgRepo: msgRepository,
			Client:  cacheClient,
		})
	} else {
		msgCache = inmemorymsgcacheimpl.NewMessageCacheService(&inmemorymsgcacheimpl.NewMessageCacheServiceParams{MsgRepo: msgRepository})
	}
	
	var iconSvc service.IconStoreService
	isS3, Type, errs := cfg.IsS3()
	if isS3 {
		var StorageClient *s3.Client
		if Type == "aws_s3" {

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
		if len(errs) != 6 {
			fmt.Println("S3の初期化に必要な情報が不足しています。")
			for _, err := range errs {
				fmt.Println(err)
			}
		}
		iconSvc = localiconstoreimpl.NewLocalIconStoreImpl(&localiconstoreimpl.NewLocalIconStoreImplParams{
			DirPath: cfg.LocalIconDir,
		})
	}

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
		IconSvc:       iconSvc,
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
		MsgCache:         msgCache,
		WsClientRepo:     wsClientRepository,
		WebsocketManager: wsManager,
		MsgIDFactory:     MsgIDFactory,
		ClientIDFactory:  clientDFactory,
	})
	msgUseCase := usecase.NewMessageUseCase(usecase.NewMessageUseCaseParams{
		MsgRepo:  msgRepository,
		MsgCache: msgCache,
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
		Cache:       cacheClient,
		UserHandler: userHandler,
		RoomHandler: roomHandler,
		WsHandler:   wsHandler,
		MsgHandler:  msgHansler,
	}
}
