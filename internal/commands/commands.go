package commands

import (
	"ccmemcached/internal/cache"
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
	Execute(cache *cache.Cache) error
	GetBaseCommand() BaseCommand
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

func (cmd SetCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd GetCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}

func (cmd ReplaceCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}
func (cmd AddCommand) GetBaseCommand() BaseCommand {
	return cmd.BaseCommand
}
func (cmd SetCommand) Execute(memcache *cache.Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := createBuffer(value.AmountOfBytes, cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.DataBlock = dataBlock

	memcache.Set(value.Key, value)

	if cmd.noReply != "noreply" {
		_, err := cmd.connection.Write([]byte("STORED\r\n"))
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
	}
	return nil
}

func (cmd AddCommand) Execute(memcache *cache.Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := createBuffer(value.AmountOfBytes, cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.DataBlock = dataBlock

	_, ok := memcache.Get(cmd.key)
	if ok {
		writeNotStored(cmd.connection)
	} else {
		memcache.Set(value.Key, value)
		if cmd.noReply != "noreply" {
			_, err := cmd.connection.Write([]byte("STORED\r\n"))
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
		}
	}

	return nil
}

func (cmd ReplaceCommand) Execute(memcache *cache.Cache) error {
	value := cacheValueParser(cmd)
	dataBlock, err := createBuffer(value.AmountOfBytes, cmd)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	value.DataBlock = dataBlock

	_, ok := memcache.Get(cmd.key)

	if !ok {
		writeNotStored(cmd.connection)
	} else {
		memcache.Set(value.Key, value)
		if cmd.noReply != "noreply" {
			_, err := cmd.connection.Write([]byte("STORED\r\n"))
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
		}
	}
	return nil
}

func (cmd GetCommand) Execute(memcache *cache.Cache) error {
	val, ok := memcache.Get(cmd.key)
	if ok {
		if !cache.ExpiryCheck(val) {
			cmd.connection.Write([]byte("VALUE " + val.Key + " " + val.Flags + " " + strconv.Itoa(val.AmountOfBytes) + "\n"))
			cmd.connection.Write([]byte(val.DataBlock + "\n"))
		} else {
			memcache.Delete(cmd.key)
		}
		cmd.connection.Write([]byte("END\r\n"))
	}
	return nil
}

func (cmd *BaseCommand) readDataBlock() (string, error) {
	buf := make([]byte, cmd.byteCount)
	_, err := cmd.connection.Read(buf)
	if err != nil {
		return "", fmt.Errorf("error reading data block: %w", err)
	}
	return string(buf), nil
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

func createBuffer(amountOfBytes int, cmd Command) (string, error) {
	byteSize := amountOfBytes
	buf := make([]byte, byteSize)
	_, err := cmd.GetBaseCommand().connection.Read(buf)
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}

	return string(buf), nil
}
