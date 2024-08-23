package main

import (
	"fmt"
	"memcached/pkg/cache"
	"memcached/pkg/network"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_TYPE = "tcp"
)

func main() {
	memcache := cache.NewCache()
	serverPort := network.EvaluatePort()

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
			err := network.HandleConnection(conn, memcache)
			if err != nil {
				fmt.Printf("error: %w", err)
			}
		}()
	}
}
