package main

import "testing"

func TestCacheUtils(t *testing.T) {
	cacheValue := CacheValue{ExpiryTime: 1}
	hasExpired := ExpiryCheck(cacheValue)
	if !hasExpired {
		t.Error("HasExpired should return false")
	}
	cacheValue = CacheValue{ExpiryTime: 0}
	hasExpired = ExpiryCheck(cacheValue)
	if hasExpired {
		t.Error("HasExpired should return true")
	}
}
