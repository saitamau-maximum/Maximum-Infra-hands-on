package config

import (
	"os"
	"strconv"
	"time"
)

// *stringは、nilのときに実装を切り替えるために使う
type Config struct {
	Port             string        // サーバーのポート番号
	CORSOrigin       string        // CORSのオリジン
	DBPath           string        // SQLite用データベースファイルの場所
	SecretKey        string        // JWTトークンの署名に使用する秘密鍵
	HashCost         int           // パスワードハッシュ化に使用するcost値
	TokenExpiry      time.Duration // JWTトークンの有効期限
	MySQLDSN         *string       // MySQL用データベースのDSN
	MemcachedAddr    *string       // Memcachedのアドレス
	LocalIconDir     string       // ユーザーアイコンのローカル保存先
	IconStoreBaseURL *string       // ユーザーアイコンの保存先URL
	IconStoreBucket  *string       // ユーザーアイコンの保存先バケット名
	IconStorePrefix	 *string       // ユーザーアイコンの保存先プレフィックス
}

func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		CORSOrigin:       getEnv("CORS_ORIGIN", "http://localhost:5173"),
		DBPath:           getEnv("DB_PATH", "database.db"),
		SecretKey:        getEnv("SECRET_KEY", "secret"),
		HashCost:         parseInt(getEnv("HASH_COST", "10")),
		TokenExpiry:      paraseDuration(getEnv("TOKEN_EXPIRY", "24h")),
		MySQLDSN:         parseStringPointer(getEnv("MYSQL_DSN", "")),
		MemcachedAddr:    parseStringPointer(getEnv("MEMCACHED_ADDR", "")),
		LocalIconDir:     getEnv("LOCAL_ICON_DIR", "./images/icons"),
		IconStoreBaseURL: parseStringPointer(getEnv("ICON_STORE_BASE_URL", "")),
		IconStoreBucket: parseStringPointer(getEnv("ICON_STORE_BUCKET", "")),
		IconStorePrefix:  parseStringPointer(getEnv("ICON_STORE_PREFIX", "")),
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
