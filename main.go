package main

import gorillahandlers "github.com/gorilla/handlers"

import (
	"bouncer/config"
	"bouncer/handlers"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"log"
	"os/signal"
	"syscall"
)

func setupRouter(cfg config.Config) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Root)
	router.HandleFunc("/experiments/", handlers.ListExperiments(cfg))
	router.HandleFunc("/groups/", handlers.ListGroups(cfg))
	router.HandleFunc("/features/", handlers.ListFeatures(cfg))
	router.HandleFunc("/participate/", handlers.Participate(cfg))

	return gorillahandlers.LoggingHandler(os.Stdout, router)
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
	logFilePtr := flag.String("log", "./paricipation.log", "log file")

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
	setupMetricsLogging(*logFilePtr)
	http.Handle("/", setupRouter(bouncerConfig))

	fmt.Println("Listening on", *listenPtr, "and logging to", *logFilePtr)
	http.ListenAndServe(*listenPtr, nil)
}
