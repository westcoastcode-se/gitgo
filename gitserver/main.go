package main

import (
	"github.com/westcoastcode-se/gitgo/gitserver/server"
	"log"
)

func main() {
	log.Println("INFO: Starting git server")

	cfg := server.LoadConfig()
	s, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("ERROR: Could not create a new git server: %v", err)
	}

	err = s.AcceptClients()
	if err != nil {
		log.Fatalf("ERRR: Could not start git server. %v", err)
	}

	log.Println("INFO: Shutting the server down")
}
