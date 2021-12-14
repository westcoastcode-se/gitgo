package main

import (
	"gitgo/api"
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

	webServer, err := web.NewServer(cfg)
	if err != nil {
		log.Fatalf("ERROR: Could not create web server: %v", err)
	}

	err = webServer.ServeTLS()
	if err != nil {
		log.Fatalf("ERROR: Could not start web server. %e\n", err)
	}
	log.Println("INFO: Shutting the server down")
}
