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
	serverPort := EvaluatePort()

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
		go func() {
			err := HandleConnection(conn, memcache)
			if err != nil {
				fmt.Printf("error: %v", err)
			}
		}()
	}
}
