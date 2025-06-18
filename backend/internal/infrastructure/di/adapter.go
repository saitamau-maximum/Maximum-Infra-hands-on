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

func InitializeAdapter(cfg *config.Config) adapter.Adapter {
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

	return adapter.Adapter{
		HasherAdapter:       hasher,
		TokenServiceAdapter: tokenService,
		LoggerAdapter:       logger,
		Upgrader:            upgrader,
	}
}
