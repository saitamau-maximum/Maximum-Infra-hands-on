package main

import (
	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/infrastructure/server"
)

func main() {
	// 設定の読み込み
	cfg := config.LoadConfig()

	// サーバーの起動
	e, db, client := server.ServerStart(cfg) // Echoインスタンスを取得
	defer db.Close()
	if client != nil {
		// No Close method for Memcached client; cleanup logic not required
	}
	e.Logger.Fatal(e.Start(":" + cfg.Port)) // サーバーを起動
}
