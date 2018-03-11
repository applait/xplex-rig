package rest

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/applait/xplex-rig/account"
	"github.com/applait/xplex-rig/common"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// accountHandler providers handler for `/accounts` HTTP API
func accountHandler(r *mux.Router) {
	// Route for creating new user
	r.Methods("POST").Path("/").Handler(required(userCreate, "username", "password", "email"))
	// Route for listing /accounts API info
	r.Methods("GET").Path("/").HandlerFunc(userHome)
	// Route for authenticating user using username and password
	r.Methods("POST").Path("/auth/local").Handler(required(userAuth, "username", "password"))
	// Route for generating new user invite
	r.Methods("POST").Path("/invite/verify").Handler(required(userInviteVerify, "inviteToken", "email"))

	rpost := r.Methods("POST").Subrouter()
	// Ensure all other routes here are authenticated as a user
	rpost.Use(ensureAuthenticatedUser)
	// Route for updating password
	rpost.Handle("/password", required(userPassword, "oldPassword", "newPassword"))
	// Route for generating new user invite
	rpost.Handle("/invite", required(userInvite, "email"))
}

// userCreate handles new user creation
func userCreate(w http.ResponseWriter, r *http.Request) {
	u := common.UserAccount{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	err := account.CreateUser(&u)
	if err != nil {
		errorRes(w, "Unable to create user", http.StatusExpectationFailed)
		return
	}
	log.Printf("User created. ID: %s, Username: %s\n", u.ID, u.Username)

	msg := "User created"
	payload := map[string]string{
		"userID":   u.ID.String(),
		"username": r.FormValue("username"),
		"email":    r.FormValue("email"),
	}
	success(w, msg, http.StatusOK, payload)
}

func userHome(w http.ResponseWriter, r *http.Request) {
	res := Res{
		Msg:    "Users API",
		Status: 200,
		Payload: []string{
			"POST /accounts/ - Create new user account",
			"POST /accounts/password - Update user password",
			"POST /accounts/auth/local - Authenticate using username and password",
			"POST /accounts/invite - Create an invite for a new user account",
			"POST /accounts/invite/verify - Verify an invite using email and invite token",
		},
	}
	res.Send(w)
}

// userPassword updates a given password for current user in the database.
// Current user is selected from the iss field of JWT used for Authorization
func userPassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	id := uuid.FromStringOrNil(claims.Issuer)
	if err := account.ChangePassword(id, r.FormValue("oldPassword"), r.FormValue("newPassword")); err != nil {
		errorRes(w, "Cannot update user password", http.StatusBadRequest)
		return
	}
	success(w, "User password changed", http.StatusOK, map[string]string{
		"userName": claims.Subject,
	})
}

// userAuth handles authentication of users using username and password
func userAuth(w http.ResponseWriter, r *http.Request) {
	t, err := account.AuthLocal(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		errorRes(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	success(w, "Authentication successful", http.StatusOK, map[string]string{
		"username": r.FormValue("username"),
		"token":    t,
	})
}

// userInvite handles generating invite codes
func userInvite(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	_, err := account.GetUserByEmail(r.FormValue("email"))
	if err == sql.ErrNoRows {
		t, err := account.NewInviteToken(claims.Issuer, r.FormValue("email"))
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

// userInviteVerify validates invite token
func userInviteVerify(w http.ResponseWriter, r *http.Request) {
	t, err := account.ParseUserToken(r.FormValue("inviteToken"))
	if err != nil || t.IssuerType != "invite" || t.Subject != r.FormValue("email") {
		errorRes(w, "Error verifying invite token", http.StatusNotAcceptable)
		return
	}
	success(w, "Invite verified.", http.StatusOK, map[string]string{
		"email": t.Subject,
	})
}
