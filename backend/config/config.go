package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string // サーバーのポート番号
	DBPath      string // SQLite用データベースファイルの場所
	SecretKey   string // JWTトークンの署名に使用する秘密鍵
	HashCost    int // パスワードハッシュ化に使用するcost値
	TokenExpiry time.Duration // JWTトークンの有効期限
	MySQLDSN    *string // MySQL用データベースのDSN
	// MySQL用データベースのDSNをポインタにしているのは，環境変数が設定されていない場合にnilにするため
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "database.db"),
		SecretKey:   getEnv("SECRET_KEY", "secret"),
		HashCost:    parseInt(getEnv("HASH_COST", "10")),
		TokenExpiry: paraseDuration(getEnv("TOKEN_EXPIRY", "24h")),
		MySQLDSN:    parseStringPointer(getEnv("MYSQL_DSN", "")),
	}
}

func getEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return i
}

func paraseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return d
}

func parseStringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
