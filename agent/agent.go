package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func HomeHandler (w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("xplex rig agent\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8082", r)
}
