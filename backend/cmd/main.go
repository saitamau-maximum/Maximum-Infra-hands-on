package main

import (
	"log"

	"example.com/infrahandson/config"
	mysqlgatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/mysql"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	"example.com/infrahandson/internal/infrastructure/server"
	"example.com/infrahandson/internal/interface/gateway"
)

func main() {
	// 設定の読み込み
	cfg := config.LoadConfig()
	log.Println("DSN: ", *cfg.MySQLDSN)
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
	// データベースの初期化
	db, err := initializer.Init()
	if err != nil {
		panic("failed to initialize database: " + err.Error())
	}
	defer db.Close()
	// スキーマの初期化
	if err := initializer.InitSchema(db); err != nil {
		panic("failed to initialize schema: " + err.Error())
	}
	// サーバーの起動
	e := server.ServerStart(cfg, db)        // Echoインスタンスを取得
	e.Logger.Fatal(e.Start(":" + cfg.Port)) // サーバーを起動
}
