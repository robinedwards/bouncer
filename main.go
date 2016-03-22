package main

import gorillahandlers "github.com/gorilla/handlers"

import (
	"bouncer/config"
	"bouncer/handlers"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

func setupRouter(cfg func() *config.Config, logger func(interface{})) http.Handler {
	ourStats := stats.New()

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Root)
	router.HandleFunc("/experiments/", handlers.ListExperiments(cfg))
	router.HandleFunc("/groups/", handlers.ListGroups(cfg))
	router.HandleFunc("/features/", handlers.ListFeatures(cfg))
	router.HandleFunc("/participate/", handlers.Participate(cfg, logger))
	router.HandleFunc("/error/", func(w http.ResponseWriter, r *http.Request) { panic("error") })
	router.HandleFunc("/stats/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(ourStats.Data())
		w.Write(b)
	})

	return gorillahandlers.LoggingHandler(os.Stdout, handlers.ErrorHandler(
		ourStats.Handler(router)))
}

func setupFluentd(fluentHost string) func(interface{}) {
	parts := strings.Split(fluentHost, ":")
	if len(parts) != 2 {
		fmt.Fprintf(os.Stderr, "invalid fluentd host format, should be <hostname>:<port>")
		os.Exit(1)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid port for fluentd")
		os.Exit(1)
	}

	cfg := fluent.Config{FluentHost: parts[0], FluentPort: port}
	logger, err := fluent.New(cfg)

	return func(message interface{}) {
		body, err := json.Marshal(message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: couldn't encode json response: %v", err)
		}

		fmt.Println(body)
		logger.Post("bouncer", message)
	}
}

func setupConfigReload(filename string) {
	go func() {
		reload := make(chan os.Signal, 1)
		signal.Notify(reload, syscall.SIGUSR2)

		for {
			<-reload
			tmp, err := config.LoadConfigFile(filename)

			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reloading config", err)
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
	fluentPtr := flag.String("fluent", "localhost:24220", "td-agent host and port")
	flag.Parse()

	if len(*sentryPtr) > 0 {
		fmt.Println("Setting up Sentry")
		raven.SetDSN(*sentryPtr)
	}

	// load config
	tmp, err := config.LoadConfigFile(*configPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't load config '%s': %s\n", *configPtr, err)
		os.Exit(1)
	}

	bouncerConfig = tmp

	setupConfigReload(*configPtr)
	logger := setupFluentd(*fluentPtr)

	fmt.Println("Listening on", *listenPtr, "and logging to", *fluentPtr, "pid", os.Getpid())
	err = http.ListenAndServe(*listenPtr, setupRouter(getConfig, logger))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(74) // io error
	}
}
