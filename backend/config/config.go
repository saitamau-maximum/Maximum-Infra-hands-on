package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// *stringは、nilのときに実装を切り替えるために使う
type Config struct {
	Port        string        // サーバーのポート番号
	CORSOrigin  string        // CORSのオリジン
	DBPath      string        // SQLite用データベースファイルの場所
	SecretKey   string        // JWTトークンの署名に使用する秘密鍵
	HashCost    int           // パスワードハッシュ化に使用するcost値
	TokenExpiry time.Duration // JWTトークンの有効期限
	// DB
	MySQLDSN *string // MySQL用データベースのDSN
	// Cache
	MemcachedAddr *string // Memcachedのアドレス
	// IconStore
	LocalIconDir       string  // ユーザーアイコンのローカル保存先
	IconStoreEndpoint  *string // ユーザーアイコンの保存先エンドポイント
	IconStoreBucket    *string // ユーザーアイコンの保存先バケット名
	IconStoreAccessKey *string // ユーザーアイコンの保存先アクセスキー
	IconStoreSecretKey *string // ユーザーアイコンの保存先シークレットキー
	IconStoreSecure    *string // ユーザーアイコンの保存先セキュアフラグ
	IconStoreBaseURL   *string // ユーザーアイコンの保存先URL
	IconStorePrefix    *string // ユーザーアイコンの保存先プレフィックス
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:5173"),
		DBPath:      getEnv("DB_PATH", "database.db"),
		SecretKey:   getEnv("SECRET_KEY", "secret"),
		HashCost:    parseInt(getEnv("HASH_COST", "10")),
		TokenExpiry: paraseDuration(getEnv("TOKEN_EXPIRY", "24h")),
		// DB
		MySQLDSN: parseStringPointer(getEnv("MYSQL_DSN", "")),
		// Cache
		MemcachedAddr: parseStringPointer(getEnv("MEMCACHED_ADDR", "")),
		//IconStore
		LocalIconDir:       getEnv("LOCAL_ICON_DIR", "./images/icons"),
		IconStoreEndpoint:  parseStringPointer(getEnv("ICON_STORE_ENDPOINT", "")),
		IconStoreBucket:    parseStringPointer(getEnv("ICON_STORE_BUCKET", "")),
		IconStoreAccessKey: parseStringPointer(getEnv("ICON_STORE_ACCESS_KEY", "")),
		IconStoreSecretKey: parseStringPointer(getEnv("ICON_STORE_SECRET_KEY", "")),
		IconStoreSecure:    parseStringPointer(getEnv("ICON_STORE_SECURE", "")),
		IconStoreBaseURL:   parseStringPointer(getEnv("ICON_STORE_BASE_URL", "")),
		IconStorePrefix:    parseStringPointer(getEnv("ICON_STORE_PREFIX", "")),
	}
}

func (c *Config) IsS3() (bool, []error) {
	var err []error
	if c.IconStoreEndpoint == nil {
		err = append(err, errors.New("IconStoreEndpoint is nil"))
	}
	if c.IconStoreBucket == nil {
		err = append(err, errors.New("IconStoreBucket is nil"))
	}
	if c.IconStoreAccessKey == nil {
		err = append(err, errors.New("IconStoreAccessKey is nil"))
	}
	if c.IconStoreSecretKey == nil {
		err = append(err, errors.New("IconStoreSecretKey is nil"))
	}
	if c.IconStoreBaseURL == nil {
		err = append(err, errors.New("IconStoreBaseURL is nil"))
	}
	if c.IconStorePrefix == nil {
		err = append(err, errors.New("IconStorePrefix is nil"))
	}
	
	if err != nil {
		return false, err
	}
	return true, nil
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
