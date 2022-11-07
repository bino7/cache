package cache

import "sync"

type Key interface {
	Key() string
}

type KeySet interface {
	Key
	Len() int
	Contains(key string) bool
	Remove(key string)
}
type keySet struct {
	key  string
	keys []string
	mu   sync.Mutex
}

func NewKeySet(key string, keys []string) KeySet {
	return &keySet{key: key, keys: keys}
}

func (ks *keySet) Key() string {
	return ks.key
}
func (ks *keySet) Len() int {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	return len(ks.keys)
}
func (ks *keySet) Contains(key string) bool {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	for _, k := range ks.keys {
		if k == key {
			return true
		}
	}
	return false
}
func (ks *keySet) Add(key string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keys = append(ks.keys, key)
}
func (ks *keySet) Remove(key string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	j := -1
	for i, k := range ks.keys {
		if k == key {
			j = i
		}
	}
	if j >= 0 {
		ks.keys = append(ks.keys[:j], ks.keys[j+1:]...)
	}
}
