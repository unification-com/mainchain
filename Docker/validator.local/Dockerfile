FROM golang:1.13.1-alpine

ARG GO_VERSION

RUN apk update && \
    apk upgrade && \
    apk add git make jq

ENV GOPATH="/root/.go"
ENV GOROOT="/usr/local/go"
ENV GO111MODULE="on"
ENV PATH="/usr/local/go/bin:/root/.go/bin:${PATH}"
ENV LEDGER_ENABLED="false"

RUN mkdir -p $GOPATH/src/github.com/unification-com
WORKDIR $GOPATH/src/github.com/unification-com

RUN git clone https://github.com/unification-com/mainchain-cosmos.git

COPY . $GOPATH/src/github.com/unification-com/mainchain-cosmos

WORKDIR $GOPATH/src/github.com/unification-com/mainchain-cosmos
RUN make install

COPY Docker/assets /root/.und_mainchain
RUN /root/.go/bin/und validate-genesis /root/.und_mainchain/node0/config/genesis.json

RUN echo "/root/.go/bin/und" >> /root/.bash_history && \
    echo "alias ll='ls -la'" >> /root/.bashrc
