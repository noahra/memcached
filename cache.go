package main

import "sync"

type CacheStore interface {
	Get(key string) (CacheValue, bool)
	Set(key string, value CacheValue)
	Delete(key string)
}

type Cache struct {
	store map[string]CacheValue
	mu    sync.RWMutex
}

func (cache Cache) Get(key string) (CacheValue, bool) {
	val, ok := cache.store[key]
	return val, ok
}

func (cache Cache) Set(key string, value CacheValue) {
	cache.mu.Lock()
	cache.store[key] = value
	cache.mu.Unlock()
}
func (cache Cache) Delete(key string) {
	cache.mu.Lock()
	delete(cache.store, key)
	cache.mu.Unlock()
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]CacheValue),
	}
}
