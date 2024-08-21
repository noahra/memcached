package main

import (
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_TYPE = "tcp"
)

func main() {
	memcache := NewCache()
	serverPort := evaluatePort()

	ln, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+serverPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn, memcache)
	}
}
