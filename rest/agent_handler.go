package rest

import (
	"net/http"

	"github.com/applait/xplex-rig/stream"

	"github.com/gorilla/mux"
)

func agentHandler(r *mux.Router) {
	// Routes here can only be accessed by agents
	// TODO: implement agent authentication
	// r.Use(ensureAuthenticatedAgent)

	r.Methods("GET").Path("/config/{streamKey}").HandlerFunc(agentStreamConfig)
}

// agentStreamConfig returns the config for a stream given its stream key.
func agentStreamConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamKey := vars["streamKey"]
	s, err := stream.GetStreamByStreamKey(streamKey)
	if err != nil {
		ErrGetResource.Send(w)
		return
	}
	s.Destinations, err = stream.GetDestinations(s.ID)
	if err != nil {
		ErrGetResource.Send(w)
		return
	}
	if s.IsStreaming {
		ErrResourceInUse.Send(w)
		return
	}
	var res Success
	res.Message = "Stream Config"
	res.Payload = s.Destinations
	res.Send(w)
}
