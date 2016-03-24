FROM golang:1.6-alpine

COPY . /go/src/bouncer
ENV GOPATH /go/src/bouncer

WORKDIR $GOPATH

RUN go get bouncer/... && \
    go test bouncer/...  && \
    go build -o bouncer main.go && \
    chmod +x bouncer

EXPOSE 9000

ENTRYPOINT ["/go/src/bouncer/bouncer"]
CMD ["-config", "/etc/bouncer/config.json", "-listen", "0.0.0.0:9000", "-fluent", "localhost:24220"]
