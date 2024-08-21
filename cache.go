package main

import "sync"

type CacheStore interface {
	Get(key string) (cacheValue, bool)
	Set(key string, value cacheValue)
	Delete(key string)
}

type Cache struct {
	store map[string]cacheValue
	mu    sync.RWMutex
}

func (cache Cache) Get(key string) (cacheValue, bool) {
	val, ok := cache.store[key]
	return val, ok
}

func (cache Cache) Set(key string, value cacheValue) {
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
		store: make(map[string]cacheValue),
	}
}
