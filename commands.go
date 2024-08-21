package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type BaseCommand struct {
	key        string
	flags      string
	expiryTime float64
	byteCount  int
	noReply    string
	connection net.Conn
}

type Command interface {
	Execute(cache *Cache) error
}

type SetCommand struct {
	BaseCommand
}

type AddCommand struct {
	BaseCommand
}

type ReplaceCommand struct {
	BaseCommand
}

type GetCommand struct {
	BaseCommand
}

func (cmd SetCommand) Execute(cache *Cache) error {
	value := cacheValue{
		key:           cmd.key,
		flags:         cmd.flags,
		expiryTime:    cmd.expiryTime,
		amountOfBytes: cmd.byteCount,
		createdAt:     time.Now(),
	}
	byteSize := value.amountOfBytes
	buf := make([]byte, byteSize)
	_, err := cmd.connection.Read(buf)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.dataBlock = string(buf)

	cache.Set(value.key, value)

	if cmd.noReply != "noreply" {
		_, err := cmd.connection.Write([]byte("STORED\r\n"))
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}
	return nil
}

func (cmd AddCommand) Execute(cache *Cache) error {
	value := cacheValue{
		key:           cmd.key,
		flags:         cmd.flags,
		expiryTime:    cmd.expiryTime,
		amountOfBytes: cmd.byteCount,
		createdAt:     time.Now(),
	}
	byteSize := value.amountOfBytes
	buf := make([]byte, byteSize)
	_, err := cmd.connection.Read(buf)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.dataBlock = string(buf)

	_, ok := cache.Get(cmd.key)
	if ok {
		writeNotStored(cmd.connection)
	} else {
		cache.Set(value.key, value)
		if cmd.noReply != "noreply" {
			_, err := cmd.connection.Write([]byte("STORED\r\n"))
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
		}
	}

	return nil
}

func (cmd ReplaceCommand) Execute(cache *Cache) error {
	value := cacheValue{
		key:           cmd.key,
		flags:         cmd.flags,
		expiryTime:    cmd.expiryTime,
		amountOfBytes: cmd.byteCount,
		createdAt:     time.Now(),
	}
	byteSize := value.amountOfBytes
	buf := make([]byte, byteSize)
	_, err := cmd.connection.Read(buf)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.dataBlock = string(buf)

	_, ok := cache.Get(cmd.key)

	if !ok {
		writeNotStored(cmd.connection)
	} else {
		cache.Set(value.key, value)
		if cmd.noReply != "noreply" {
			_, err := cmd.connection.Write([]byte("STORED\r\n"))
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
		}
	}
	return nil
}

func (cmd GetCommand) Execute(cache *Cache) error {
	val, ok := cache.Get(cmd.key)
	if ok {
		if !expiryCheck(val) {
			cmd.connection.Write([]byte("VALUE " + val.key + " " + val.flags + " " + strconv.Itoa(val.amountOfBytes) + "\n"))
			cmd.connection.Write([]byte(val.dataBlock + "\n"))
		} else {
			cache.Delete(cmd.key)
		}
		cmd.connection.Write([]byte("END\r\n"))
	}
	return nil
}
