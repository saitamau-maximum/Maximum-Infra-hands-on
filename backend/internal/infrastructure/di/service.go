package di

import (
	"fmt"

	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/infrastructure/gatewayImpl/cache"
	"example.com/infrahandson/internal/infrastructure/gatewayImpl/s3client"
	"example.com/infrahandson/internal/infrastructure/serviceImpl/iconStoreServiceImpl/localiconsvc"
	"example.com/infrahandson/internal/infrastructure/serviceImpl/iconStoreServiceImpl/s3iconsvc"
	"example.com/infrahandson/internal/infrastructure/serviceImpl/messageCacheImpl/memcachedmsg"
	"example.com/infrahandson/internal/infrastructure/serviceImpl/messageCacheImpl/memmsgcache"
	"example.com/infrahandson/internal/infrastructure/serviceImpl/websocketManagerImpl/memwsmanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bradfitz/gomemcache/memcache"
)

func ServiceInitialize(
	cfg *config.Config,
	repo repository.Repository,
) (service.Service, *memcache.Client) {
	// Serviceの初期化
	wsManager := memwsmanager.NewInMemoryWebSocketManager()
	var msgCache service.MessageCacheService
	var cacheClient *memcache.Client
	if cfg.MemcachedAddr != nil {
		// Memcachedの初期化
		fmt.Println("Memcached address:", *cfg.MemcachedAddr)
		cacheInitializer := cache.NewCacheInitializer(&cache.NewCacheInitializerParams{Cfg: cfg})

		cacheClient, err := cacheInitializer.Init()
		if err != nil {
			panic("failed to initialize memcached: " + err.Error())
		}
		fmt.Println("Memcached client initialized successfully")
		msgCache = memcachedmsg.NewMessageCacheService(&memcachedmsg.NewMessageCacheServiceParams{
			MsgRepo: repo.MessageRepository,
			Client:  cacheClient,
		})
	} else {
		msgCache = memmsgcache.NewMessageCacheService(&memmsgcache.NewMessageCacheServiceParams{MsgRepo: repo.MessageRepository})
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

		iconSvc = s3iconsvc.NewS3IconStoreImpl(s3iconsvc.NewS3IconStoreImplParams{
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
		iconSvc = localiconsvc.NewLocalIconStoreImpl(&localiconsvc.NewLocalIconStoreImplParams{
			DirPath: cfg.LocalIconDir,
		})
	}

	return service.Service{
		IconStoreService: iconSvc,
		MessageCacheService: msgCache,
		WebsocketManager: wsManager,
	}, cacheClient
}
