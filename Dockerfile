ARG IMG_TAG=latest

# Compile the gaiad binary
FROM golang:1.22-alpine AS und-builder
WORKDIR /src/app/
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN apk add --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN CGO_ENABLED=0 make install

# Add to a distroless container
FROM cgr.dev/chainguard/static:$IMG_TAG
ARG IMG_TAG
COPY --from=und-builder /go/bin/und /usr/local/bin/
EXPOSE 26656 26657 1317 9090
USER 0

ENTRYPOINT ["und", "start"]
