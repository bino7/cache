package cache

import (
	rediscache "github.com/go-redis/cache/v8"
	goredis "github.com/go-redis/redis/v8"
)

func Redis() *rediscache.Cache {
	return redis
}

var redis *rediscache.Cache

func newRedisCache(client *goredis.Client) *rediscache.Cache {
	redis = rediscache.New(&rediscache.Options{
		Redis: client,
	})
	return redis
}
