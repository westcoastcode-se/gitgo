package main

import (
	"gitgo/api"
	"gitgo/server/git"
	"gitgo/server/server"
	"gitgo/server/web"
	"log"
)

func main() {
	_ = api.User{
		Name:       "",
		Password:   "",
		PublicKeys: nil,
	}

	log.Println("INFO: Starting server")
	cfg := server.LoadConfig()
	gitServer := git.NewServer(cfg)
	webServer := web.NewServer(cfg)

	go func() {
		err := gitServer.ListenAndServe()
		if err != nil {
			log.Fatalf("ERRR: Could not start git server. %e\n", err)
		}
	}()

	err := webServer.ListenAndServe()
	if err != nil {
		log.Fatalf("ERROR: Could not start web server. %e\n", err)
	}
	log.Println("INFO: Shutting the server down")
}
