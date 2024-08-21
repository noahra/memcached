package main

import (
	"time"
)

type cacheValue struct {
	key           string
	flags         string
	expiryTime    float64
	amountOfBytes int
	dataBlock     string
	createdAt     time.Time
}

func expiryCheck(cacheValue cacheValue) bool {
	if cacheValue.expiryTime != 0 {
		duration := time.Now().Sub(cacheValue.createdAt)
		if duration.Seconds() > cacheValue.expiryTime {
			return true
		}
	}
	return false
}
