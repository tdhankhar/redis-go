package server

import (
	"log"
	"net"
	"syscall"

	"github.com/tdhankhar/redis-go/config"
)

const maxClients int = 20000

func RunAsyncTcpServer() {
	log.Println("starting async TCP server on", config.Host, config.Port)

	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(serverFd)

	if err = syscall.SetNonblock(serverFd, true); err != nil {
		panic(err)
	}

	ip4 := net.ParseIP(config.Host)
	socketAddress := syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}
	if err = syscall.Bind(serverFd, &socketAddress); err != nil {
		panic(err)
	}

	if err = syscall.Listen(serverFd, maxClients); err != nil {
		panic(err)
	}

	eventLoop := initEventLoop(serverFd)
	eventLoop.Run()
}
