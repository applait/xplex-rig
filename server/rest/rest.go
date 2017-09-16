package rest

import (
	"net/http"

	"github.com/applait/xplex-rig/server/config"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// Start bootstraps the REST API
func Start(db *pg.DB, config *config.Config) *mux.Router {
	r := mux.NewRouter()
	users(r.PathPrefix("/users").Subrouter())
	return r
}

func users(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Users!\n"))
	})
}
