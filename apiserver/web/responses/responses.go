package responses

import (
	"encoding/json"
	"github.com/westcoastcode-se/gitgo/api"
	"net/http"
)

type RequestError interface {
	Reason() string
	StatusCode() int
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Reason() string {
	return "Not Found"
}

func (e *NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// WriteError writes the supplied error
func WriteError(uri string, err RequestError, rw http.ResponseWriter) (int, error) {
	bytes, _ := json.Marshal(&api.Error{URI: uri, Reason: err.Reason()})
	rw.WriteHeader(err.StatusCode())
	return rw.Write(bytes)
}
