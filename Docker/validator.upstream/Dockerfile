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

RUN git clone https://github.com/unification-com/mainchain.git
WORKDIR $GOPATH/src/github.com/unification-com/mainchain

RUN make install

COPY Docker/assets /root/.und_mainchain

RUN echo "/root/.go/bin/und" >> /root/.bash_history && \
    echo "alias ll='ls -la'" >> /root/.bashrc
