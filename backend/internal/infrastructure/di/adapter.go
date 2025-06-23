// adapter インターフェースの実体を組み立てる
package di

import (
	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/infrastructure/adapterImpl/hasherAdapterImpl/bcrypt"
	"example.com/infrahandson/internal/infrastructure/adapterImpl/loggerAdapterImpl/fmtLogger"
	"example.com/infrahandson/internal/infrastructure/adapterImpl/tokenServiceAdapterImpl/JWT"
	"example.com/infrahandson/internal/infrastructure/adapterImpl/upgraderAdapterImpl/gorillaupgrader"
	"example.com/infrahandson/internal/interface/adapter"
)

// InitializeAdapter はアダプターの初期化を行います。
// cfg でアプリケーション設定を受け取ります（詳細：config/config.go）
// 返り値 adapter.Adapter は adapter 層をまとめた構造体です。（詳細：internal/interface/adapter/adapter.go）
func InitializeAdapter(cfg *config.Config) *adapter.Adapter {
	// ハッシュアダプターの初期化
	hasher := bcrypt.NewHasherAdapter(bcrypt.NewHasherAddapterParams{
		Cost: cfg.HashCost,
	})

	// トークンサービスアダプターの初期化
	tokenService := jwt.NewTokenServiceAdapter(jwt.NewTokenServiceAdapterParams{
		SecretKey:     cfg.SecretKey,
		ExpireMinutes: int(cfg.TokenExpiry),
	})

	// Loggerの設定
	logger := fmtlogger.NewFmtLogger()

	// WebSocketアップグレーダーの初期化
	upgrader := gorillaupgrader.NewGorillaWebSocketUpgrader()

	return &adapter.Adapter{
		HasherAdapter:       hasher,
		TokenServiceAdapter: tokenService,
		LoggerAdapter:       logger,
		Upgrader:            upgrader,
	}
}
