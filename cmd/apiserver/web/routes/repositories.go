package routes

import (
	"net/http"
)

type Repositories struct {
}

func (h *Repositories) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

type Repository struct {
}

func (h *Repository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

type RepositoryHasAccess struct {
}

func (h *RepositoryHasAccess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
