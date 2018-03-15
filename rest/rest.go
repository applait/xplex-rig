package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Response defines an interface used by the different types of API responses
type Response interface {
	Send(w http.ResponseWriter) (int, error)
}

// Success defines a successful response structure for the HTTP API
type Success struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

// Send sends out a JSON response
func (r Success) Send(w http.ResponseWriter) (int, error) {
	m, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return w.Write(m)
}

// Start bootstraps the REST API
func Start() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(ensureContentType)
	r.NotFoundHandler = notFoundHandler
	r.MethodNotAllowedHandler = methodNotAllowedHandler

	r.HandleFunc("/", homeHandler).Methods("GET")
	accountHandler(r.PathPrefix("/accounts").Subrouter())
	streamHandler(r.PathPrefix("/streams").Subrouter())
	agentHandler(r.PathPrefix("/agents").Subrouter())
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var s Success
	s.Message = "xplex-rig HTTP API v1"
	// TODO: This should send out data based on client access level
	s.Payload = []string{
		"GET / - Get list of api",
		"POST /accounts/ - Create new user account",
		"POST /accounts/password - Update user password",
		"POST /accounts/auth/local - Authenticate using username and password",
		"POST /accounts/invite - Create an invite for a new user account",
		"POST /accounts/invite/verify - Verify an invite using email and invite token",
		"POST /streams/ - Create new stream for current user",
		"GET /streams/ - Get all streams configured for current user",
		"GET /streams/{streamID} - Get details of specific stream",
		"POST /streams/{streamID}/destination - Add a destination for a stream",
		"POST /streams/{streamID}/changeKey - Update streaming key for an existing stream",
		"GET /agents/config/{streamKey} - Get destinations for a given stream key",
	}
	s.Send(w)
}
