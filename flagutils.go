package main

import "flag"

func evaluatePort() string {
	portFromUser := flag.String("port", "11211", "make memcached server listen to this port")
	flag.StringVar(portFromUser, "p", "11211", "make memcached server listen to this port (shorthand)")
	flag.Parse()
	return *portFromUser
}
