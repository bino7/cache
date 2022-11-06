package cache

import "sync"

const (
	ADD_EVENT    = "add"
	UPDATE_EVENT = "update"
	DEL_EVENT    = "del"
)

type CacheResourceEvent struct {
	Event string
	Key   string
}

type NewFunc func(string) (Item, error)

type Store struct {
	new    NewFunc
	caches map[string]Item
	mu     sync.Mutex
}

func NewStore(new NewFunc) *Store {
	return &Store{new: new, caches: make(map[string]Item)}
}
func (s *Store) GetItem(key string) (Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if item, ok := s.caches[key]; ok {
		return item, nil
	}
	if s.new != nil {
		item, err := s.new(key)
		if err != nil {
			return nil, err
		}
		s.caches[key] = item
		return item, nil
	}
	return nil, nil
}
func (s *Store) AddItem(cache Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.caches[cache.Key()] = cache
}
func (s *Store) Get(key string) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cache, ok := s.caches[key]
	if ok {
		return cache.Get()
	}
	return nil, nil
}
func (s *Store) Notify(event *CacheResourceEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch event.Event {
	case ADD_EVENT:
		for _, cache := range s.caches {
			if !cache.Static() {
				cache.Invalid()
			}
		}
	case UPDATE_EVENT:
		for _, cache := range s.caches {
			if !cache.Static() {
				cache.Invalid()
			} else {
				/*if cache.Contains(event.Key) {

				  }
				  if cache.Len() == 0 {
				      cache.Invalid()
				      s.caches[cache.Key()] = nil
				  }*/
			}
		}
	case DEL_EVENT:
		for _, cache := range s.caches {
			if !cache.Static() {
				cache.Invalid()
			} else {
				/* cache.Remove(event.Key)
				   if cache.Len() == 0 {
				       cache.Invalid()
				       s.caches[cache.Key()] = nil
				   }*/
			}
		}
	}

}
