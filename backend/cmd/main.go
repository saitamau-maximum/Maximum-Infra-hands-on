package main

import (
	"example.com/infrahandson/config"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	"example.com/infrahandson/internal/infrastructure/server"
)

func main() {
	// 設定の読み込み
	cfg := config.LoadConfig()

	initializer := sqlitegatewayimpl.NewSQLiteInitializer(&sqlitegatewayimpl.NewSQLiteInitializerParams{
		Path:           cfg.DBPath,
		MigrationsPath: "./internal/infrastructure/gatewayImpl/db/sqlite/migrations",
	})
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
