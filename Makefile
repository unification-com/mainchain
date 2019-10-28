PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=UndMainchain \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=und \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=undcli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags)"

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)' -gcflags 'all=-N -l'

include Makefile.ledger
all: lint install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/und
		go install $(BUILD_FLAGS) ./cmd/undcli

build: go.sum
		go build -mod=readonly $(BUILD_FLAGS) -o build/und ./cmd/und
		go build -mod=readonly $(BUILD_FLAGS) -o build/undcli ./cmd/undcli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

clean:
	rm -rf build/

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
