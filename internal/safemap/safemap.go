package safemap

import "sync"

type SafeMap[K string, V any] struct {
	sync.RWMutex
	m map[K]V
}

func New[K string, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: make(map[K]V),
	}
}

func NewFromMap[K string, V any](m map[K]V) *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: m,
	}
}

func (s *SafeMap[K, V]) Put(k K, v V) {
	s.Lock()
	s.m[k] = v
	s.Unlock()
}

func (s *SafeMap[K, V]) Get(k K) (v V, ok bool) {
	s.RLock()
	v, ok = s.m[k]
	s.RUnlock()
	return
}

func (s *SafeMap[K, V]) MustGet(k K) (v V) {
	s.RLock()
	v = s.m[k]
	s.RUnlock()
	return
}

func (s *SafeMap[K, V]) Delete(k K) {
	s.Lock()
	delete(s.m, k)
	s.Unlock()
}

func (s *SafeMap[K, V]) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

func (s *SafeMap[K, V]) Range(f func(K, V) bool) {
	s.RLock()
	for k, v := range s.m {
		if !f(k, v) {
			break
		}
	}
	s.RUnlock()
}
