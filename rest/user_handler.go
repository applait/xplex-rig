package rest

import (
	"log"
	"net/http"

	"github.com/applait/xplex-rig/config"
	"github.com/applait/xplex-rig/models"
	"github.com/applait/xplex-rig/token"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// UserHandler providers handler for `/users` HTTP API
func UserHandler(r *mux.Router, db *pg.DB, conf *config.Config) {
	// Route for creating new user
	r.Handle("/", newChain(required("username", "password", "email")).use(userCreate(db, conf))).Methods("POST")

	// Route for `GET /users/`
	r.HandleFunc("/", userHome).Methods("GET")
}

// userCreate handles new user creation
func userCreate(db *pg.DB, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := models.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
		}
		if err := u.SetPassword(r.FormValue("password")); err != nil {
			log.Printf("Error setting user password. Reason: %s\n", err)
			errorRes(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		if err := u.Insert(db); err != nil {
			log.Printf("Error saving new user to DB. Reason: %s\n", err)
			errorRes(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		log.Printf("User created. ID: %d, Username: %s\n", u.ID, u.Username)

		var t string
		var err error

		msg := "User created"
		payload := map[string]string{
			"username": r.FormValue("username"),
			"email":    r.FormValue("email"),
			"token":    "",
		}
		if t, err = token.NewUserToken(u.ID, conf.Server.JWTSecret); err != nil {
			msg = "User created, but could not generate token"
			log.Printf("Error creating token for new user ID %d. Reason: %s\n", u.ID, err)
		} else {
			log.Printf("Generated token for new user ID %d\n", u.ID)
		}
		payload["token"] = t
		success(w, msg, http.StatusOK, payload)
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
