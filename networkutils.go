package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type cacheValue struct {
	key           string
	flags         string
	expiryTime    int
	amountOfBytes string
	dataBlock     string
	createdAt     time.Time
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
		case "add":
			handleAddCommand(conn, words, memcache)
		case "replace":
			handleReplaceCommand(conn, words, memcache)
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
		hasExpired := false
		if val.expiryTime != 0 {
			duration := time.Now().Sub(val.createdAt)
			if int(duration.Seconds()) > val.expiryTime {
				hasExpired = true
			}
		}
		if !hasExpired {
			conn.Write([]byte("VALUE " + val.key + " " + val.flags + " " + val.amountOfBytes + "\n"))
			conn.Write([]byte(val.dataBlock + "\n"))
		} else {
			delete(memcache, key)
		}
		conn.Write([]byte("END\r\n"))
	}
}

func handleSetCommand(conn net.Conn, words []string, memcache map[string]cacheValue) {
	exptime, _ := strconv.Atoi(words[3])
	value := cacheValue{
		key:           words[1],
		flags:         words[2],
		expiryTime:    exptime,
		amountOfBytes: words[4],
		createdAt:     time.Now(),
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

func handleReplaceCommand(conn net.Conn, words []string, memcache map[string]cacheValue) {
	exptime, _ := strconv.Atoi(words[3])
	value := cacheValue{
		key:           words[1],
		flags:         words[2],
		expiryTime:    exptime,
		amountOfBytes: words[4],
		createdAt:     time.Now(),
	}

	byteSize, _ := strconv.Atoi(value.amountOfBytes)
	buf := make([]byte, byteSize)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	value.dataBlock = string(buf)
	_, ok := memcache[words[1]]
	if !ok {
		_, err := conn.Write([]byte("NOT_STORED\r\n"))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		memcache[value.key] = value
		if words[len(words)-2] != "noreply" {
			_, err := conn.Write([]byte("STORED\r\n"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

func handleAddCommand(conn net.Conn, words []string, memcache map[string]cacheValue) {
	exptime, _ := strconv.Atoi(words[3])
	value := cacheValue{
		key:           words[1],
		flags:         words[2],
		expiryTime:    exptime,
		amountOfBytes: words[4],
		createdAt:     time.Now(),
	}

	byteSize, _ := strconv.Atoi(value.amountOfBytes)
	buf := make([]byte, byteSize)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	value.dataBlock = string(buf)
	_, ok := memcache[words[1]]
	if ok {
		_, err := conn.Write([]byte("NOT_STORED\r\n"))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		memcache[value.key] = value
		if words[len(words)-2] != "noreply" {
			_, err := conn.Write([]byte("STORED\r\n"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
