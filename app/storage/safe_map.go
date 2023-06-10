package storage

import (
	"sync"
)

type SafeMap struct {
	mu sync.RWMutex
	m  map[string][]byte
}

func (s *SafeMap) GetLen() int {
	return len(s.m)
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string][]byte),
	}
}

func (s *SafeMap) Write(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *SafeMap) Read(key string) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[key]
}
