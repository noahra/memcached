package main

import (
	"net"
	"strconv"
)

func handleGetCommand(conn net.Conn, key string, memcache map[string]cacheValue) {
	val, ok := memcache[key]
	if ok {
		if !expiryCheck(val) {
			conn.Write([]byte("VALUE " + val.key + " " + val.flags + " " + strconv.Itoa(val.amountOfBytes) + "\n"))
			conn.Write([]byte(val.dataBlock + "\n"))
		} else {
			delete(memcache, key)
		}
		conn.Write([]byte("END\r\n"))
	}
}
