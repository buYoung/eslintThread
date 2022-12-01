package module

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"log"
	"reflect"
	"time"
)

type cacheInstance struct {
	Cache *bigcache.BigCache
}

func (c *cacheInstance) Init() error {
	config := bigcache.DefaultConfig(10 * time.Second)
	config.Shards = 64
	config.OnRemove = c.onRemove
	config.StatsEnabled = false
	config.HardMaxCacheSize = 10
	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return err
	}
	c.Cache = cache
	return nil
}

func (c *cacheInstance) onRemove(key string, entry []byte) {
	log.Println("removeKey", key, entry)
}

func (c *cacheInstance) Set(key string, value interface{}) error {
	log.Println(reflect.TypeOf(value).Name(), c.isFind(key))
	return nil
}

func (c *cacheInstance) isFind(key string) bool {
	_, err := c.Cache.Get(key)
	if err != nil {
		return false
	}
	return true
}
