package main

import (
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

		if len(words) > 1 && !containsCommand(words[0], commands) {
			buf = make([]byte, 1024)
			totalBytesRead = 0
			continue
		}

		if len(words) > 6 {
			buf = make([]byte, 1024)
			totalBytesRead = 0
			continue
		}

		if len(words) == 2 && words[0] == "get" {
			_, ok := memcache.Get(words[1])
			if !ok {
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
			buf = make([]byte, 1024)
			totalBytesRead = 0
			continue
		}

		if len(words) >= 5 && containsCommand(words[0], commands) {
			if len(words) == 6 && words[len(words)-1] != "noreply" {
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
			buf = make([]byte, 1024)
			totalBytesRead = 0
			continue
		}
	}
}

func containsCommand(word string, commands map[string]struct{}) bool {
	_, exists := commands[word]
	return exists
}
