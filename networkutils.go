package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type cacheValue struct {
	flags         string
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

		if words[0] == "set" {
			fmt.Println("reached")
			handleSetCommand(conn, words, memcache)
		}
	}
}

func evaluatePort() string {
	portFromUser := flag.String("port", "11211", "make memcached server listen to this port")

	flag.StringVar(portFromUser, "p", "11211", "make memcached server listen to this port (shorthand)")
	flag.Parse()

	return *portFromUser
}

func handleSetCommand(conn net.Conn, words []string, memcache map[string]cacheValue) {
	key := words[1]
	value := cacheValue{
		flags:         words[2],
		amountOfBytes: words[3],
	}

	byteSize, _ := strconv.Atoi(value.amountOfBytes)
	buf := make([]byte, byteSize)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	value.dataBlock = string(buf)
	memcache[key] = value
	fmt.Printf("Memcache: %s\n", memcache)
	if words[len(words)-2] != "noreply" {
		fmt.Println("STORED")
		conn.Write([]byte("STORED"))
	}

}
