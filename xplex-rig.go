package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/applait/xplex-rig/config"
	"github.com/applait/xplex-rig/models"
	"github.com/applait/xplex-rig/rest"
)

func main() {
	confPath := flag.String("conf", "config.json", "Path to configuration JSON file relative to current directory.")
	createConfig := flag.Bool("createConfig", false, "Set this flag to generate a dummy config file and exit")

	flag.Parse()

	if *createConfig {
		_, err := config.CreateConfig(*confPath)
		if err != nil {
			log.Fatalf("Unable to generate config file. Reason: %s\n", err)
		}
		log.Printf("Created config file at %s\n", *confPath)
		os.Exit(0)
	}

	// TODO: parse config based on environment
	conf, err := config.ParseConfig(*confPath)
	if err != nil {
		log.Fatalf("Error loading config. Reason: %s", err)
	}

	// Run server
	db, err := models.ConnectPG(conf.PostgresURL)
	if err != nil {
		log.Fatalf("Error connecting to database. Reason: %s", err)
	}
	log.Printf("Starting HTTP server on port %d", conf.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), rest.Start(db, &conf)))
}
