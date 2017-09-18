package rest

import (
	"log"
	"net/http"
	"strconv"

	"github.com/applait/xplex-rig/config"
	"github.com/applait/xplex-rig/models"
	"github.com/applait/xplex-rig/token"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
)

// UserHandler providers handler for `/users` HTTP API
func UserHandler(r *mux.Router, db *pg.DB, conf *config.Config) {
	// Route for `GET /users/`
	r.HandleFunc("/", userHome).Methods("GET")

	rpost := r.Methods("POST").Subrouter()
	// Route for updating password
	rpost.Handle("/password", newChain(required("password"), auth(conf.Server.JWTSecret, "user")).
		use(userPassword(db)))

	// Route for authenticating user using username and password
	rpost.Handle("/auth", newChain(required("username", "password")).use(userAuth(db, conf)))

	// Route for generating new user invite
	rpost.Handle("/invite", newChain(required("email"), auth(conf.Server.JWTSecret, "user")).
		use(userInvite(db, conf)))

	// Route for creating new user
	rpost.Handle("/", newChain(required("username", "password", "email")).use(userCreate(db, conf)))
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
		if t, err = token.NewUserToken(&u, conf.Server.JWTSecret); err != nil {
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

// userPassword updates a given password for current user in the database.
// Current user is selected from the iss field of JWT used for Authorization
func userPassword(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ctxClaims).(*token.Claims)
		userID, err := strconv.Atoi(claims.Issuer)
		if err != nil {
			log.Printf("Error converting user claim issuer. Reason: %s", err)
			errorRes(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}
		u := models.User{
			ID:       userID,
			Username: claims.Subject,
		}
		if err = u.Find(db); err != nil {
			errorRes(w, "Error updating user password.", http.StatusInternalServerError)
			return
		}
		if err = u.UpdatePassword(db, r.FormValue("password")); err != nil {
			log.Printf("Error hashing and setting new user password. Reason: %s", err)
			errorRes(w, "Error updating user password.", http.StatusInternalServerError)
			return
		}
		success(w, "User password changed", http.StatusOK, map[string]string{
			"userName": claims.Subject,
		})
	}
}

// userAuth handles authentication of users using username and password
func userAuth(db *pg.DB, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := models.User{
			Username: r.FormValue("username"),
		}
		err := u.Find(db)
		if err != nil && err != pg.ErrNoRows {
			log.Printf("Error retrieving user information. Reason: %s", err)
			errorRes(w, "Error fetching user information", http.StatusInternalServerError)
			return
		}
		if err == pg.ErrNoRows || u.MatchPassword(r.FormValue("password")) == false {
			errorRes(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		t, err := token.NewUserToken(&u, conf.Server.JWTSecret)
		if err != nil {
			log.Printf("Error creating auth token for user ID %d. Reason: %s\n", u.ID, err)
			errorRes(w, "Unable to create auth token", http.StatusInternalServerError)
			return
		}
		success(w, "Authentication successful", http.StatusOK, map[string]string{
			"token": t,
		})
	}
}

// userInvite handles generating invite codes
func userInvite(db *pg.DB, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := models.User{Email: r.FormValue("email")}
		claims := r.Context().Value(ctxClaims).(*token.Claims)
		err := u.Find(db)
		if err == pg.ErrNoRows {
			t, err := token.NewInviteToken(claims.Issuer, r.FormValue("email"), conf.Server.JWTSecret)
			if err != nil {
				log.Printf("Error creating invite token. senderId: %s, email: %s. Reason: %s",
					r.FormValue("senderId"), r.FormValue("email"), err)
				errorRes(w, "Unable to create invite", http.StatusInternalServerError)
			}
			success(w, "Invite created", http.StatusOK, map[string]string{
				"email": r.FormValue("email"),
				"token": t,
			})
			return
		}
		if err == nil {
			errorRes(w, "Email is already registered.", http.StatusPreconditionFailed)
			return
		}
		errorRes(w, "Unable to create invite.", http.StatusInternalServerError)
	}
}
