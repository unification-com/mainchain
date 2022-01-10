#!/usr/bin/make -f

PACKAGES=$(shell go list ./... )

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BINDIR ?= $(GOPATH)/bin
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::') # grab everything after the space in "github.com/tendermint/tendermint v0.34.7"

LATEST_RELEASE := $(shell curl --silent https://api.github.com/repos/unification-com/mainchain/releases/latest | grep -Po '"tag_name": \"\K.*?(?=\")')

#export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=UndMainchain \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=und \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

export UND_LDFLAGS = $(ldflags)

include devtools.mk
include ledger.mk
include cleveldb.mk

BUILD_FLAGS := -tags="$(build_tags)" -ldflags '$(ldflags)'

ifeq ($(WITH_DELVE),yes)
  BUILD_FLAGS += -gcflags 'all=-N -l'
endif

all: lint install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/und

build: clean go.sum
	go build  $(BUILD_FLAGS) -o build/und ./cmd/und

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

include sims.mk

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

devnet-latest-release:
	@echo "${LATEST_RELEASE}" > ./.vers_docker
	docker-compose -f Docker/docker-compose.upstream.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.upstream.yml up --build

devnet-latest-release-down:
	docker-compose -f Docker/docker-compose.upstream.yml down

devnet-master:
	@echo "master" > ./.vers_docker
	docker-compose -f Docker/docker-compose.upstream.yml down --remove-orphans
	docker-compose -f Docker/docker-compose.upstream.yml up --build

devnet-master-down:
	docker-compose -f Docker/docker-compose.upstream.yml down

# Used during active development

check-updates:
	@echo "checking for module updates"
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null
	@echo "run:"
	@echo "go get github.com/user/repo to update. E.g. go get github.com/cosmos/cosmos-sdk"

snapshot: goreleaser
	UND_BUILD_TAGS="$(build_tags)" goreleaser --snapshot --skip-publish --rm-dist --debug

release: goreleaser
	UND_BUILD_TAGS="$(build_tags)" goreleaser --rm-dist

###############################################################################
###                                Protobuf                                 ###
###############################################################################

HTTPS_GIT = https://github.com/unification-com/mainchain

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace \
	--workdir /workspace tendermintdev/docker-build-proto \
	find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -i {} \;

proto-swagger-gen:
	@echo "Generating Swagger files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protoc-swagger-gen.sh

proto-check-breaking:
	@$(DOCKER_BUF) buf breaking --against-input $(HTTPS_GIT)#branch=stargate

TM_URL              = https://raw.githubusercontent.com/tendermint/tendermint/v0.34.0-rc6/proto/tendermint
GOGO_PROTO_URL      = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_PROTO_URL    = https://raw.githubusercontent.com/regen-network/cosmos-proto/master
CONFIO_URL          = https://raw.githubusercontent.com/confio/ics23/v0.6.3

TM_CRYPTO_TYPES     = third_party/proto/tendermint/crypto
TM_ABCI_TYPES       = third_party/proto/tendermint/abci
TM_TYPES            = third_party/proto/tendermint/types
TM_VERSION          = third_party/proto/tendermint/version
TM_LIBS             = third_party/proto/tendermint/libs/bits
TM_P2P              = third_party/proto/tendermint/p2p

GOGO_PROTO_TYPES    = third_party/proto/gogoproto
COSMOS_PROTO_TYPES  = third_party/proto/cosmos_proto
CONFIO_TYPES        = third_party/proto/confio

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)
	@curl -sSL $(COSMOS_PROTO_URL)/cosmos.proto > $(COSMOS_PROTO_TYPES)/cosmos.proto

## Importing of tendermint protobuf definitions currently requires the
## use of `sed` in order to build properly with cosmos-sdk's proto file layout
## (which is the standard Buf.build FILE_LAYOUT)
## Issue link: https://github.com/tendermint/tendermint/issues/5021
	@mkdir -p $(TM_ABCI_TYPES)
	@curl -sSL $(TM_URL)/abci/types.proto > $(TM_ABCI_TYPES)/types.proto

	@mkdir -p $(TM_VERSION)
	@curl -sSL $(TM_URL)/version/types.proto > $(TM_VERSION)/types.proto

	@mkdir -p $(TM_TYPES)
	@curl -sSL $(TM_URL)/types/types.proto > $(TM_TYPES)/types.proto
	@curl -sSL $(TM_URL)/types/evidence.proto > $(TM_TYPES)/evidence.proto
	@curl -sSL $(TM_URL)/types/params.proto > $(TM_TYPES)/params.proto
	@curl -sSL $(TM_URL)/types/validator.proto > $(TM_TYPES)/validator.proto
	@curl -sSL $(TM_URL)/types/block.proto > $(TM_TYPES)/block.proto

	@mkdir -p $(TM_CRYPTO_TYPES)
	@curl -sSL $(TM_URL)/crypto/proof.proto > $(TM_CRYPTO_TYPES)/proof.proto
	@curl -sSL $(TM_URL)/crypto/keys.proto > $(TM_CRYPTO_TYPES)/keys.proto

	@mkdir -p $(TM_LIBS)
	@curl -sSL $(TM_URL)/libs/bits/types.proto > $(TM_LIBS)/types.proto

	@mkdir -p $(TM_P2P)
	@curl -sSL $(TM_URL)/p2p/types.proto > $(TM_P2P)/types.proto

	@mkdir -p $(CONFIO_TYPES)
	@curl -sSL $(CONFIO_URL)/proofs.proto > $(CONFIO_TYPES)/proofs.proto
## insert go package option into proofs.proto file
## Issue link: https://github.com/confio/ics23/issues/32
	@sed -i '4ioption go_package = "github.com/confio/ics23/go";' $(CONFIO_TYPES)/proofs.proto

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps
