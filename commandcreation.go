package main

import (
	"fmt"
	"net"
	"strconv"
)

func createCommand(words []string, conn net.Conn) (Command, error) {
	commandType := words[0]
	switch commandType {
	case "get":
		return &GetCommand{
			BaseCommand: BaseCommand{
				key:        words[1],
				connection: conn,
			},
		}, nil
	case "set":
		return &SetCommand{
			BaseCommand{
				key:        words[1],
				flags:      words[2],
				expiryTime: parseExpiry(words[3]),
				byteCount:  parseByteCount(words[4]),
				noReply:    words[len(words)-2],
				connection: conn,
			},
		}, nil
	case "add":
		return &AddCommand{
			BaseCommand{
				key:        words[1],
				flags:      words[2],
				expiryTime: parseExpiry(words[3]),
				byteCount:  parseByteCount(words[4]),
				noReply:    words[len(words)-1],
				connection: conn,
			}}, nil
	case "replace":
		return &ReplaceCommand{
			BaseCommand: BaseCommand{
				key:        words[1],
				flags:      words[2],
				expiryTime: parseExpiry(words[3]),
				byteCount:  parseByteCount(words[4]),
				noReply:    words[len(words)-1],
				connection: conn,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", commandType)
	}
}

func parseExpiry(expiryTimeString string) float64 {
	expiryTime, _ := strconv.ParseFloat(expiryTimeString, 64)
	return expiryTime
}
func parseByteCount(byteCountString string) int {
	byteCount, _ := strconv.ParseInt(byteCountString, 10, 32)
	return int(byteCount)
}

func writeNotStored(conn net.Conn) {
	_, err := conn.Write([]byte("NOT_STORED\r\n"))
	if err != nil {
		fmt.Println(err)
	}
}
