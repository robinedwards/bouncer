package main

import (
	"bouncer/config"
	"bouncer/handlers"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// parse arguments
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	configPtr := flag.String("config", "", "config file")

	flag.Parse()

	// load config
	bouncerConfig, err := config.LoadConfigFile(*configPtr)
	if err != nil {
		panic(err)
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


	http.Handle("/", router)

	fmt.Println(bouncerConfig)

	fmt.Println("Listening on " + *listenPtr)
	http.ListenAndServe(*listenPtr, nil)
}
