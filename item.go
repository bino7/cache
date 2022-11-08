package cache

import (
	"context"
	rediscache "github.com/go-redis/cache/v8"
	"time"
)

type Item interface {
	KV
	Get() (any, error)
	Update() error
	Del() error
	CacheErr() error
	Static() bool
	Valid() bool
	Invalid()
	TTL() time.Duration
	PrimaryKey() KV
}

type SaveFunc func(value interface{}) error
type GetFunc func(dest interface{}) error
type DelFunc func(interface{}) error
type KeyFunc func(args ...interface{}) (string, error)

func (getFunc GetFunc) Then(fn ...func()) GetFunc {
	return func(dest interface{}) error {
		err := getFunc(dest)
		if err != nil {
			return err
		}
		for _, f := range fn {
			f()
		}
		return nil
	}

}

type Option struct {
	GetFunc  GetFunc
	SaveFunc SaveFunc
	DelFunc  DelFunc
	TTL      time.Duration
}

func New(key string, value any, o *Option) Item {
	return &CacheItem{
		cache:    redis,
		getFunc:  o.GetFunc,
		saveFunc: o.SaveFunc,
		delFunc:  o.DelFunc,
		ttl:      o.TTL,
		key:      key,
		value:    value,
	}
}

type CacheItem struct {
	cache      *rediscache.Cache
	getFunc    GetFunc
	keyFunc    KeyFunc
	saveFunc   SaveFunc
	delFunc    DelFunc
	key        string
	value      any
	ttl        time.Duration
	err        error
	static     bool
	valid      bool
	len        int
	primaryKey KV
	associates []Item
}

func (c *CacheItem) TTL() time.Duration {
	return c.ttl
}
func (c *CacheItem) PrimaryKey() KV {
	return c.primaryKey
}
func (c *CacheItem) Static() bool {
	return c.static
}
func (c *CacheItem) Valid() bool {
	return c.valid
}
func (c *CacheItem) Invalid() {
	c.valid = false
}
func (c *CacheItem) defaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	return ctx
}
func (c *CacheItem) Get() (dest any, err error) {
	return c.get()
}
func (c *CacheItem) get() (any, error) {
	if c.primaryKey != nil {
		return c.primaryKey.Value(), nil
	}
	key := c.key
	c.err = c.cache.Get(c.defaultContext(), key, c.value)
	if c.err == nil {
		return c.value, nil
	}
	err := c.getFunc(c.value)
	if err != nil {
		return nil, err
	}
	c.err = c.set(c.key, c.value)
	if c.err == nil {
		c.valid = true
	}
	return c.value, nil
}
func (c *CacheItem) set(key string, dest any) error {
	c.err = c.cache.Set(&rediscache.Item{
		Ctx:   defaultContext(),
		Key:   key,
		Value: dest,
		TTL:   c.ttl,
	})
	return c.err
}
func (c *CacheItem) Update() error {
	key := c.key
	if c.saveFunc != nil {
		err := c.saveFunc(c.value)
		if err != nil {
			return err
		}
	}
	c.err = c.cache.Delete(c.defaultContext(), key)
	_ = c.set(key, c.value)
	return nil
}
func (c *CacheItem) Del() error {
	key := c.key
	c.err = c.cache.Delete(c.defaultContext(), key)
	if c.delFunc != nil {
		err := c.delFunc(c.value)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *CacheItem) CacheErr() error {
	return c.err
}
func (c *CacheItem) Key() string {
	return c.key
}
func (c *CacheItem) Value() (any, error) {
	return c.get()
}
