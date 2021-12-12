package server

import (
	"gitgo/server/user"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
)

const Version = "SSH-2.0-GitGo-1.0.0"

type App struct {
	// Config contains information on how the server should act
	Config  Config
	Version string

	listener net.Listener
	hostKey  ssh.Signer

	// UserDatabase can be used to resolve user-specific properties
	UserDatabase user.Database
}

func (a *App) ListenAndServe() error {
	log.Printf("Listening for requests on %s\n", a.Config.Address)

	var err error
	a.listener, err = net.Listen("tcp", a.Config.Address)
	if err != nil {
		return nil
	}

	for {
		conn, err := a.listener.Accept()
		if err != nil {
			log.Printf("ERROR: could not accept incoming connection. %e\n", err)
			continue
		}
		log.Printf("INFO: new connection established %s\n", conn.RemoteAddr().String())

		// Create a new session for this connection and handle it
		session := NewSession(a, conn)
		go session.HandleConnection()
	}
}

// HostKey can be used to fetch the private key used by the ssh server
func (a *App) HostKey() ssh.Signer {
	return a.hostKey
}

func NewApp(cfg Config) *App {
	privateBytes, err := ioutil.ReadFile(cfg.PrivateKey)
	if err != nil {
		log.Fatalf("could not read private key: %v", err)
	}

	hostKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("could not parse private key: %v", err)
	}

	result := &App{
		Config:       cfg,
		Version:      Version,
		hostKey:      hostKey,
		UserDatabase: user.New(),
	}
	return result
}
