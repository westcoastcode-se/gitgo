package routes

import (
	"net/http"
)

type Tokens struct {
}

func (h *Tokens) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
