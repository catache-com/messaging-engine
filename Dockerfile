FROM golang:1.19 as builder

WORKDIR /go/src/messaging-engine/
COPY . .
RUN go get
RUN go build -v -ldflags "-w -s -X main.Version=1.0.0 -X main.Build=`date +%FT%T%z`" -o bin/messaging-engine-linux-amd64

FROM debian:buster-slim

MAINTAINER catache.com

RUN apt update
RUN apt-get -y install ca-certificates

WORKDIR /application

COPY --from=builder /go/src/messaging-engine/bin/messaging-engine-linux-amd64 .

ENTRYPOINT ./messaging-engine-linux-amd64
