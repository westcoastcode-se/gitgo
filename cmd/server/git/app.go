package git

import (
	"fmt"
	"gitgo/server/server"
	"gitgo/server/user"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
)

const Version = "SSH-2.0-GitGo-1.0.0"

type Server struct {
	// Config contains information on how the server should act
	Config  server.Config
	Version string

	listener net.Listener
	hostKey  ssh.Signer

	// UserDatabase can be used to resolve user-specific properties
	UserDatabase user.Database
}

// Listen for incoming connection requests
func (a *Server) listen() error {
	var err error
	a.listener, err = net.Listen("tcp", a.Config.GitConfig.Address)
	if err != nil {
		return err
	}
	return nil
}

func (a *Server) AcceptClients() error {
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
func (a *Server) HostKey() ssh.Signer {
	return a.hostKey
}

// NewServer creates a new git ssh server
func NewServer(cfg server.Config) (*Server, error) {
	log.Printf("INFO: Creating git server on %s\n", cfg.GitConfig.Address)

	privateBytes, err := ioutil.ReadFile(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %v", err)
	}

	hostKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %v", err)
	}

	result := &Server{
		Config:       cfg,
		Version:      Version,
		hostKey:      hostKey,
		UserDatabase: user.New(),
	}
	err = result.listen()
	if err != nil {
		return nil, fmt.Errorf("could not listen for incoming requests: %v", err)
	}
	return result, nil

}
