package web

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/westcoastcode-se/gitgo/apiserver/server"
	"github.com/westcoastcode-se/gitgo/apiserver/user"
	"github.com/westcoastcode-se/gitgo/apiserver/web/routes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Config server.Config

	listener net.Listener
	server   *http.Server
}

func (s *Server) ServeTLS() error {
	if err := s.server.ServeTLS(s.listener, s.Config.CertPath, s.Config.PrivateKey); err != nil {
		return err
	}
	return nil
}

func (s Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	request := routes.FromHttpRequest(rw, r)
	var commonName = TryExtractCommonName(r.TLS)
	request.User = &user.User{
		Name:       commonName,
		Password:   "",
		PublicKeys: nil,
	}
	var usersRoute = &routes.Users{}
	usersRoute.ServeRoute(request)
}

func NewServer(cfg server.Config) (*Server, error) {
	log.Printf("INFO: Creating web server on %s\n", cfg.Address)

	// Listen for requests
	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("could not listen for requests on %v", err)
	}

	log.Println("INFO: Reading CA user to verify client-side certificates")
	cert, err := ioutil.ReadFile(cfg.CAPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	s := &http.Server{
		Addr:         cfg.Address,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		TLSConfig: &tls.Config{
			ClientAuth: tls.VerifyClientCertIfGiven,
			ClientCAs:  caCertPool,
		},
	}
	result := &Server{
		Config:   cfg,
		server:   s,
		listener: l,
	}
	s.Handler = result
	return result, nil
}
