package main

import (
	"bouncer/config"
	"bouncer/handlers"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

import gorillahandlers "github.com/gorilla/handlers"

func main() {
	// parse arguments
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	configPtr := flag.String("config", "config.json", "config file")

	flag.Parse()

	// load config
	bouncerConfig, err := config.LoadConfigFile(*configPtr)
	if err != nil {
		os.Exit(1)
	}

	if len(*configPtr) == 0 {
		fmt.Println("No config file supplied with --config")
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Root)
	router.HandleFunc("/experiments/", handlers.ListExperiments(bouncerConfig))
	router.HandleFunc("/groups/", handlers.ListGroups(bouncerConfig))
	router.HandleFunc("/features/", handlers.ListFeatures(bouncerConfig))
	router.HandleFunc("/participate/", handlers.Participate(bouncerConfig))

	loggedRouter := gorillahandlers.LoggingHandler(os.Stdout, router)
	http.Handle("/", loggedRouter)

	fmt.Println(bouncerConfig)

	fmt.Println("Listening on " + *listenPtr)
	http.ListenAndServe(*listenPtr, nil)
}
