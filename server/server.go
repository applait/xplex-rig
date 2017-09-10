package main

import (
	"log"
	"net/http"

	"github.com/applait/xplex-rig/server/datastore"
	"github.com/gorilla/mux"
)

// HomeHandler handles home URL route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

func main() {
	connuri := "postgres://xplex:1234@localhost/xplex_rig_dev?sslmode=disable"
	db, err := datastore.ConnectPG(connuri)
	if err != nil {
		log.Fatalf("Error connecting to DB. Reason: %s", err)
	}
	err = datastore.CreateSchema(db)
	if err != nil {
		log.Fatalf("Error creating Postgres schema. Reason: %s", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8081", r)
}
