package di

import (
	"example.com/infrahandson/config"
	mysqlgatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/mysql"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	"example.com/infrahandson/internal/interface/gateway"
	"example.com/infrahandson/internal/interface/handler"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DB      *sqlx.DB
	Cache   *memcache.Client
	Handler *handler.Handler
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
// 詳細は internal/infrastructure/di/usecase.go を参照
	usecases := UseCaseInitialize(&UseCaseDependency{
		Adapter: adapters,
		Factory: factorys,
		Repo:    repositories,
		Svc:     services,
	})

	// Handlerの初期化
	// 詳細は internal/infrastructure/di/handler.go を参照
	handlers := HandlerInitialize(&HandlerInitializeParams{
		Adapter: adapters,
		Factory: factorys,
		UseCase: usecases,
	})

	return &Dependencies{
		DB:      db,
		Cache:   cacheClient,
		Handler: handlers,
	}
}
