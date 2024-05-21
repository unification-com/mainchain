#!/usr/bin/make -f

PACKAGES=$(shell go list ./... )

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BINDIR ?= $(GOPATH)/bin
DOCKER := $(shell which docker)
TM_CORE_SEM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::') # grab everything after the space in "github.com/cometbft/cometbft v0.34.7"

LATEST_RELEASE := $(shell curl --silent https://api.github.com/repos/unification-com/mainchain/releases/latest | grep -Po '"tag_name": \"\K.*?(?=\")')

#export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=UndMainchain \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=und \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_CORE_SEM_VERSION) \
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
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
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
	TM_CORE_SEM_VERSION="${TM_CORE_SEM_VERSION}" goreleaser --snapshot --skip-publish --clean --debug

release: goreleaser
	TM_CORE_SEM_VERSION="${TM_CORE_SEM_VERSION}" goreleaser --clean

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
#protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user $(shell id -u):$(shell id -g) $(protoImageName)

proto-all: proto-format proto-lint proto-gen

# NOTE: when using rootless docker, this will fail. Before running, run:
#   chmod 777 proto/buf.lock
#   mkdir github.com && chmod 777 github.com
# After running, run:
#   sudo chown -R $(id -u):$(id -g) github.com
#   cp -r github.com/unification-com/mainchain/* ./
#   rm -rf github.com
proto-gen:
	@echo "Generating Protobuf files"
	@chmod 777 proto/buf.lock
	@mkdir github.com && chmod 777 github.com
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@mkdir tmp-swagger-gen && chmod 777 tmp-swagger-gen
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh
	#$(MAKE) update-swagger-docs

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

proto-update-deps:
	@echo "Updating Protobuf dependencies"
	$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps