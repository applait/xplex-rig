package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Res defines response structure
type Res struct {
	Msg     string      `json:"msg"`
	Status  int         `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

// Send sends out a JSON response
func (r Res) Send(w http.ResponseWriter) (int, error) {
	m, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	return w.Write(m)
}

// errorRes is a shorthand for sending error response
func errorRes(w http.ResponseWriter, msg string, status int) (int, error) {
	res := Res{
		Msg:    msg,
		Status: status,
	}
	return res.Send(w)
}

// success is a shorthand for sending success response
func success(w http.ResponseWriter, msg string, status int, payload interface{}) (int, error) {
	res := Res{
		Msg:     msg,
		Status:  status,
		Payload: payload,
	}
	return res.Send(w)
}

// Start bootstraps the REST API
func Start() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler).Methods("GET")
	accountHandler(r.PathPrefix("/accounts").Subrouter())
	// StreamHandler(r.PathPrefix("/streams").Subrouter())
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	success(w, "xplex-rig HTTP API v1", http.StatusOK, []string{
		"GET /",
		"GET /users",
		// "GET /streams",
	})
}
