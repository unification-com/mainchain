# Simple usage with a mounted data directory:
# > docker build -t undd .
# > docker run -it -p 26657:26657 -p 26656:26656 -v $HOME/.und_mainchain:/root/.und_mainchain undd und init [node_name]
# > docker run -it -p 26657:26657 -p 26656:26656 -v $HOME/.und_mainchain:/root/.und_mainchain undd und start
FROM golang:alpine3.12 AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev py-pip curl --upgrade grep

RUN apk add --update --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/unification-com

# Add source files
RUN git clone -b $(curl --silent "https://api.github.com/repos/unification-com/mainchain/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")') https://github.com/unification-com/mainchain

WORKDIR /go/src/github.com/unification-com/mainchain

# Install minimum necessary dependencies, build und & undcli, remove packages
RUN make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/und /usr/bin/und

# Run und by default, omit entrypoint to ease using container with undcli
CMD ["und"]
