package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/tdhankhar/redis-go/config"
	"github.com/tdhankhar/redis-go/core"
)
func readCommand(clientConnection net.Conn) (*core.RedisCmd, error) {
	buffer := make([]byte, 512)
	bytes, err := clientConnection.Read(buffer)
	if err != nil {
		return nil, err
	}
	tokens, err := core.DecodeArrayString(buffer[:bytes])
	if err != nil {
		return nil, err
	}
	return &core.RedisCmd{
		Cmd: strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respond(clientConnection net.Conn, cmd *core.RedisCmd) {
	err := core.EvalAndRespond(clientConnection, cmd)
	if err != nil {
		clientConnection.Write(core.Encode(err))
	}
}

func RunSyncTcpServer() {
	log.Println("starting sync TCP server on ", config.Host, config.Port)

	concurrentClients := 0
	listener, err := net.Listen("tcp", config.Host + ":" + strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		concurrentClients += 1
		log.Println("client connected with address", clientConnection.RemoteAddr())
		log.Println("concurrent clients", concurrentClients)

		for {
			cmd, err := readCommand(clientConnection)
			if err != nil {
				clientConnection.Close()

				concurrentClients -= 1
				log.Println("client disconnected", clientConnection.RemoteAddr())
				log.Println("concurrent clients", concurrentClients)

				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			respond(clientConnection, cmd)
		}
	}
}