package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/westcoastcode-se/gitgo/api"
	"github.com/westcoastcode-se/gitgo/gitserver/apiserver"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os/exec"
	"path/filepath"
	"sync"
)

var PublicKeyNotFoundError = errors.New("could not find user matching the supplied publicKey")

// Session is an active SSH session
type Session struct {
	context *Context
	cancel  context.CancelFunc

	// Connection that represents a connection to this client
	connection net.Conn

	// hostKey represents the private key
	hostKey ssh.Signer

	// apiServerClient can be used to talk to an api server
	apiServerClient *apiserver.Client

	// User an authorized user if set, nil if no user is found
	User *api.User

	// gitBinDir points to where git binaries are located
	gitBinDir string

	// RepositoryPath points to where repositories are located
	repositoryPath string

	// EnvironmentVars contains all environment variables requested by the client
	// to be used when executing the actual git commands
	environmentVars []string
}

func (s *Session) Close() {
	s.cancel()
}

func (s *Session) HandleConnection() {
	// Prepare configuration for this connection
	sshConfig := ssh.ServerConfig{
		Config: ssh.Config{},
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return s.publicKeyCallback(conn, key)
		},
		ServerVersion: Version,
	}
	sshConfig.AddHostKey(s.hostKey)

	// Before use, a handshake must be performed on the incoming net.Conn.
	sConn, newChannels, reqs, err := ssh.NewServerConn(s.connection, &sshConfig)
	if err != nil {
		log.Printf("WARN: failed to complete handshaking: %v\n", err)
		return
	}
	log.Printf("INFO: new ssh connection from %s(%s)", sConn.RemoteAddr(), sConn.ClientVersion())

	// It's important to "service" requests, otherwise the connection will hang.
	// We only care about requests received over a session channel on this git server
	go ssh.DiscardRequests(reqs)
	go s.processNewChannels(newChannels)
}

func (s *Session) publicKeyCallback(_ ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	fingerprint := ssh.FingerprintSHA256(key)

	// Resolve the user using the public key
	var err error
	s.User, err = s.apiServerClient.FindUserUsingPublicKey(fingerprint)
	if err != nil {
		return nil, err
	}
	if s.User == nil {
		return nil, PublicKeyNotFoundError
	}
	return &ssh.Permissions{}, nil
}

func (s *Session) processNewChannels(newChannels <-chan ssh.NewChannel) {
	// service requests received on a channel
	for newChannel := range newChannels {
		channelType := newChannel.ChannelType()
		if channelType != "session" {
			log.Printf("WARN: unsupported channel type: %s\n", channelType)
			_ = newChannel.Reject(ssh.UnknownChannelType, "unsupported channel type")
			continue
		}
		s.handleNewChannel(newChannel)
	}
}

func (s *Session) handleNewChannel(newChannel ssh.NewChannel) {
	ch, reqs, err := newChannel.Accept()
	if err != nil {
		log.Printf("WARN: could not accept channel %v\n", err)
		return
	}

	go func(in <-chan *ssh.Request) {
		defer ch.Close()
		for req := range in {
			switch req.Type {
			case "env":
				s.processEnvRequest(req)
			case "exec":
				if err = s.processExecRequest(ch, req); err != nil {
					log.Printf("WARN: could not execute incoming request %v\n", err)
				}
				return
			default:
				ch.Write([]byte("Unsupported request type.\r\n"))
				log.Println("ssh: unsupported req type:", req.Type)
				return
			}
		}
	}(reqs)
}

func (s *Session) processEnvRequest(req *ssh.Request) {
	var env struct {
		Name  string
		Value string
	}
	if err := ssh.Unmarshal(req.Payload, &env); err != nil {
		log.Printf("WARN: Invalid env payload %q: %v", req.Payload, err)
		return
	}
	log.Printf("INFO: incoming env request: %q\n", env.Name)

	if !isEnvironmentAllowed(env.Name) {
		log.Printf("INFO: environment variable %s, sent by %s, is not allowed", env.Name,
			s.connection.RemoteAddr())
		return
	}

	s.environmentVars = append(s.environmentVars, env.Name+"="+env.Value)
}

func (s *Session) processExecRequest(ch ssh.Channel, req *ssh.Request) error {
	payload := ReadPayload(req)
	log.Printf("INFO: incoming exec request: %q\n", payload)

	command, err := Parse(payload)
	if err != nil {
		log.Printf("WARN: ignoring %q because it's not a valid git command", payload)
		return err
	}

	if !RepositoryExists(filepath.Join(s.repositoryPath, command.Repository)) {
		return fmt.Errorf("could not find repository %s", command.Repository)
	}

	commandFilename := filepath.Join(s.gitBinDir, command.Command)
	cmd := exec.CommandContext(s.context, commandFilename, command.Repository)
	cmd.Dir = s.repositoryPath
	cmd.Env = s.environmentVars

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("WARN: could not open stdout pipe: %v", err)
		return err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("WARN: could not open stderr pipe: %v", err)
		return err
	}
	defer stderr.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("WARN: could not open stdin pipe: %v", err)
		return err
	}
	defer stdin.Close()

	// We want to wait for stdout and stderr until the command is complete.
	// The command might finish before we've copied all the data to the client, so this helps
	// to ensure that the pipes are not closed beforehand
	wg := &sync.WaitGroup{}
	wg.Add(2)

	if err = cmd.Start(); err != nil {
		log.Printf("WARN: could not start git command: %v", err)
		return err
	}

	go func() {
		defer stdin.Close()
		if _, err := io.Copy(stdin, ch); err != nil {
			log.Printf("WARN: failed to write session to stdin. %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		defer stdout.Close()
		if _, err := io.Copy(ch, stdout); err != nil {
			log.Printf("WARN: failed to write stdout to session. %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		defer stderr.Close()
		if _, err := io.Copy(ch.Stderr(), stderr); err != nil {
			log.Printf("WARN: failed to write stderr to session. %s", err)
		}
	}()

	e := req.Reply(true, nil)
	if e != nil {
		return fmt.Errorf("could not send reply to client: %v", err)
	}

	// Wait for all pipes to finish and then wait to the command to close
	wg.Wait()
	if err = cmd.Wait(); err != nil {
		log.Printf("WARN: command failed: %v", err)
		return err
	}

	_, _ = ch.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
	return nil
}

// Contains environment variables that are allowed by the server
var allowedEnvironmentVariables = []string{"GIT_PROTOCOL"}

// isEnvironmentAllowed verifies that the supplied environment key is allowed
func isEnvironmentAllowed(key string) bool {
	for _, allowed := range allowedEnvironmentVariables {
		if allowed == key {
			return true
		}
	}
	return true
}
