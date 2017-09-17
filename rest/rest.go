package rest

import (
	"encoding/json"
	"net/http"

	"github.com/applait/xplex-rig/config"
	"github.com/go-pg/pg"
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

// middleware takes in a http.Handler and calls its ServeHTTP method only if it
// can move to the next level
type middleware func(http.Handler) http.Handler

// chain is a sugar for chaining middlewares. It intercepts requests, executes
// mounted middleware handlers in sequence and calls given router.
//
// A middleware definition and body will look like this:
//
// ````
// func exampleHandler(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("Before call")
// 		defer log.Printf("After call")
// 		h.ServeHTTP(w, r) // This will invoke the next handler. Do a premature return to break stack
// 	})
// }
// ````
type chain struct {
	middlewares []middleware
}

// newChain returns a new empty middleware chain
func newChain(m ...middleware) *chain {
	return &chain{
		middlewares: m,
	}
}

// add adds a handler to the middleware chain and returns the `Chain` so that
// you can chain the `Chain`
func (c *chain) add(h ...middleware) *chain {
	c.middlewares = append(c.middlewares, h...)
	return c
}

// handle takes a final http.handler and returns a http.Handler with the entire
// stack set up
func (c chain) handle(h http.Handler) http.Handler {
	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}
	return h
}

// use takes a http.HandlerFunc and wraps it as http.Handler and calls Handle
func (c chain) use(f http.HandlerFunc) http.Handler {
	return c.handle(http.HandlerFunc(f))
}

// Start bootstraps the REST API
func Start(db *pg.DB, config *config.Config) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler).Methods("GET")
	UserHandler(r.PathPrefix("/users").Subrouter(), db, config)
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	success(w, "xplex-rig HTTP API v1", http.StatusOK, []string{
		"GET /",
		"GET /users",
		"GET /streams",
	})
}
