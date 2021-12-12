package routes

import (
	"gitgo/server/web/responses"
	"net/http"
)

// NotFound is a route which is called when the requested URI aren't found
type NotFound struct {
}

func (h *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = responses.WriteError(r.RequestURI, &responses.NotFoundError{}, w)
}
