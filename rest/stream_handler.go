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

// StreamHandler providers handler for `/streams` HTTP API
func StreamHandler(r *mux.Router, db *pg.DB, conf *config.Config) {
	r.HandleFunc("/", streamhome).Methods("GET")
	rpost := r.Methods("POST").Subrouter()
	rpost.Handle("/create", newChain(auth(conf.Server.JWTSecret, "user")).use(streamCreate(db)))
}

func streamhome(w http.ResponseWriter, r *http.Request) {
	res := Res{
		Msg:    "Streams API",
		Status: 200,
		Payload: []string{
			"POST /create",
			"POST /config/new",
			"GET /config",
			"POST /updateKey",
			"GET /list",
		},
	}
	res.Send(w)
}

type streamCreateRes struct {
	streamID  int
	streamKey string
	active    bool
	user      int
}

// streamCreate creates new multistream entry for user
func streamCreate(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ctxClaims).(*token.Claims)
		uid, _ := strconv.Atoi(claims.Issuer)
		m := models.MultiStream{
			UserID: uid,
		}
		if err := m.Create(db); err != nil {
			log.Printf("Error creating multistream. Reason: %s", err)
			errorRes(w, "Error creating multistream", http.StatusBadRequest)
			return
		}
		log.Printf("New multistream created for user %d", uid)
		success(w, "Stream created.", http.StatusOK, streamCreateRes{m.ID, m.Key, m.IsActive, m.UserID})
	}
}
