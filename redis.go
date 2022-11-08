package cache

import (
	"context"
	"fmt"
	goredis "github.com/go-redis/redis/v8"
	redigo "github.com/gomodule/redigo/redis"
)

func newRedisPool(conf *Config) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: conf.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			return dialWithDB("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), conf.Password, conf.DB)
		},
	}
}

func dialWithDB(network, address, password string, DB int) (redigo.Conn, error) {
	c, err := dial(network, address, password)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("SELECT", DB); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

func dial(network, address, password string) (redigo.Conn, error) {
	c, err := redigo.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

var client *goredis.Client

func newRedisClient(conf *Config) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:        fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		DB:          conf.DB,
		Password:    conf.Password,
		PoolSize:    conf.MaxActive,
		IdleTimeout: conf.IdleTimeout,
	})
}

var ctx context.Context

func defaultContext() context.Context {
	return ctx
}
