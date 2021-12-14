package server

import (
	"fmt"
	"gitgo/gitserver/apiserver"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
)

// Version represents the SSH server version
const Version = "SSH-2.0-GitGo-1.0.0"

// Server is the actual server instance
type Server struct {
	config   *Config
	listener net.Listener
	hostKey  ssh.Signer

	// apiServerClient is a client that we can use when calling the api server, for example, when
	// checking if a specific fingerprint is allowed to read and write to a specific repository
	apiServerClient *apiserver.Client
}

func (a *Server) HandleConnection(conn net.Conn) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Printf("WARN: failed to process new connection: %v", err)
		_ = conn.Close()
		return
	}

	context, cancel := NewContext(uuid.String())
	session := &Session{
		context:         context,
		cancel:          cancel,
		connection:      conn,
		hostKey:         a.hostKey,
		apiServerClient: a.apiServerClient,
		repositoryPath:  a.config.RepositoriesPath,
		environmentVars: []string{},
	}
	session.HandleConnection()
}

func (a *Server) AcceptClients() error {
	for {
		conn, err := a.listener.Accept()
		if err != nil {
			log.Printf("WARN: could not accept incoming connection. %v\n", err)
			continue
		}
		log.Printf("INFO: new connection established %s\n", conn.RemoteAddr().String())

		// Handle the new connection attempt
		go a.HandleConnection(conn)
	}
}

// NewServer returns a new server
func NewServer(cfg *Config) (*Server, error) {
	apiServerClient, err := apiserver.NewClient(cfg.APIServerAddress, cfg.ClientCertPath, cfg.ClientKeyPath,
		cfg.ClientCAPath, cfg.InsecureSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("could not create api server client: %v", err)
	}

	privateBytes, err := ioutil.ReadFile(cfg.SSHKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %v", err)
	}

	hostKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %v", err)
	}

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("could not listen on address %s: %v", cfg.Address, err)
	}

	s := &Server{
		config:          cfg,
		listener:        listener,
		hostKey:         hostKey,
		apiServerClient: apiServerClient,
	}
	return s, nil
}
