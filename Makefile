PACKAGES=$(shell go list ./... )

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BINDIR ?= $(GOPATH)/bin

export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=UndMainchain \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=und \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=undcli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

export UND_LDFLAGS = $(ldflags)

include Makefile.devtools
include Makefile.ledger

BUILD_FLAGS := -tags="$(build_tags)" -ldflags '$(ldflags)'

ifeq ($(WITH_DELVE),yes)
  BUILD_FLAGS += -gcflags 'all=-N -l'
endif

all: lint install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/und
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/undcli

build: clean go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/und ./cmd/und
	go build -mod=readonly $(BUILD_FLAGS) -o build/undcli ./cmd/undcli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

lint:
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	go mod verify

test: test-unit

test-unit:
	@go test -mod=readonly ./...

test-race:
	@go test -mod=readonly -race ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic ./...

test-cover-html: test-cover
	@go tool cover -html=coverage.txt

test-no-cache:
	@go clean -testcache
	@go test -v -mod=readonly ./...

clean:
	rm -rf build/

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/lcd/swagger-ui -dest=client/lcd -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
    	echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
.PHONY: update-swagger-docs

# Docker compositions

devnet:
	docker-compose -f Docker/docker-compose.local.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.local.yml up --build

devnet-down:
	docker-compose -f Docker/docker-compose.local.yml down --remove-orphans

devnet-pristine:
	docker-compose -f Docker/docker-compose.upstream.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.upstream.yml up --build

devnet-pristine-down:
	docker-compose -f Docker/docker-compose.upstream.yml down

devnet-systemtest:
	docker-compose -f Docker/docker-compose.systemtest.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.systemtest.yml up --build

devnet-systemtest-down:
	docker-compose -f Docker/docker-compose.systemtest.yml down

# Used during active development

include Makefile.sims

deps:
	go get -u ./...

update-sdk:
	go get github.com/cosmos/cosmos-sdk@master

build-update-sdk:
	go build $(BUILD_FLAGS) -o build/und ./cmd/und
	go build $(BUILD_FLAGS) -o build/undcli ./cmd/undcli

snapshot: goreleaser
	UND_BUILD_TAGS="$(build_tags)" goreleaser --snapshot --skip-publish --rm-dist

release: goreleaser
	UND_BUILD_TAGS="$(build_tags)" goreleaser --rm-dist
