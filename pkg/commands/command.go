package commands

import (
	"memcached/pkg/cache"
	"net"
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

type AppendCommand struct {
	BaseCommand
}

type PrependCommand struct {
	BaseCommand
}
