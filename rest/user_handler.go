package rest

import (
	"net/http"

	"github.com/applait/xplex-rig/config"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// UserHandler providers handler for `/users` HTTP API
func UserHandler(r *mux.Router, db *pg.DB, conf *config.Config) {
	// Route for creating new user
	r.Handle("/", newChain(required("username", "password", "email")).use(userCreate(db))).Methods("POST")

	// Route for `GET /users/`
	r.HandleFunc("/", userHome).Methods("GET")
}

func userCreate(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Do actual stuff
		o := Res{
			Msg:    "User created",
			Status: http.StatusOK,
			Payload: map[string]string{
				"username": r.FormValue("username"),
				"email":    r.FormValue("email"),
				"token":    "youruniquejwt",
			},
		}
		o.Send(w)
	}
}

func userHome(w http.ResponseWriter, r *http.Request) {
	res := Res{
		Msg:    "Users API",
		Status: 200,
		Payload: []string{
			"POST /",
			"POST /password",
			"POST /auth",
			"POST /invite",
			"GET /invite/verify",
		},
	}
	res.Send(w)
}
