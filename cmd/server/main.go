package main

import (
	"gitgo/api"
	"gitgo/server/server"
	"log"
)

func main() {
	_ = api.User{
		Name:       "",
		Password:   "",
		PublicKeys: nil,
	}

	log.Println("Starting server")
	cfg := server.LoadConfig()
	a := server.NewApp(cfg)
	err := a.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start server. %e\n", err)
	}
	log.Println("Shutting the server down")
}
