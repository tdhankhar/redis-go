package main

import (
	"flag"
	"log"

	"github.com/tdhankhar/redis-go/config"
	"github.com/tdhankhar/redis-go/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for redis-go server")
	flag.IntVar(&config.Port, "port", 6379, "port for redis-go server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("redis-go server started")
	server.RunSyncTcpServer()
}