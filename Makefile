#!/usr/bin/make -f

PACKAGES=$(shell go list ./... )

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BINDIR ?= $(GOPATH)/bin
DOCKER := $(shell which docker)
#TM_CORE_SEM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::') # grab everything after the space in "github.com/cometbft/cometbft v0.34.7"
COSMOS_SDK_SEM_VERSION := $(shell go list -m github.com/cosmos/cosmos-sdk | sed 's:.* ::') # used in Swagger dep download
IBC_GO_SEM_VERSION := $(shell grep 'github.com/cosmos/ibc-go/v8' go.mod | sed 's:.* ::') # used in Swagger dep download
LATEST_RELEASE := $(shell curl --silent https://api.github.com/repos/unification-com/mainchain/releases/latest | grep -Po '"tag_name": \"\K.*?(?=\")')

#export GO111MODULE = on

build_tags = netgo

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=UndMainchain \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=und \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

export UND_LDFLAGS = $(ldflags)

include scripts/makefiles/devtools.mk
include scripts/makefiles/ledger.mk
include scripts/makefiles/proto.mk
include scripts/makefiles/unittests.mk
include scripts/makefiles/sims.mk
include scripts/makefiles/devnet.mk

BUILD_FLAGS := -tags="$(build_tags)" -ldflags '$(ldflags)'

ifeq ($(WITH_DELVE),yes)
  BUILD_FLAGS += -gcflags 'all=-N -l'
endif

all: lint install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/und

build: clean go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/und ./cmd/und

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

lint:
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	go mod verify

gofmt:
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	go mod verify


clean:
	rm -rf build/

# Used during active development

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
PHONY: go-mod-cache

check-updates:
	@echo "checking for module updates"
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null
	@echo "run:"
	@echo "go get github.com/user/repo to update. E.g. go get github.com/cosmos/cosmos-sdk"

snapshot: goreleaser
	goreleaser --snapshot --skip=publish --clean --debug

release: goreleaser
	goreleaser --clean

