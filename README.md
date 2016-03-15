# Bouncer

## Getting started

    $ go get bouncer/...
    $ export GOPATH=$(pwd)

Run all tests:

    $ go test bouncer/...

Execute:

    $ go run main.go
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
