package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type CacheStore interface {
	Get(key string) (cacheValue, bool)
	Set(key string, value cacheValue)
	Delete(key string)
}

type Cache struct {
	store map[string]cacheValue
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]cacheValue),
	}
}

type Command interface {
	Execute(cache Cache) error
}

type Storable interface {
	Store(cache Cache) error
}

type Removable interface {
	Remove(cache Cache) error
}

type SetCommand struct {
	memCacheCommand
}

type AddCommand struct {
	memCacheCommand
}

type ReplaceCommand struct {
	memCacheCommand
}

func (cmd SetCommand) Execute(cache Cache) error {
	// Implement the logic for the "set" command
}

func (cmd AddCommand) Execute(cache Cache) error {
	// Implement the logic for the "add" command
}

func (cmd ReplaceCommand) Execute(cache Cache) error {
	// Implement the logic for the "replace" command
}

func (c memCacheCommand) CacheOperation(memcache map[string]cacheValue) {

	value := cacheValue{
		key:           c.key,
		flags:         c.flags,
		expiryTime:    c.expiryTime,
		amountOfBytes: c.byteCount,
		createdAt:     time.Now(),
	}
	byteSize := value.amountOfBytes
	buf := make([]byte, byteSize)
	_, err := c.connection.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	value.dataBlock = string(buf)

	_, ok := memcache[c.key]

	switch c.commandType {
	case "set":
		memcache[value.key] = value
	case "add":
		if ok {
			write_NOT_STORED(c.connection)
		} else {
			memcache[value.key] = value
		}
	case "replace":
		if !ok {
			write_NOT_STORED(c.connection)
		} else {
			memcache[value.key] = value
		}
	default:
	}

	if c.noReply != "noreply" {
		_, err := c.connection.Write([]byte("STORED\r\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

type memCacheCommand struct {
	connection  net.Conn
	commandType string
	key         string
	flags       string
	expiryTime  float64
	byteCount   int
	noReply     string
}
