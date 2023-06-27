FROM golang:1.20

MAINTAINER Pouyan Heyratpour <pouyan@janstun.com>

ARG ACCESS_TOKEN
RUN git config --global url."https://${ACCESS_TOKEN}@github.com".insteadOf "https://github.com"

RUN apt-get update && \
    apt-get install -y protobuf-compiler golang-golang-x-tools rsync

RUN go install -v github.com/pouyanh/polywatch/cmd/polywatch@latest && \
    go install -v github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT ["/go/bin/polywatch"]
