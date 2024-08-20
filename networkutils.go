package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type cacheValue struct {
	key           string
	flags         string
	expiryTime    int
	amountOfBytes string
	dataBlock     string
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
		if len(words) == 0 {
			continue // Skip if no command is received
		}

		switch words[0] {
		case "set":
			handleSetCommand(conn, words, memcache)
		case "get":
			handleGetCommand(conn, words[len(words)-2], memcache)
		default:
		}
	}
}

func evaluatePort() string {
	portFromUser := flag.String("port", "11211", "make memcached server listen to this port")

	flag.StringVar(portFromUser, "p", "11211", "make memcached server listen to this port (shorthand)")
	flag.Parse()

	return *portFromUser
}
func handleGetCommand(conn net.Conn, key string, memcache map[string]cacheValue) {
	val, ok := memcache[key]
	// If the key exists
	if ok {
		conn.Write([]byte("VALUE " + val.key + " " + val.flags + " " + val.amountOfBytes + "\n"))
		conn.Write([]byte(val.dataBlock + "\n"))
	}
	conn.Write([]byte("END\r\n"))
}
func handleSetCommand(conn net.Conn, words []string, memcache map[string]cacheValue) {
	exptime, _ := strconv.Atoi(words[3])
	value := cacheValue{
		key:           words[1],
		flags:         words[2],
		expiryTime:    exptime,
		amountOfBytes: words[4],
	}

	byteSize, _ := strconv.Atoi(value.amountOfBytes)
	buf := make([]byte, byteSize)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	value.dataBlock = string(buf)
	memcache[value.key] = value
	if words[len(words)-2] != "noreply" {
		_, err := conn.Write([]byte("STORED\r\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}
