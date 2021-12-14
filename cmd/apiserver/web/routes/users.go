package routes

import (
	"encoding/json"
	"errors"
	"gitgo/api"
	"log"
	"net/http"
)

type Users struct {
}

func (h *Users) ServeRoute(request *Request) error {
	var fingerprint = request.Query("fingerprint")
	if len(fingerprint) == 0 {
		return errors.New("missing query parameter 'fingerprint'")
	}

	log.Println("Testing ", fingerprint)

	user := api.User{
		Name:     "apiserver_client",
		Password: "asdf",
		PublicKeys: []api.PublicKey{
			{
				Name:        "OSX",
				Fingerprint: fingerprint,
				PublicKey:   "???",
			},
		},
	}
	bytes, _ := json.Marshal(user)
	_, _ = request.Ok(bytes)
	return nil
}

func (h *Users) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
