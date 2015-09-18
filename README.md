# Bouncer

## Getting started

    $ go get github.com/gorilla/mux github.com/unrolled/render github.com/serialx/hashring

    $ export GOPATH=$(pwd)

    $ go test bouncer/experiment_test
    $ go test bouncer/handler_test

    $ go run main.go
    Listening on localhost:8000
