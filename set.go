package cache

import (
	"sync"
	"sync/atomic"
)

var set *Set

type Set struct {
	key  string
	kvs  *sync.Map
	size *atomic.Int64
}

func (s *Set) Get() (any, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Set) Update() error {
	//TODO implement me
	panic("implement me")
}

func (s *Set) Del() error {
	//TODO implement me
	panic("implement me")
}

func (s *Set) Key() string {
	return s.key
}

func (s *Set) Value() (any, error) {
	values := make([]any, s.size.Load())
	i := 0
	s.kvs.Range(func(k, v any) bool {
		values[i] = v
		i++
		return true
	})
	return values, nil
}

func (s *Set) Add(kv KV) {
	_, ok := s.kvs.Load(kv.Key())
	s.kvs.Store(kv.Key(), kv)
	if !ok {
		s.size.Add(1)
	}
}
func (s *Set) Remove(kv KV) {
	_, ok := s.kvs.Load(kv.Key())
	if ok {
		s.kvs.Delete(kv.Key())
		s.size.Add(-1)
	}
}

func NewSet(key string, ks ...KV) KV {
	kvs := new(sync.Map)
	for _, kv := range ks {
		kvs.Store(kv.Key(), kv)
	}
	return &Set{key: key, kvs: kvs, size: new(atomic.Int64)}
}
