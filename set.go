package cache

import (
	"sync"
)

type Set struct {
	items *sync.Map
}

func NewSet(items ...Item) *Set {
	m := new(sync.Map)
	for _, item := range items {
		m.Store(item.Key(), item)
	}
	return &Set{m}
}
func (s *Set) Put(item Item) {
	s.items.Store(item.Key(), item)
}

func (s *Set) Get(key string) Item {
	v, ok := s.items.Load(key)
	if ok {
		return v.(Item)
	}
	return nil
}
func (s *Set) Del() error {
	var err error
	s.items.Range(func(key, value any) bool {
		err = value.(Item).Delete()
		return true
	})
	return err
}
func (s *Set) GetValue(key string) (any, error) {
	item := s.Get(key)
	if item == nil {
		return nil, nil
	}
	v, err := item.Get()
	if err != nil {
		return nil, err
	}
	return v, nil
}
func (s *Set) Values() ([]any, error) {
	var err error
	values := make([]any, 0)
	s.items.Range(func(key any, v any) bool {
		item := v.(*CacheItem)
		value, getErr := item.Get()
		err = getErr
		values = append(values, value)
		return true
	})
	return values, err
}
func (s *Set) Contains(key string) bool {
	_, ok := s.items.Load(key)
	return ok
}
func (s *Set) FirstValue() (any, error) {
	return nil, nil
}
func (s *Set) Len() int {
	c := 0
	s.items.Range(func(key, v any) bool {
		c++
		return true
	})
	return c
}
func (s *Set) Remove(key string) {
	s.items.Delete(key)
}