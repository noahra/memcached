package network

import "flag"

func EvaluatePort() string {
	portFromUser := flag.String("port", "11211", "make memcached server listen to this port")
	flag.StringVar(portFromUser, "p", "11211", "make memcached server listen to this port (shorthand)")
	flag.Parse()
	return *portFromUser
}
