package main

import (
	"net"
	"strings"
)

func HandleConnection(conn net.Conn, memcache *Cache) error {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			return err

		}

		words := strings.Fields(string(buf))
		if len(words) < 2 {
			continue
		}

		command, err := CreateCommand(words, conn)
		if err != nil {
			return err
		}
		err = command.Execute(memcache)
		if err != nil {
			return err
		}
	}
}
