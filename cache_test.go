package main

import (
	"testing"
	"time"
)

func TestCacheOperations(t *testing.T) {
	testCache := NewCache()
	testCache.Set("key1", CacheValue{
		Key:           "key1",
		Flags:         "0",
		ExpiryTime:    0,
		AmountOfBytes: 4,
		DataBlock:     "test",
		CreatedAt:     time.Now(),
	})
	_, doesExist := testCache.Get("key1")
	if !doesExist {
		t.Error("key does not exist")
	}
	testCache.Delete("key1")
	_, doesExist = testCache.Get("key1")
	if doesExist {
		t.Error("key does exist")
	}
}
