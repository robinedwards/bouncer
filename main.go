package main

import (
	"bouncer/config"
	"bouncer/experiment"
	"bouncer/handlers"
	"flag"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
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

	bouncerConfig := config.LoadConfig(`{
  "groups": [
    {
      "name": "admins",
      "uids": [
        "8",
        "32"
      ]
    },
    {
      "name": "testers",
      "uids": [
        "8",
        "32"
      ]
    }
  ],
  "features": {
    "audio_mode": {
      "groups": {
        "admins": 1,
        "users": 0
      },
      "enabled": 0.5
    }
  },
  "experiments": [
    {
      "name": "progress_bar",
      "groups": [
        {
          "admins": "green",
          "testers": "red"
        }
      ],
      "alternatives": [
        {
          "name": "green",
          "weight": 1
        },
        {
          "name": "red",
          "weight": 9
        }
      ]
    }
  ]
}`)

	fmt.Println(bouncerConfig)
}

func main() {
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ListExperiments(db))
	router.HandleFunc("/participate/", handlers.Participate(db))

	http.Handle("/", router)

	fmt.Println("Listening on " + *listenPtr)
	http.ListenAndServe(*listenPtr, nil)
}
