package rest

import (
	"encoding/json"
	"net/http"

	"github.com/applait/xplex-rig/server/config"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// Res defines response structure
type Res struct {
	Msg     string      `json:"msg"`
	Status  int         `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

// Send sends out a JSON response
func (r Res) Send(w http.ResponseWriter) (int, error) {
	m, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	return w.Write(m)
}

// Start bootstraps the REST API
func Start(db *pg.DB, config *config.Config) *mux.Router {
	r := mux.NewRouter()
	users(r.PathPrefix("/users").Subrouter().StrictSlash(true))
	r.HandleFunc("/", home)
	return r
}

func users(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Users!\n"))
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	o := Res{
		Msg:    "xplex-rig HTTP API v1",
		Status: 200,
		Payload: []string{
			"GET /",
			"GET /users",
			"GET /streams",
		},
	}
	o.Send(w)
}
