![Unification](https://raw.githubusercontent.com/unification-com/mainchain/master/unification_logoblack.png "Unification")

## UND Mainchain

Official Golang implementation of Unification Mainchain.

See [Documentation](docs/README.md) for full guides.

[![Go Report Card](https://goreportcard.com/badge/github.com/unification-com/mainchain)](https://goreportcard.com/report/github.com/unification-com/mainchain)
[![Join the chat at https://gitter.im/unification-com/mainchain](https://badges.gitter.im/unification-com/mainchain.svg)](https://gitter.im/unification-com/mainchain?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## HD Wallet Path

BIP-0044 Path for our HD Wallets is as follows:

`44'/5555'/0'/0`   

SLIP-0044 Coin ID is `5555`

## Pre-compiled binaries

The quickest way to obtain and run the `und` and `undcli` applications is to download
the pre-compiled binaries from [latest release](https://github.com/unification-com/mainchain/releases)

## Build

Compile `und` and `undcli` binaries and output to ./build

```bash
make build
```

## Install

Install `und` and `undcli` binaries into `$GOPATH`

```bash
make install
```

### Dockerised `und` and `undcli`

The Dockerised binaries can be used instead of installing locally. The Docker container
will use the latest release tag to build the binaries.

Build the container:

```bash
docker build -t undd .
```

Example commands, with mounted data directories:

```bash
$ docker run -it -p 26657:26657 -p 26656:26656 -v ~/.und_mainchain:/root/.und_mainchain -v ~/.und_cli:/root/.und_cli undd und init [node_name]
$ docker run -it -p 26657:26657 -p 26656:26656 -v ~/.und_mainchain:/root/.und_mainchain -v ~/.und_cli:/root/.und_cli undd und start
```

## DevNet Enviroment

A complete DevNet environment, comprising of 3 EVs, a REST server, a reverse proxy server
and several test wallets loaded with UND is available via Docker Compose compositions 
for development and testing purposes.

### Local build

The local build copies the current local codebase to the Docker containers, and is used during
development to test changes before committing to the repository.

```
docker-compose -f Docker/docker-compose.local.yml up --build
docker-compose -f Docker/docker-compose.local.yml down --remove-orphans
```

or using the `make` target:

```bash
make devnet
make devnet-down
```

### Pure Upstream build

Pure upstream uses the `master` branch of this repository to build the binaries, and is useful
for testing the latest master version.

```
docker-compose -f Docker/docker-compose.upstream.yml up --build
docker-compose -f Docker/docker-compose.upstream.yml down --remove-orphans
```

or using the `make` target:

```bash
make devnet-pristine
make devnet-pristine-down
```

#### REST API Endpoints

With DevNet up, the REST API endpoints can be seen via http://localhost:1318/swagger-ui/

### Importing docker composition keys

The DevNet node keys can be imported locally for ease of testing:

```
undcli keys import node1 Docker/assets/keys/node1.key
undcli keys import node2 Docker/assets/keys/node2.key
undcli keys import node3 Docker/assets/keys/node3.key
```

A full list of DevNet accounts and keys for importing can be found in 
[Docker/README.md](https://github.com/unification-com/mainchain/blob/master/Docker/README.md)

### Useful Defaults

`undcli` defaults for DevNet can be set as follows. This will set the corresponding values in
`$HOME/.und_cli/config/config.toml`

```
undcli config chain-id UND-Mainchain-DevNet
undcli config node tcp://localhost:26661
```

## `undcli` command overview

### Query account

```
undcli query account und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
```

### Interacting with WRKChain module

Register:
undcli tx wrkchain register --moniker="wrkchain1" --genesis="genesishashkjwnedjknwed" --name="Wrkchain 1" --base="geth" --from wrktest --gas=auto --gas-adjustment=1.15

then run:
```
undcli query tx [TX HASH]
```

this will return the generated WRKChain ID integer

Query metadata
```
undcli query wrkchain get 1
```

Filter WRKChain metadata

```
undcli query wrkchain search --moniker wrkchain1
undcli query wrkchain search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
undcli query wrkchain search --page=2 --limit=100
```

Record block hash(es)
```
undcli tx wrkchain record 1 --wc_height=1 --block_hash="d04b98f48e8" --parent_hash="f8bcc15c6ae" --hash1="5ac050801cd6" --hash2="dcfd428fb5f9e" --hash3="65c4e16e7807340fa" --from wrktest --gas=auto --gas-adjustment=1.15
```

Query a block
```
undcli query wrkchain block 1 1
```

#### WRKChain REST

http://localhost:1318/wrkchain/wrkchains  
http://localhost:1318/wrkchain/wrkchains?moniker=wrkchain1  
http://localhost:1318/wrkchain/wrkchains?owner=[bech32address]  
http://localhost:1318/wrkchain/1  
http://localhost:1318/wrkchain/1/block/1  

### Interacting with BEACON module

Register:
undcli tx beacon register --moniker=beacon1 --name="Beacon 1" --from wrktest

then run:
```
undcli query tx [TX HASH]
```

this will return the generated Beacon ID integer

Query metadata
```
undcli query beacon beacon 1
```

Record Timestamp hash
```
undcli tx beacon record 1 --hash=d04b98f48e8 --subtime=$(date +%s) --from wrktest --gas=auto --gas-adjustment=1.15
```

Query a Timestamp
```
undcli query beacon timestamp 1 1
```

#### BEACON REST

http://localhost:1318/beacon/beacons  
http://localhost:1318/beacon/beacons?moniker=beacon1  
http://localhost:1318/beacon/beacons?owner=[bech32address]  
http://localhost:1318/beacon/1   
http://localhost:1318/beacon/1/timestamp/1  

### Purchase Enterprise UND

Raise purchase order:
```
undcli tx enterprise purchase 1002000000000nund --from wrktest --gas=auto --gas-adjustment=1.15
```

List purchase orders:
```
undcli query enterprise get-all-pos
```

get specific purchase order:
```
undcli query enterprise get 1
```

Accept purchase order (must be sent from specified enterprise account)
```
undcli tx enterprise process 1 accept --from ent --gas=auto --gas-adjustment=1.15
```

Reject purchase order:
```
undcli tx enterprise process 1 reject --from ent --gas=auto --gas-adjustment=1.15
```

Query total locked enterprise UND
```
undcli query enterprise total-locked
```

Query locked enterprise UND for an account
```
undcli query enterprise locked [address]
```

#### Enterprise REST

http://localhost:1317/enterprise/params  
http://localhost:1317/enterprise/locked  
http://localhost:1317/enterprise/unlocked  
http://localhost:1317/enterprise/pos  
http://localhost:1317/enterprise/po/1  
http://localhost:1317/enterprise/[bech32addr]/locked

## Invariance checking

Start a full node with the `--inv-check-period` flag. Value of 1 will
check every block for invariances:

```
und start --inv-check-period 1
```

Invariance Tx can be sent using:

```
undcli tx crisis invariant-broken enterprise module-account --from wrktest
```

## Simulation

### test full app simulation

```
go test -mod=readonly ./simapp \
    -run=TestFullAppSimulation \
    -Enabled=true \
    -NumBlocks=500 \
    -BlockSize=300 \
    -Commit=true \
    -Seed=24 \
    -Period=1 \
    -PrintAllInvariants=true \
    -ExportParamsPath=/path/to/.simapp/params.json \
    -ExportStatePath=/path/to/.simapp/state.json \
    -ExportStatsPath=/path/to/.simapp/statistics.json \
    -Verbose=true \
    -v \
    -timeout 24h
```

### benchmark test

```
go test -mod=readonly -benchmem -run=^$ github.com/unification-com/mainchain/simapp -bench ^BenchmarkFullAppSimulation -Commit=true -cpuprofile /path/to/.simapp/cpu.out -v -timeout 24h
```

#### pprof tools

```
go tool pprof /path/to/.simapp/cpu.out
(pprof) top
(pprof) list [function]
(pprof) web
(pprof) quit
```
