FROM golang:1.19-alpine3.15

RUN apk update && \
    apk upgrade && \
    apk add git make jq gcc libc-dev curl nano --no-cache --upgrade grep

ENV GOPATH="/root/.go"
ENV GOROOT="/usr/local/go"
ENV PATH="/usr/local/go/bin:/root/.go/bin:${PATH}"
ENV LEDGER_ENABLED="false"

RUN mkdir -p /root/src/mainchain && \
    export GO111MODULE=off && go get -u github.com/go-delve/delve/cmd/dlv

WORKDIR /root/src/mainchain

COPY ./go.mod ./go.sum ./
RUN cd /root/src/mainchain && go mod download && go mod tidy

COPY . ./

ARG WITH_DELVE="no"
ENV WITH_DELVE=$WITH_DELVE
RUN export WITH_DELVE=$WITH_DELVE && make install && /root/.go/bin/und version

COPY Docker/assets /root/.und_mainchain

RUN echo "/root/.go/bin/und" >> /root/.bash_history && \
    echo "alias ll='ls -la'" >> /root/.bashrc
