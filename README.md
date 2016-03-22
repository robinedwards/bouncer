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
    -fluent string
        fluentd host and port
    -sentry string
        Sentry DSN

Signals:
    SIGUSR2 - reload configuration file
