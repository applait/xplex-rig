package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// HomeHandler handles home URL route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8081", r)
}
