package main

import (
	"gitgo/api"
	"gitgo/apiserver/git"
	"gitgo/apiserver/server"
	"gitgo/apiserver/web"
	"log"
)

func main() {
	_ = api.User{
		Name:       "",
		Password:   "",
		PublicKeys: nil,
	}

	log.Println("INFO: Starting GitGo")
	cfg := server.LoadConfig()
	var err error

	// Create a git server
	gitServer, err := git.NewServer(cfg)
	if err != nil {
		log.Fatalf("ERROR: Could not create git server: %v", err)
	}

	webServer, err := web.NewServer(cfg)
	if err != nil {
		log.Fatalf("ERROR: Could not create web server: %v", err)
	}

	go func() {
		err := gitServer.AcceptClients()
		if err != nil {
			log.Fatalf("ERRR: Could not start git server. %e\n", err)
		}
	}()

	err = webServer.ServeTLS()
	if err != nil {
		log.Fatalf("ERROR: Could not start web server. %e\n", err)
	}
	log.Println("INFO: Shutting the server down")
}
