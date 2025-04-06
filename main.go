package main

import (
	"flag"
	"log"
	"os"

	"github.com/2-coffee/Distributed-Logs/server"
	"github.com/2-coffee/Distributed-Logs/web"
)

var (
	filename = flag.String("filename", "", "Filename to store our logs data")
	inmem    = flag.Bool("inmem", false, "Whether or not use In-Memory storage. Default: On-Disk")
)

func main() {
	flag.Parse()

	var backend web.Storage

	if *inmem {
		backend = &server.InMemory{} // pointer to in-memory service
	} else {
		// storage setup for on-disk
		if *filename == "" {
			log.Fatalf("The flag `--filename` must be provided")
		}
		fp, err := os.OpenFile(*filename, os.O_CREATE|os.O_RDWR, 0660)
		if err != nil {
			log.Fatalf("Could not create file %q: %s", *filename, err)
		}
		defer fp.Close() // make sure to close connection to file at the end

		backend = server.NewOnDisk(fp) // method to create pointer to on-disk service
	}
	s := web.NewServer(backend)

	log.Println("Listening for connections")

	s.Serve()
}
