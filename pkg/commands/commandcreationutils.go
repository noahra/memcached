package commands

import (
	"fmt"
	"memcached/pkg/cache"
	"strconv"
	"time"
)

func parseExpiry(expiryTimeString string) float64 {
	expiryTime, _ := strconv.ParseFloat(expiryTimeString, 64)
	return expiryTime
}
func parseByteCount(byteCountString string) int {
	byteCount, _ := strconv.ParseInt(byteCountString, 10, 32)
	return int(byteCount)
}

func cacheValueParser(command Command) cache.CacheValue {
	baseCmd := command.GetBaseCommand()
	value := cache.CacheValue{
		Key:           baseCmd.key,
		Flags:         baseCmd.flags,
		ExpiryTime:    baseCmd.expiryTime,
		AmountOfBytes: baseCmd.byteCount,
		CreatedAt:     time.Now(),
	}
	return value
}

func parseDataBlock(cmd Command) (string, error) {
	buf := make([]byte, cmd.GetBaseCommand().byteCount)
	_, err := cmd.GetBaseCommand().connection.Read(buf)
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}

	return string(buf), nil
}
