package main

import (
	"fmt"
	"strconv"
	"time"
)

func parseExpiry(seconds int) time.Time {
	currentTime := time.Now()
	duration := time.Duration(seconds) * time.Second
	return currentTime.Add(duration)
}
func parseByteCount(byteCountString string) int {
	byteCount, _ := strconv.ParseInt(byteCountString, 10, 32)
	return int(byteCount)
}

func cacheValueParser(command Command) CacheValue {
	baseCmd := command.GetBaseCommand()
	value := CacheValue{
		Key:           baseCmd.key,
		Flags:         baseCmd.flags,
		ExpiryTime:    parseExpiry(baseCmd.expiryTime),
		AmountOfBytes: baseCmd.byteCount,
		CreatedAt:     time.Now(),
	}
	return value
}

func parseDataBlock(cmd Command) ([]byte, error) {
	buf := make([]byte, cmd.GetBaseCommand().byteCount)
	_, err := cmd.GetBaseCommand().connection.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return buf, nil
}
