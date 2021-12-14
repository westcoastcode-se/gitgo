package routes

import (
	"github.com/westcoastcode-se/gitgo/apiserver/server"
	"github.com/westcoastcode-se/gitgo/apiserver/user"
	"net/http"
)

type Request struct {
	Context  *server.Context
	User     *user.User
	Original *http.Request
	Response http.ResponseWriter
	URI      string
}

func (r *Request) Query(param string) string {
	return r.Original.URL.Query().Get(param)
}

func (r *Request) Ok(body []byte) (int, error) {
	r.Response.WriteHeader(http.StatusOK)
	r.Response.Header().Set("Content-Type", "application/json")
	return r.Response.Write(body)
}

func FromHttpRequest(rw http.ResponseWriter, r *http.Request) *Request {
	ctx, _ := server.NewContext()
	return &Request{
		Context:  ctx,
		User:     nil,
		Original: r,
		Response: rw,
		URI:      r.RequestURI,
	}
}

type Route interface {
	// ServeRoute serves a specific request for this route
	ServeRoute(request *Request) error
}
