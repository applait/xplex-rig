package rest

import (
	"log"
	"net/http"

	"github.com/applait/xplex-rig/config"
	"github.com/applait/xplex-rig/models"
	"github.com/applait/xplex-rig/token"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// StreamHandler providers handler for `/streams` HTTP API
func StreamHandler(r *mux.Router, db *pg.DB, conf *config.Config) {
	r.HandleFunc("/", streamhome).Methods("GET")
	rpost := r.Methods("POST").Subrouter()
	rpost.Handle("/", newChain(auth(conf.Server.JWTSecret, "user")).use(streamCreate(db)))
	rpost.Handle("/key", newChain(required("streamID"), auth(conf.Server.JWTSecret, "user")).use(streamUpdateKey(db)))
}

func streamhome(w http.ResponseWriter, r *http.Request) {
	res := Res{
		Msg:    "Streams API",
		Status: 200,
		Payload: []string{
			"POST /",
			"GET /config",
			"POST /config",
			"POST /key",
			"GET /list",
		},
	}
	res.Send(w)
}

type streamCreateRes struct {
	StreamID  uuid.UUID `json:"streamID"`
	StreamKey string    `json:"streamKey"`
	Active    bool      `json:"active"`
	User      uuid.UUID `json:"user"`
}

// streamCreate creates new multistream entry for user
func streamCreate(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ctxClaims).(*token.Claims)
		m := models.MultiStream{
			UserAccountID: uuid.FromStringOrNil(claims.Issuer),
		}
		if err := m.Create(db); err != nil {
			log.Printf("Error creating multistream. Reason: %s", err)
			errorRes(w, "Error creating multistream", http.StatusBadRequest)
			return
		}
		log.Printf("New multistream created. Stream ID: %s, user ID: %s", m.ID, m.UserAccountID)
		success(w, "Stream created.", http.StatusOK, streamCreateRes{m.ID, m.Key, m.IsActive, m.UserAccountID})
	}
}

type streamUpdateRes struct {
	StreamID  uuid.UUID `json:"streamID"`
	StreamKey string    `json:"streamKey"`
}

// streamUpdate updates multistreaming key for given multistream
func streamUpdateKey(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ctxClaims).(*token.Claims)
		m := models.MultiStream{
			ID:            uuid.FromStringOrNil(r.FormValue("streamID")),
			UserAccountID: uuid.FromStringOrNil(claims.Issuer),
		}
		if err := m.Find(db); err != nil {
			log.Printf("Error updating multistream key. Reason: %s", err)
			errorRes(w, "Invalid stream details.", http.StatusBadRequest)
			return
		}
		if ok, err := m.UpdateKey(db); !ok {
			log.Printf("Error updating multistream key. Reason: %s", err)
			errorRes(w, "Error updating multistream key.", http.StatusBadRequest)
			return
		}
		log.Printf("Updated streamkey for stream %s, user %s", m.ID, m.UserAccountID)
		success(w, "Stream key updated.", http.StatusOK, streamUpdateRes{m.ID, m.Key})
	}
}
