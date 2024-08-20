package main

import (
	"net"
	"strconv"
	"time"
)

func expiryCheck(cacheValue cacheValue) bool {
	if cacheValue.expiryTime != 0 {
		duration := time.Now().Sub(cacheValue.createdAt)
		if duration.Seconds() > cacheValue.expiryTime {
			return true
		}

	}
	return false
}

func serializeWordsIntoCacheValue(words []string, conn net.Conn) memCacheCommand {
	expiryTime, _ := strconv.ParseFloat(words[3], 64)
	byteCount, _ := strconv.ParseInt(words[4], 10, 32)
	return memCacheCommand{
		connection:  conn,
		commandType: words[0],
		key:         words[1],
		flags:       words[2],
		expiryTime:  expiryTime,
		byteCount:   int(byteCount),
		noReply:     words[len(words)-2],
	}
}
