package main

import (
	ds "applait/xplex-rig/server/datastore"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HomeHandler handles home URL route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

func main() {
	connuri := "postgres://xplex:1234@localhost/xplex_rig_dev?sslmode=disable"
	_, err := ds.ConnectPG(connuri)
	if err != nil {
		log.Fatalf("Error connecting to DB. Reason: %s", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8081", r)
}
