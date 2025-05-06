package cache

import (
	"errors"

	"example.com/infrahandson/config"
	"example.com/infrahandson/internal/interface/gateway"
	"github.com/bradfitz/gomemcache/memcache"
)

type CacheInitializer struct {
	cfg    *config.Config
}

type NewCacheInitializerParams struct {
	Cfg    *config.Config
}

func (p *NewCacheInitializerParams) Validate() error {
	if p.Cfg == nil {
		return errors.New("config is required")
	}
	return nil
}

func NewCacheInitializer(p *NewCacheInitializerParams) gateway.CacheInitializer {
	if err := p.Validate(); err != nil {
		panic(err)
	}
	return &CacheInitializer{
		cfg:    p.Cfg,
	}
}

func (ci *CacheInitializer) Init() (*memcache.Client, error) {
	// これが使われる際には、すでにMemcachedがnil出ないことが保証されている。
	addr := *ci.cfg.MemcachedAddr
	client := memcache.New(addr)
	if client == nil {
		return nil, errors.New("failed to create memcache client")
	}
	return client, nil
}
