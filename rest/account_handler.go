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
		e := ErrCreateResource
		e.Message = "Cannot create user account."
		e.Send(w)
		return
	}
	log.Printf("User created. ID: %s, Username: %s\n", u.ID, u.Username)
	var s Success
	s.Message = "User created"
	s.Payload = map[string]string{
		"userID":   u.ID.String(),
		"username": r.FormValue("username"),
		"email":    r.FormValue("email"),
	}
	s.Send(w)
}

// userPassword updates a given password for current user in the database.
// Current user is selected from the iss field of JWT used for Authorization
func userPassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	id := uuid.FromStringOrNil(claims.Issuer)
	if err := account.ChangePassword(id, r.FormValue("oldPassword"), r.FormValue("newPassword")); err != nil {
		e := ErrUpdateResource
		e.Message = "Cannot update user password"
		e.Send(w)
		return
	}
	var s Success
	s.Message = "User password changed"
	s.Payload = map[string]string{
		"userName": claims.Subject,
	}
	s.Send(w)
}

// userAuth handles authentication of users using username and password
func userAuth(w http.ResponseWriter, r *http.Request) {
	t, err := account.AuthLocal(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		ErrInvalidCredentials.Send(w)
		return
	}
	var s Success
	s.Message = "Authentication successful"
	s.Payload = map[string]string{
		"username": r.FormValue("username"),
		"token":    t,
	}
	s.Send(w)
}

// userInvite handles generating invite codes
func userInvite(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	e := ErrCreateResource
	_, err := account.GetUserByEmail(r.FormValue("email"))
	if err == sql.ErrNoRows {
		t, err := account.NewInviteToken(claims.Issuer, r.FormValue("email"))
		if err != nil {
			log.Printf("Error creating invite token. senderId: %s, email: %s. Reason: %s",
				r.FormValue("senderId"), r.FormValue("email"), err)
			e.Message = "Unable to create invite"
			e.Status = http.StatusInternalServerError
			e.Send(w)
			return
		}
		var s Success
		s.Message = "Invite created"
		s.Payload = map[string]string{
			"email": r.FormValue("email"),
			"token": t,
		}
		s.Send(w)
		return
	}
	if err == nil {
		e.Message = "Email is already registered"
		e.Send(w)
		return
	}
	e.Send(w)
}

// userInviteVerify validates invite token
func userInviteVerify(w http.ResponseWriter, r *http.Request) {
	t, err := account.ParseUserToken(r.FormValue("inviteToken"))
	if err != nil || t.IssuerType != "invite" || t.Subject != r.FormValue("email") {
		e := ErrInvalidInput
		e.Message = "Error verifying invite token"
		e.Send(w)
		return
	}
	var s Success
	s.Message = "Invite verified."
	s.Payload = map[string]string{
		"email": t.Subject,
	}
	s.Send(w)
}
