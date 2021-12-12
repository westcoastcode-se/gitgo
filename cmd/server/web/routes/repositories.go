package routes

import "net/http"

type Repositories struct {
}

func (h *Repositories) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
