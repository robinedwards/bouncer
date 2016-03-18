# Bouncer

## Getting started

Install dependencies

    $ export GOPATH=$(pwd)
    $ go get bouncer/...

Run all tests:

    $ go test bouncer/...

Execute:

    $ go run main.go -config yourconfig.json
    Listening on localhost:8000

Flags:

    -help
    -config string
        config file (default "config.json")
    -listen string
        host and port to listen on (default "localhost:8000")
    -log string
        log file (default "./participation.log")
    -sentry string
        Sentry DSN

Signals:
    SIGHUP - log rotate
    SIGUSR2 - reload configuration file
