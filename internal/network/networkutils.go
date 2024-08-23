package network

import (
	"ccmemcached/internal/cache"
	"ccmemcached/internal/commands"
	"fmt"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn, memcache *cache.Cache) error {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			return fmt.Errorf("err: %w", err)

		}
		words := strings.Fields(string(buf))
		if len(words) < 2 {
			continue // Skip if no command is received
		}

		command, err := commands.CreateCommand(words, conn)
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		err = command.Execute(memcache)
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
	}
}
