package cache

import (
	"context"
	redis_cache "github.com/go-redis/cache/v8"
	"go.uber.org/fx"
)

var Module = fx.Module("cache",
	fx.Provide(newRedisClient, newRedisCache, newRedisPool),
	fx.Invoke(func(lc fx.Lifecycle, redisCache *redis_cache.Cache) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})
	}),
)
