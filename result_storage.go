package grep

import "sync"

type ResultStorage interface {
	Update(key, value interface{})
	Delete(key interface{})
	Get(key interface{}) (interface{}, bool)
}

type LightSyncMap struct {
	storage map[interface{}]interface{}
	mux     *sync.RWMutex
}

func MakeLightSyncMap() *LightSyncMap {
	return &LightSyncMap{
		storage: make(map[interface{}]interface{}),
		mux:     &sync.RWMutex{},
	}
}

func (s *LightSyncMap) Update(key, value interface{}) {
	s.mux.Lock()
	s.storage[key] = value
	s.mux.Unlock()
}

func (s *LightSyncMap) Delete(key interface{}) {
	s.mux.Lock()
	delete(s.storage, key)
	s.mux.Unlock()
}

func (s *LightSyncMap) Get(key interface{}) (interface{}, bool) {
	m, f := s.storage[key]
	return m, f
}

func (s *LightSyncMap) Len() int {
	return len(s.storage)
}
