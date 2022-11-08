package cache

import rediscache "github.com/go-redis/cache/v8"

type KV interface {
	Key() string
	Value() (any, error)
	Update() error
	Del() error
}

type kv struct {
	key string
}

func (kv *kv) Update(value any) error {
	//TODO implement me
	panic("implement me")
}

func (kv *kv) Del() error {
	//TODO implement me
	panic("implement me")
}

func (kv *kv) Key() string {
	return kv.key
}

func (kv *kv) Value() (any, error) {
	var dest interface{}
	err := kv.get(dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (kv *kv) get(dest any) error {
	return redis.Get(defaultContext(), kv.key, dest)
}
func (kv *kv) set(value any) error {
	c.err = redis.Set(&rediscache.Item{
		Ctx:   defaultContext(),
		Key:   kv.key,
		Value: value,
		TTL:   ttl,
	})
}

func Find(pattern string) (KV, error) {
	keys, err := client.Keys(defaultContext(), pattern).Result()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}
	if len(keys) == 1 {
		return &kv{key: keys[0]}, nil
	}
	return NewSet(pattern, keys...), nil
}
