package main

import (
	"flag"
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_TYPE = "tcp"
)

func main() {

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
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Received: %s\n", buf)
}

func evaluatePort() string {
	portFromUser := flag.String("port", "11211", "make memcached server listen to this port")

	flag.StringVar(portFromUser, "p", "11211", "make memcached server listen to this port (shorthand)")
	flag.Parse()

	return *portFromUser
}
