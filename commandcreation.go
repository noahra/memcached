package main

import (
	"fmt"
	"net"
	"strconv"
)

func CreateCommand(words []string, conn net.Conn) (Command, error) {
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
		baseCommand, err := createBaseCommand(words, conn)
		if err != nil {
			return nil, err
		}
		return &SetCommand{
			baseCommand,
		}, nil
	case "add":
		baseCommand, err := createBaseCommand(words, conn)
		if err != nil {
			return nil, err
		}
		return &AddCommand{
			baseCommand,
		}, nil
	case "append":
		baseCommand, err := createBaseCommand(words, conn)
		if err != nil {
			return nil, err
		}
		return &AppendCommand{
			baseCommand,
		}, nil
	case "prepend":
		baseCommand, err := createBaseCommand(words, conn)
		if err != nil {
			return nil, err
		}
		return &PrependCommand{
			baseCommand,
		}, nil
	case "replace":
		baseCommand, err := createBaseCommand(words, conn)
		if err != nil {
			return nil, err
		}
		return &ReplaceCommand{
			baseCommand,
		}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", commandType)
	}
}
func writeNotStored(conn net.Conn) {
	_, err := conn.Write([]byte("NOT_STORED\r\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func createBaseCommand(words []string, conn net.Conn) (BaseCommand, error) {
	flags, err := strconv.ParseInt(words[2], 10, 32)
	if err != nil {
		return BaseCommand{}, fmt.Errorf("could not parse flags from command: %w", err)
	}
	i, err := strconv.Atoi(words[3])
	if err != nil {
		return BaseCommand{}, fmt.Errorf("could not parse int from command: %w", err)
	}
	return BaseCommand{
		key:        words[1],
		flags:      flags,
		expiryTime: i,
		byteCount:  parseByteCount(words[4]),
		noReply:    words[len(words)-1],
		connection: conn,
	}, nil
}
