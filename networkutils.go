package main

import (
	"bytes"
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
		totalBytesRead += bytesRead
		if idx := bytes.Index(buf[:totalBytesRead], []byte("\r\n")); idx != -1 {
			// Process the first command
			err := processCommand(buf[:idx], memcache, conn, commands)
			if err != nil {
				return fmt.Errorf("Error processing command: %v", err)
			}

			// Calculate the number of remaining bytes
			remainingBytes := totalBytesRead - (idx + 2)

			// Move the remaining data to the beginning of the buffer
			copy(buf, buf[idx+2:totalBytesRead])

			// Reset totalBytesRead
			totalBytesRead = remainingBytes
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

func processCommand(buf []byte, memcache *Cache, conn net.Conn, commands map[string]struct{}) error {
	individualFields := bytes.Fields(buf)

	if len(individualFields) == 0 {
		return nil
	}
	words := strings.Fields(string(buf))

	if len(individualFields) == 2 && bytes.Equal(individualFields[0], []byte("get")) {
		_, ok := memcache.Get(string(individualFields[1]))
		if !ok {
			return nil
		}

		err := createAndExecuteCommand(words, conn, memcache)
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}
		return nil
	}

	if len(words) >= 5 && containsCommand(words[0], commands) {
		err := createAndExecuteCommand(words, conn, memcache)
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}
	}
	return nil
}
