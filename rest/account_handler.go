package rest

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/applait/xplex-rig/common"

	"github.com/applait/xplex-rig/account"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// accountHandler providers handler for `/accounts` HTTP API
func accountHandler(r *mux.Router) {
	// Route for creating new user
	r.Methods("POST").Path("/").HandlerFunc(userCreate)
	// Route for authenticating user using username and password
	r.Methods("POST").Path("/auth/local").HandlerFunc(userAuth)
	// Route for generating new user invite
	r.Methods("POST").Path("/invite/verify").HandlerFunc(userInviteVerify)

	rpost := r.Methods("POST").Subrouter()
	// Ensure all other routes here are authenticated as a user
	rpost.Use(ensureAuthenticatedUser)
	// Route for updating password
	// rpost.Handle("/password", required(userPassword, "oldPassword", "newPassword"))
	rpost.HandleFunc("/password", userPassword)
	// Route for generating new user invite
	rpost.HandleFunc("/invite", userInvite)
}

// userCreateReq defines request data type for user create
type userCreateReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userCreate handles new user creation
func userCreate(w http.ResponseWriter, r *http.Request) {
	var req userCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrInvalidInput.Send(w)
	}
	u := common.UserAccount{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
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

// userPasswordReq defines request data type for user password
type userPasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// userPassword updates a given password for current user in the database.
// Current user is selected from the iss field of JWT used for Authorization
func userPassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	var req userPasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrInvalidInput.Send(w)
		return
	}
	id := uuid.FromStringOrNil(claims.Issuer)
	if err := account.ChangePassword(id, req.OldPassword, req.NewPassword); err != nil {
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

// userAuthReq defines request type for user auth using local strategy
type userAuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// userAuth handles authentication of users using username and password
func userAuth(w http.ResponseWriter, r *http.Request) {
	var req userAuthReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrInvalidInput.Send(w)
		return
	}
	t, err := account.AuthLocal(req.Username, req.Password)
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

// userInviteReq defines request type for generating invite codes
type userInviteReq struct {
	Email string `json:"email"`
}

// userInvite handles generating invite codes
func userInvite(w http.ResponseWriter, r *http.Request) {
	var req userInviteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrInvalidInput.Send(w)
		return
	}
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	e := ErrCreateResource
	_, err := account.GetUserByEmail(req.Email)
	if err == sql.ErrNoRows {
		t, err := account.NewInviteToken(claims.Issuer, req.Email)
		if err != nil {
			log.Printf("Error creating invite token. senderId: %s, email: %s. Reason: %s",
				r.FormValue("senderId"), req.Email, err)
			e.Message = "Unable to create invite"
			e.Status = http.StatusInternalServerError
			e.Send(w)
			return
		}
		var s Success
		s.Message = "Invite created"
		s.Payload = map[string]string{
			"email": req.Email,
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

// userInviteVerifyReq defines request type for verifying invite codes
type userInviteVerifyReq struct {
	Email       string `json:"email"`
	InviteToken string `json:"inviteToken"`
}

// userInviteVerify validates invite token
func userInviteVerify(w http.ResponseWriter, r *http.Request) {
	var req userInviteVerifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrInvalidInput.Send(w)
		return
	}
	t, err := account.ParseUserToken(req.InviteToken)
	if err != nil || t.IssuerType != "invite" || t.Subject != req.Email {
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
