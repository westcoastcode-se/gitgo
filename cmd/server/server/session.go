package server

import (
	"context"
	"errors"
	"fmt"
	"gitgo/server/git"
	"gitgo/server/user"
	"gitgo/server/utils"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os/exec"
	"path/filepath"
	"sync"
)

var PublicKeyNotFoundError = errors.New("could not find user matching the supplied publicKey")

type Session struct {
	App        *App
	Connection net.Conn

	context *Context
	cancel  context.CancelFunc

	// User associated with this session
	user *user.User

	// EnvironmentVars contains all environment variables requested by the client
	// to be used when executing the actual git commands
	EnvironmentVars []string
}

func (s *Session) PublicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	fingerprint := ssh.FingerprintSHA256(key)
	user := s.App.UserDatabase.GetUserUsingPublicKey(fingerprint)
	if user == nil {
		log.Printf("WARN: could not find user with fingerprint %s\n", fingerprint)
		return nil, PublicKeyNotFoundError
	}
	s.user = user
	s.context.SetValue(ContextUser, s.user)
	return &ssh.Permissions{}, nil
}

// HandleConnection does the actual processing of all requests made on the associated connection
func (s *Session) HandleConnection() {
	// Prepare configuration for this connection
	sshConfig := ssh.ServerConfig{
		Config: ssh.Config{},
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return s.PublicKeyCallback(conn, key)
		},
		ServerVersion: s.App.Version,
	}
	sshConfig.AddHostKey(s.App.HostKey())

	// Before use, a handshake must be performed on the incoming net.Conn.
	sConn, newChannels, reqs, err := ssh.NewServerConn(s.Connection, &sshConfig)
	if err != nil {
		log.Printf("WARN: failed to complete handshaking: %v\n", err)
		return
	}
	log.Printf("INFO: new ssh connection from %s(%s)", sConn.RemoteAddr(), sConn.ClientVersion())

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)
	for newChannel := range newChannels {
		channelType := newChannel.ChannelType()
		if channelType != "session" {
			log.Printf("WARN: unsupported channel type: %s\n", channelType)
			_ = newChannel.Reject(ssh.UnknownChannelType, "unsupported channel type")
			continue
		}
		s.HandleNewChannel(newChannel)
	}
}

func (s *Session) HandleNewChannel(newChannel ssh.NewChannel) {
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
		log.Printf("INFO: environment variable %s, sent by %s, is not allowed", env.Name, s.Connection.RemoteAddr())
		return
	}

	s.EnvironmentVars = append(s.EnvironmentVars, env.Name+"="+env.Value)
}

func (s *Session) processExecRequest(ch ssh.Channel, req *ssh.Request) error {
	payload := utils.ReadPayload(req)
	log.Printf("INFO: incoming exec request: %q\n", payload)

	command, err := git.Parse(payload)
	if err != nil {
		log.Printf("WARN: ignoring %q because it's not a valid git command", payload)
		if req.WantReply {
			_ = req.Reply(false, nil)
		}
		return err
	}

	if !git.RepositoryExists(filepath.Join(s.App.Config.RepositoryPath, command.Repository)) {
		log.Printf("WARN: repository not found for payload %q", payload)
		return fmt.Errorf("could not find repository %s", command.Repository)
	}

	cmd := exec.Command(command.Command, command.Repository)
	cmd.Dir = s.App.Config.RepositoryPath
	cmd.Env = s.EnvironmentVars

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
		log.Printf("WARN: command failed: %v", err)
	}

	// Wait for all input and output until we wait for the command to finish
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

// NewSession creates a new session for a specific connection
func NewSession(a *App, conn net.Conn) *Session {
	ctx, cancel := NewContext()
	ctx.SetValue(ContextLocalAddr, conn.LocalAddr())
	ctx.SetValue(ContextRemoteAddr, conn.RemoteAddr())
	return &Session{
		App:        a,
		Connection: conn,
		context:    ctx,
		cancel:     cancel,
	}
}
