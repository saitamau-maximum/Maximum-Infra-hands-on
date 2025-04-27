package main

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/infrastructure/gateway_impl"
	"example.com/webrtc-practice/internal/server"
)

func main() {
	// 設定の読み込み
	cfg := config.LoadConfig()
	
	initializer := gateway_impl.NewSQLiteInitializer(cfg.DBPath)
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
	server.ServerStart(cfg, db)
}
