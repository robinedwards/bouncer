package main

import gorillahandlers "github.com/gorilla/handlers"

import (
	"bouncer/config"
	"bouncer/handlers"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func setupRouter(cfg config.Config) http.Handler {
	ourStats := stats.New()

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Root)
	router.HandleFunc("/experiments/", handlers.ListExperiments(cfg))
	router.HandleFunc("/groups/", handlers.ListGroups(cfg))
	router.HandleFunc("/features/", handlers.ListFeatures(cfg))
	router.HandleFunc("/participate/", handlers.Participate(cfg))
	router.HandleFunc("/error/", func(w http.ResponseWriter, r *http.Request) { panic("error") })
	router.HandleFunc("/stats/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(ourStats.Data())
		w.Write(b)
	})

	return gorillahandlers.LoggingHandler(os.Stdout, handlers.ErrorHandler(
		ourStats.Handler(router)))
}

func setupMetricsLogging(logfile string) {
	logger := &lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	log.SetOutput(logger)
	log.SetFlags(0)
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	go func() {
		for {
			<-sighup
			logger.Rotate()
		}
	}()
}

func main() {
	// parse arguments
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	configPtr := flag.String("config", "config.json", "config file")
	sentryPtr := flag.String("sentry", "", "Sentry DSN")
	logFilePtr := flag.String("log", "./paricipation.log", "log file")
	flag.Parse()

	if len(*sentryPtr) > 0 {
		fmt.Println("Setting up Sentry")
		raven.SetDSN(*sentryPtr)
	}

	if len(*configPtr) == 0 {
		log.Fatalf("No config file supplied with --config")
	}

	// load config
	bouncerConfig, err := config.LoadConfigFile(*configPtr)

	if err != nil {
		log.Fatalf("Couldn't load config '%s': %s", *configPtr, err)
	}

	setupMetricsLogging(*logFilePtr)

	fmt.Println("Listening on", *listenPtr, "and logging to", *logFilePtr)

	err = http.ListenAndServe(*listenPtr, setupRouter(bouncerConfig))
	if err != nil {
		fmt.Errorf("Error listening: %s", err)
		os.Exit(1)
	}
}
