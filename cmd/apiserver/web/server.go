package web

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gitgo/apiserver/server"
	"gitgo/apiserver/web/routes"
	"github.com/gorilla/mux"
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
	var commonName = TryExtractCommonName(r.TLS)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(commonName))
	log.Println("Common name: ", commonName)
}

func NewServer(cfg server.Config) (*Server, error) {
	log.Printf("INFO: Creating web server on %s\n", cfg.Address)

	// Listen for requests
	l, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("could not listen for requests on %v", err)
	}

	// Register all routes available on this server
	router := mux.NewRouter()
	router.
		Handle("/api/v1/users", &routes.Users{}).
		Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	router.
		Handle("/api/v1/repositories", &routes.Repositories{}).
		Methods(http.MethodGet, http.MethodPut)
	router.
		Handle("/api/v1/repositories/{name}", &routes.Repository{}).
		Methods(http.MethodGet, http.MethodDelete)
	router.
		Handle("/api/v1/tokens", &routes.Tokens{}).
		Methods(http.MethodGet, http.MethodPost)
	router.NotFoundHandler = &routes.NotFound{}

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
