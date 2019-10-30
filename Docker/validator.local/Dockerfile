FROM golang:1.13.1-alpine

ARG GO_VERSION

RUN apk update && \
    apk upgrade && \
    apk add git make jq gcc libc-dev

ENV GOPATH="/root/.go"
ENV GOROOT="/usr/local/go"
ENV GO111MODULE="on"
ENV PATH="/usr/local/go/bin:/root/.go/bin:${PATH}"
ENV LEDGER_ENABLED="false"

RUN mkdir -p $GOPATH/src/github.com/unification-com
WORKDIR $GOPATH/src/github.com/unification-com

RUN go get -u github.com/go-delve/delve/cmd/dlv

COPY . $GOPATH/src/github.com/unification-com/mainchain-cosmos

WORKDIR $GOPATH/src/github.com/unification-com/mainchain-cosmos

RUN more $GOPATH/src/github.com/unification-com/mainchain-cosmos/go.mod

ARG WITH_DELVE="no"
ENV WITH_DELVE=$WITH_DELVE
RUN export WITH_DELVE=$WITH_DELVE && make install

COPY Docker/assets /root/.und_mainchain
RUN /root/.go/bin/und validate-genesis /root/.und_mainchain/node1/config/genesis.json

RUN echo "/root/.go/bin/und" >> /root/.bash_history && \
    echo "alias ll='ls -la'" >> /root/.bashrc
