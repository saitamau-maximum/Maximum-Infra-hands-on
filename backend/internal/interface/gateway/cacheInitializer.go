package gateway

import "github.com/bradfitz/gomemcache/memcache"

type CacheInitializer interface {
	Init() (*memcache.Client, error)
}