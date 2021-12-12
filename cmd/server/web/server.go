package web

import (
	"fmt"
	"gitgo/server/server"
	"gitgo/server/web/routes"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Config server.Config

	listener net.Listener
	server   *http.Server
	router   *mux.Router
}

func (s *Server) ServeTLS() error {
	if err := s.server.ServeTLS(s.listener, s.Config.CertPath, s.Config.PrivateKey); err != nil {
		return err
	}
	return nil
}

func NewServer(cfg server.Config) (*Server, error) {
	log.Printf("INFO: Creating web server on %s\n", cfg.WebConfig.Address)

	// Listen for requests
	l, err := net.Listen("tcp", cfg.WebConfig.Address)
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
		Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	router.
		Handle("/api/v1/tokens", &routes.Tokens{}).
		Methods(http.MethodGet, http.MethodPost)
	router.NotFoundHandler = &routes.NotFound{}
	s := &http.Server{
		Addr:         cfg.WebConfig.Address,
		ReadTimeout:  cfg.WebConfig.ReadTimeout,
		WriteTimeout: cfg.WebConfig.WriteTimeout,
		IdleTimeout:  cfg.WebConfig.IdleTimeout,
	}
	result := &Server{
		Config:   cfg,
		server:   s,
		listener: l,
	}
	s.Handler = router
	return result, nil
}
