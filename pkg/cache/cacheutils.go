package cache

import (
	"time"
)

type CacheValue struct {
	Key           string
	Flags         string
	ExpiryTime    float64
	AmountOfBytes int
	DataBlock     string
	CreatedAt     time.Time
}

func ExpiryCheck(cacheValue CacheValue) bool {
	if cacheValue.ExpiryTime != 0 {
		duration := time.Now().Sub(cacheValue.CreatedAt)
		if duration.Seconds() > cacheValue.ExpiryTime {
			return true
		}
	}
	return false
}
