package main

import (
	"log"

	"github.com/2-coffee/Distributed-Logs/server"
	"github.com/2-coffee/Distributed-Logs/web"
)

func main() {
	s := web.NewServer(&server.InMemory{})

	log.Println("Listening for connections")

	s.Serve()
}
