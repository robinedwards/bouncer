package main

import (
	"bouncer/config"
	"bouncer/experiment"
	"bouncer/handlers"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type BouncerDB struct {
	ActiveExperiments []experiment.Experiment
}

func (db BouncerDB) GetExperiments() []experiment.Experiment {
	return db.ActiveExperiments
}

var db BouncerDB

func init() {
	db.ActiveExperiments = make([]experiment.Experiment, 0)

	// setup experiments for demo
	db.ActiveExperiments = append(db.ActiveExperiments, experiment.NewExperiment("test1",
		experiment.Alternative{Name: "a", Weight: 1},
		experiment.Alternative{Name: "b", Weight: 1}))

	db.ActiveExperiments = append(db.ActiveExperiments, experiment.NewExperiment("test2",
		experiment.Alternative{Name: "a", Weight: 1},
		experiment.Alternative{Name: "b", Weight: 1}))

}

func main() {
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	configPtr := flag.String("config", "", "config file")

	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ListExperiments(db))
	router.HandleFunc("/participate/", handlers.Participate(db))

	http.Handle("/", router)

	if len(*configPtr) == 0 {
		fmt.Println("No config file supplied with --config")
		return
	}

	bouncerConfig, err := config.LoadConfigFile(*configPtr)
	if err != nil {
		panic(err)
	}

	fmt.Println(bouncerConfig)

	fmt.Println("Listening on " + *listenPtr)
	http.ListenAndServe(*listenPtr, nil)
}
