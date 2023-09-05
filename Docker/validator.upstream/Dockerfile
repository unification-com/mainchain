FROM golang:1.19-alpine3.15

ARG GO_VERSION

RUN apk update && \
    apk upgrade && \
    apk add git make jq curl --no-cache --upgrade grep

ENV GOPATH="/root/.go"
ENV GOROOT="/usr/local/go"
ENV GO111MODULE="on"
ENV PATH="/usr/local/go/bin:/root/.go/bin:${PATH}"
ENV LEDGER_ENABLED="false"

RUN mkdir -p $GOPATH/src/github.com/unification-com
WORKDIR $GOPATH/src/github.com/unification-com

COPY .vers_docker /root

RUN rm -rf $GOPATH/src/github.com/unification-com/mainchain && rm -f /root/.go/bin/und && rm -f /root/.go/bin/undcli
RUN git clone -b $(cat /root/.vers_docker | tr -d '\n') https://github.com/unification-com/mainchain.git
WORKDIR $GOPATH/src/github.com/unification-com/mainchain

RUN make install

COPY Docker/assets /root/.und_mainchain

RUN echo "/root/.go/bin/und" >> /root/.bash_history && \
    echo "alias ll='ls -la'" >> /root/.bashrc
