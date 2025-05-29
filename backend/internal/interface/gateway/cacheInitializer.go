// キャッシュの初期化インターフェース
// 具体実装はbackend/internal/infrastructure/gatewayImpl/cache
package gateway

import "github.com/bradfitz/gomemcache/memcache"

type CacheInitializer interface {
	// Init はキャッシュクライアントを初期化します。
	Init() (*memcache.Client, error)
}