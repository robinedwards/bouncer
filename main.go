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
	"sync"
	"syscall"
)

func setupRouter(cfg func() *config.Config) http.Handler {
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

	go func() {
		sighup := make(chan os.Signal, 1)
		signal.Notify(sighup, syscall.SIGHUP)

		for {
			<-sighup
			logger.Rotate()
		}
	}()
}

func setupConfigReload(filename string) {
	go func() {
		reload := make(chan os.Signal, 1)
		signal.Notify(reload, syscall.SIGUSR2)

		for {
			<-reload
			tmp, err := config.LoadConfigFile(filename)

			if err != nil {
				fmt.Errorf("Error reloading config %s", err)
				return
			}

			configLock.RLock()
			*bouncerConfig = *tmp
			configLock.RUnlock()

			fmt.Println("Reloaded config:", filename)
		}
	}()
}

var bouncerConfig *config.Config
var configLock = new(sync.RWMutex)

func getConfig() *config.Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return bouncerConfig
}

func main() {
	// Parse arguments
	listenPtr := flag.String("listen", "localhost:8000", "host and port to listen on")
	configPtr := flag.String("config", "config.json", "config file")
	sentryPtr := flag.String("sentry", "", "Sentry DSN")
	logFilePtr := flag.String("log", "./participation.log", "log file")
	flag.Parse()

	if len(*sentryPtr) > 0 {
		fmt.Println("Setting up Sentry")
		raven.SetDSN(*sentryPtr)
	}

	if len(*configPtr) == 0 {
		log.Fatalf("No config file supplied with --config")
	}

	// load config
	tmp, err := config.LoadConfigFile(*configPtr)
	if err != nil {
		log.Fatalf("Couldn't load config '%s': %s", *configPtr, err)
	}

	bouncerConfig = tmp

	// setup signal handlers
	setupConfigReload(*configPtr)
	setupMetricsLogging(*logFilePtr)

	fmt.Println("Listening on", *listenPtr, "and logging to", *logFilePtr, "pid", os.Getpid())
	err = http.ListenAndServe(*listenPtr, setupRouter(getConfig))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
