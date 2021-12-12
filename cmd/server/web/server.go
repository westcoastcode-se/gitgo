package web

import (
	"gitgo/server/server"
	"gitgo/server/web/routes"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	Config server.Config

	server *http.Server
	router *mux.Router
}

func (s *Server) ListenAndServe() error {
	if err := s.server.ListenAndServeTLS(s.Config.CertPath, s.Config.PrivateKey); err != nil {
		return err
	}
	return nil
}

func NewServer(cfg server.Config) *Server {
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
		ReadTimeout:  cfg.WebConfig.ReadTimeout * time.Millisecond,
		WriteTimeout: cfg.WebConfig.WriteTimeout * time.Millisecond,
		IdleTimeout:  cfg.WebConfig.IdleTimeout * time.Millisecond,
	}
	result := &Server{
		Config: cfg,
		server: s,
	}
	s.Handler = router
	return result
}
