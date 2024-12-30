package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

func evalPING(clientConnection io.ReadWriter, args []string) error {
	if len(args) > 1 {
		return errors.New("invalid args")
	}
	var bytes []byte
	if len(args) == 0 {
		bytes = Encode("PONG")
	} else {
		bytes = Encode(bulkString(args[0]))
	}
	clientConnection.Write(bytes)
	return nil
}

func evalGET(clientConnection io.ReadWriter, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid args")
	}
	key := args[0]
	obj := get(key)
	if obj == nil {
		clientConnection.Write(Encode(nil))
		return nil
	}
	clientConnection.Write(Encode(bulkString(obj.Value.(string))))
	return nil
}

func evalSET(clientConnection io.ReadWriter, args []string) error {
	if len(args) < 2 {
		return errors.New("invalid args")
	}
	key, value := args[0], args[1]
	var expiryMs int64 = -1
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return errors.New("invalid args")
			}
			expirySec, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return errors.New("invalid args")
			}
			expiryMs = expirySec * 1000
		default:
			return errors.New("invalid args")
		}
	}
	put(key, createStoreObj(value, expiryMs))
	clientConnection.Write(Encode("OK"))
	return nil
}

func evalTTL(clientConnection io.ReadWriter, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid args")
	}
	key := args[0]
	obj := get(key)
	if obj == nil {
		clientConnection.Write(Encode(-2))
		return nil
	}
	if obj.ExpiresAt == -1 {
		clientConnection.Write(Encode(-1))
		return nil
	}
	expiryMs := obj.ExpiresAt - time.Now().UnixMilli()
	if expiryMs < 0 {
		clientConnection.Write(Encode(-2))
		return nil
	}
	clientConnection.Write(Encode(expiryMs / 1000))
	return nil
}

func evalDEL(clientConnection io.ReadWriter, args []string) error {
	if len(args) == 0 {
		return errors.New("invalid args")
	}
	deletedCount := 0
	for _, key := range args {
		deletedCount += del(key)
	}
	clientConnection.Write(Encode(deletedCount))
	return nil
}

func evalEXPIRE(clientConnection io.ReadWriter, args []string) error {
	if len(args) < 2 {
		return errors.New("invalid args")
	}
	key := args[0]
	expirySec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("invalid args")
	}
	obj := get(key)
	if obj == nil {
		clientConnection.Write(Encode(0))
		return nil
	}
	obj.ExpiresAt = time.Now().UnixMilli() + expirySec * 1000
	clientConnection.Write(Encode(1))
	return nil
}

func EvalAndRespond(clientConnection io.ReadWriter, cmd *RedisCmd) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(clientConnection, cmd.Args)
	case "GET":
		return evalGET(clientConnection, cmd.Args)
	case "SET":
		return evalSET(clientConnection, cmd.Args)
	case "TTL":
		return evalTTL(clientConnection, cmd.Args)
	case "DEL":
		return evalDEL(clientConnection, cmd.Args)
	case "EXPIRE":
		return evalEXPIRE(clientConnection, cmd.Args)
	default:
		return errors.New("invalid command")
	}
}
