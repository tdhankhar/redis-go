package core

import (
	"errors"
	"net"
)

func evalPING(clientConnection net.Conn, args []string) error {
	var bytes []byte
	if len(args) > 1 {
		return errors.New("invalid args")
	}
	if len(args) == 0 {
		bytes = Encode("PONG")
	} else {
		bytes = Encode(BulkString(args[0]))
	}
	clientConnection.Write(bytes)
	return nil
}

func EvalAndRespond(clientConnection net.Conn, cmd *RedisCmd) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(clientConnection, cmd.Args)
	default:
		return errors.New("invalid command")
	}
}