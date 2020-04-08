FROM golang:1.13.3-alpine

RUN apk update && \
    apk upgrade && \
    apk add git make jq gcc libc-dev

ENV GOPATH="/root/.go"
ENV GOROOT="/usr/local/go"
ENV GO111MODULE="on"
ENV PATH="/usr/local/go/bin:/root/.go/bin:${PATH}"
ENV LEDGER_ENABLED="false"

RUN mkdir -p $GOPATH/src/github.com/unification-com && \
    rm -rf $GOPATH/src/github.com/unification-com/mainchain
RUN RUN export GO111MODULE=off && go get -u github.com/go-delve/delve/cmd/dlv

COPY . $GOPATH/src/github.com/unification-com/mainchain

WORKDIR $GOPATH/src/github.com/unification-com/mainchain

ARG WITH_DELVE="no"
ENV WITH_DELVE=$WITH_DELVE
RUN export WITH_DELVE=$WITH_DELVE && make install

COPY Docker/assets /root/.und_mainchain

RUN echo "/root/.go/bin/und" >> /root/.bash_history && \
    echo "alias ll='ls -la'" >> /root/.bashrc
