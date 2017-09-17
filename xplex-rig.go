package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/applait/xplex-rig/config"
	"github.com/applait/xplex-rig/models"
	"github.com/applait/xplex-rig/rest"
)

// HomeHandler handles home URL route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}

func main() {
	conf, err := config.ParseConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config. Reason: %s", err)
	}
	db, err := models.ConnectPG(conf.Server.PostgresURL)
	if err != nil {
		log.Fatalf("Error connecting to database. Reason: %s", err)
	}
	if err = models.CreateSchema(db); err != nil {
		log.Fatalf("Error creating schema. Reason: %s\n", err)
		return
	}
	log.Printf("Starting HTTP server on port %d", conf.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Server.Port), rest.Start(db, &conf)))
}
