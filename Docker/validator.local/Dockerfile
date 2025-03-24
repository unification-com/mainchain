FROM golang:1.22-alpine

RUN apk update && \
    apk upgrade && \
    apk add git make jq gcc libc-dev curl nano --no-cache --upgrade grep

ENV GOPATH="/root/.go"
ENV GOROOT="/usr/local/go"
ENV PATH="/usr/local/go/bin:/root/.go/bin:${PATH}"
ENV LEDGER_ENABLED="false"

WORKDIR /root
RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN mkdir -p /root/src/mainchain && \
    export GO111MODULE=off

WORKDIR /root/src/mainchain

COPY ./go.mod ./go.sum ./
RUN cd /root/src/mainchain && go mod tidy && go mod download

COPY . ./

ARG WITH_DELVE="no"
ENV WITH_DELVE=$WITH_DELVE
RUN export WITH_DELVE=$WITH_DELVE && make install && /root/.go/bin/und version

COPY Docker/assets /root/.und_mainchain
