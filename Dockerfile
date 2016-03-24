FROM golang:1.6-alpine

COPY . /go/src/bouncer
ENV GOPATH /go/src/bouncer

WORKDIR $GOPATH

RUN apk add --no-cache git && \
    go get -d . && \
    go test bouncer/...  && \
    echo $GOPATH && \
    go build -o bouncer main.go && \
    chmod +x bouncer && \
    apk del --no-cache --purge git && \
    cp bouncer /usr/bin && \
    rm -rf ../bouncer

EXPOSE 9000

ENTRYPOINT ["/usr/bin/bouncer"]
CMD ["-config", "/etc/bouncer/config.json", "-listen", "0.0.0.0:9000", "-fluent", "localhost:24220"]
