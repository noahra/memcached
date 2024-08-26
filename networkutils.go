package main

import (
	"fmt"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn, memcache *Cache) error {
	buf := make([]byte, 1024)
	totalBytesRead := 0
	commands := map[string]struct{}{
		"set":     {},
		"add":     {},
		"append":  {},
		"prepend": {},
		"replace": {},
		"get":     {},
	}
	for {
		bytesRead, err := conn.Read(buf[totalBytesRead:])
		if err != nil {
			return err
		}
		totalBytesRead += bytesRead - 2

		words := strings.Fields(string(buf[:totalBytesRead]))
		if (len(words) > 1 && !containsCommand(words[0], commands)) || len(words) > 6 {
			resetBufferAndBytesRead(&buf, &totalBytesRead)
			continue
		}

		if len(words) == 2 && words[0] == "get" {
			_, ok := memcache.Get(words[1])
			if !ok {
				continue
			}
			err := createAndExecuteCommand(words, conn, memcache)
			if err != nil {
				return fmt.Errorf("err: %v", err)
			}
			resetBufferAndBytesRead(&buf, &totalBytesRead)
			continue
		}

		if len(words) >= 5 && containsCommand(words[0], commands) {
			if len(words) == 6 && words[len(words)-1] != "noreply" {
				continue
			}
			err := createAndExecuteCommand(words, conn, memcache)
			if err != nil {
				return fmt.Errorf("err: %v", err)
			}
			resetBufferAndBytesRead(&buf, &totalBytesRead)
			continue
		}
	}
}

func containsCommand(word string, commands map[string]struct{}) bool {
	_, exists := commands[word]
	return exists
}

func createAndExecuteCommand(words []string, conn net.Conn, memcache *Cache) error {
	command, err := CreateCommand(words, conn)
	if err != nil {
		return err
	}
	err = command.Execute(memcache)
	if err != nil {
		return err
	}
	return nil
}

func resetBufferAndBytesRead(buf *[]byte, totalBytesRead *int) {
	*buf = make([]byte, 1024)
	*totalBytesRead = 0
}
