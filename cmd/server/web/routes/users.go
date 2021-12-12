package routes

import (
	"net/http"
)

type Users struct {
}

func (h *Users) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
