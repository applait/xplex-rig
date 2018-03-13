package rest

import (
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"
	validator "gopkg.in/validator.v2"

	"github.com/applait/xplex-rig/account"
	"github.com/applait/xplex-rig/stream"
	"github.com/gorilla/mux"
)

func streamHandler(r *mux.Router) {
	// All routes here are authenticated
	r.Use(ensureAuthenticatedUser)

	r.Methods("POST").Path("/").HandlerFunc(streamCreate)
	r.Methods("GET").Path("/").HandlerFunc(streamList) // All streams of current user
	r.Methods("GET").Path("/{streamID}").HandlerFunc(streamDetails)
	r.Methods("POST").Path("/{streamID}/changeKey").HandlerFunc(streamChangeKey)
	r.Methods("POST").Path("/{streamID}/destination").HandlerFunc(streamAddDestination)
}

// streamCreate creates a new stream
func streamCreate(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	s, err := stream.CreateStream(uuid.FromStringOrNil(claims.Issuer))
	if err != nil {
		e := ErrInvalidInput
		e.Send(w)
		return
	}
	var res Success
	res.Message = "New stream created"
	res.Payload = s
	res.Send(w)
}

// streamList lists all streams of current user
func streamList(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	s, err := stream.GetStreamsOfUser(uuid.FromStringOrNil(claims.Issuer))
	if err != nil {
		ErrGetResource.Send(w)
		return
	}
	var res Success
	res.Message = "Stream list for user"
	res.Payload = s
	res.Send(w)
	return
}

// streamDetails shows detail of a single stream
func streamDetails(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	userID := uuid.FromStringOrNil(claims.Issuer)
	vars := mux.Vars(r)
	streamID := uuid.FromStringOrNil(vars["streamID"])
	s, err := stream.GetStreamByID(streamID)
	if err != nil {
		ErrGetResource.Send(w)
		return
	}
	if !uuid.Equal(userID, s.User.ID) {
		ErrInvalidInput.Send(w)
		return
	}
	s.Destinations, err = stream.GetDestinations(streamID)
	if err != nil {
		ErrGetResource.Send(w)
	}
	var res Success
	res.Message = "Stream detail"
	res.Payload = s
	res.Send(w)
}

// streamChangeKey changes stream key for a given stream
func streamChangeKey(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	userID := uuid.FromStringOrNil(claims.Issuer)
	vars := mux.Vars(r)
	streamID := uuid.FromStringOrNil(vars["streamID"])
	s, err := stream.GetStreamByID(streamID)
	if err != nil || !uuid.Equal(userID, s.User.ID) {
		ErrInvalidInput.Send(w)
		return
	}
	k, err := stream.ChangeStreamKey(streamID, userID)
	if err != nil {
		e := ErrUpdateResource
		e.Message = "Unable to update stream key"
		e.Send(w)
		return
	}
	var res Success
	res.Message = "New stream key"
	res.Payload = map[string]string{
		"key":      k,
		"streamID": streamID.String(),
	}
	res.Send(w)
}

type streamAddDestinationReq struct {
	Service   string `json:"service" validate:"nonzero"`
	StreamKey string `json:"streamKey" validate:"nonzero"`
}

// streamAddDestination creates a new destination for a given stream
func streamAddDestination(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ctxClaims).(*account.Claims)
	userID := uuid.FromStringOrNil(claims.Issuer)
	vars := mux.Vars(r)
	streamID := uuid.FromStringOrNil(vars["streamID"])
	s, err := stream.GetStreamByID(streamID)
	if err != nil || !uuid.Equal(userID, s.User.ID) {
		ErrInvalidInput.Send(w)
		return
	}
	var req streamAddDestinationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrInvalidInput.Send(w)
		return
	}
	if err := validator.Validate(req); err != nil {
		e := ErrInvalidInput
		e.Details = err
		e.Send(w)
		return
	}
	dest, err := stream.AddDestination(streamID, req.Service, req.StreamKey, userID)
	if err != nil {
		e := ErrCreateResource
		e.Message = "Unable to add destination to stream"
		e.Send(w)
		return
	}
	var res Success
	res.Message = "Stream destination added"
	res.Payload = dest
	res.Send(w)
}
