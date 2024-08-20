package main

import (
	"fmt"
	"net"
	"strings"
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

func handleConnection(conn net.Conn, memcache map[string]cacheValue) {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		words := strings.Fields(string(buf))
		if len(words) < 2 {
			continue // Skip if no command is received
		}
		if words[0] == "get" {
			handleGetCommand(conn, words[len(words)-2], memcache)
		} else {
			command := serializeWordsIntoCacheValue(words, conn)
			command.CacheOperation(memcache)
		}
	}
}

func write_NOT_STORED(conn net.Conn) {
	_, err := conn.Write([]byte("NOT_STORED\r\n"))
	if err != nil {
		fmt.Println(err)
	}
}
