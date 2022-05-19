# Installing the Mainchain Software

This documentation outlines how to install the Unification Mainchain software, in
order to participate and interact with any of the Mainchain networks.

#### Contents

[[toc]]

## Installing the latest release binaries

The latest pre-compiled binaries are available from
[https://github.com/unification-com/mainchain/releases](https://github.com/unification-com/mainchain/releases).

- The `und` binary has been compiled for Linux, OSX and Windows.

Simply download the archives for your OS.

Once downloaded, you can verify the SHA256 checksum against those listed in the release's `checksums.txt`, for example:

```bash
$ openssl dgst -sha256 und_v1.5.0_linux_x86_64.tar.gz
SHA256(und_v1.5.0_linux_x86_64.tar.gz)= 98a93e757234f4cc408421b112bbc850975178900f3db53ab4a244f677041287
```

Extract them and copy the binaries to a suitable location - preferably a location in your `$PATH` environment variable,
for example `/usr/local/bin`, `/opt`, etc.

Once installed, verify:

```bash
$ und version --long
```

The output should match the latest release version tag.

## Building from Source

The Mainchain binaries can also be built from source.

### Prerequisites

`git`, `curl` and `make` are required to build the binaries. `jq` is also useful for quickly looking up 
values in `genesis.json`

These can all be installed via your package manager:

```bash
sudo apt-get install git curl make jq
```

or

```bash
sudo yum install git curl make jq
```

### Install Go

**Go 1.16+** is required to build the Mainchain binaries

Install `go` by following the [official docs](https://golang.org/doc/install).
Once Go is installed, set your `$PATH` environment variable:

```bash
$ mkdir -p $HOME/go/bin
$ echo "export PATH=$PATH:$(go env GOPATH)/bin" >> ~/.bash_profile
$ source ~/.bash_profile
```

### Build and install the binaries

::: warning IMPORTANT
unless you are contributing to Mainchain development, it is recommended you checkout and build from the latest release 
tag and **not** the `master` branch if you intend to connect to a live, public network (e.g. TestNet/MainNet).
:::

Download the **latest** tagged Mainchain release from
[https://github.com/unification-com/mainchain/releases](https://github.com/unification-com/mainchain/releases)

The **`[latest-release-tag]`** required for the command below can also be obtained by running:

```bash
curl --silent "https://api.github.com/repos/unification-com/mainchain/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")'
```

```bash
$ git clone -b [latest-release-tag] https://github.com/unification-com/mainchain
$ cd mainchain
$ make install
```

This will install the two binaries `und` and `und` into your `$HOME/go/bin`

### Verify the installation

Run the following commands:

```bash
$ und version --long
```

If they have installed correctly, you should see output similar to the following:

```bash
name: UndMainchain
server_name: und
version: 1.4.8-197-gf1df9cd
commit: f1df9cdf078acebb1bb1b26f0ac6c64d4496bae0
build_tags: netgo ledger
go: go version go1.16.2 linux/amd64
build_deps:
- github.com/99designs/keyring@v1.1.6 => github.com/cosmos/keyring@v1.1.7-0.20210622111912-ef00f8ac3d76
- github.com/ChainSafe/go-schnorrkel@v0.0.0-20200405005733-88cbf1b4c40d
- github.com/Workiva/go-datastructures@v1.0.52
- github.com/armon/go-metrics@v0.3.8
- github.com/beorn7/perks@v1.0.1
- github.com/bgentry/speakeasy@v0.1.0
- github.com/btcsuite/btcd@v0.22.0-beta
- github.com/cespare/xxhash/v2@v2.1.1
- github.com/confio/ics23/go@v0.6.6
- github.com/cosmos/btcutil@v1.0.4
- github.com/cosmos/cosmos-sdk@v0.42.11
- github.com/cosmos/go-bip39@v1.0.0
- github.com/cosmos/iavl@v0.17.3
- github.com/cosmos/ledger-cosmos-go@v0.11.1
- github.com/cosmos/ledger-go@v0.9.2
- github.com/davecgh/go-spew@v1.1.1
- github.com/dvsekhvalnov/jose2go@v0.0.0-20200901110807-248326c1351b
- github.com/felixge/httpsnoop@v1.0.1
- github.com/fsnotify/fsnotify@v1.4.9
- github.com/go-kit/kit@v0.10.0
- github.com/go-logfmt/logfmt@v0.5.0
- github.com/godbus/dbus@v0.0.0-20190726142602-4481cbc300e2
- github.com/gogo/gateway@v1.1.0
- github.com/gogo/protobuf@v1.3.3 => github.com/regen-network/protobuf@v1.3.3-alpha.regen.1
- github.com/golang/protobuf@v1.5.2
- github.com/golang/snappy@v0.0.3
- github.com/google/btree@v1.0.0
- github.com/google/orderedcode@v0.0.1
- github.com/gorilla/handlers@v1.5.1
- github.com/gorilla/mux@v1.8.0
- github.com/gorilla/websocket@v1.4.2
- github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0
- github.com/grpc-ecosystem/grpc-gateway@v1.16.0
- github.com/gsterjov/go-libsecret@v0.0.0-20161001094733-a6f4afe4910c
- github.com/gtank/merlin@v0.1.1
- github.com/gtank/ristretto255@v0.1.2
- github.com/hashicorp/go-immutable-radix@v1.0.0
- github.com/hashicorp/golang-lru@v0.5.4
- github.com/hashicorp/hcl@v1.0.0
- github.com/lib/pq@v1.2.0
- github.com/libp2p/go-buffer-pool@v0.0.2
- github.com/magiconair/properties@v1.8.5
- github.com/mattn/go-isatty@v0.0.12
- github.com/matttproud/golang_protobuf_extensions@v1.0.1
- github.com/mimoo/StrobeGo@v0.0.0-20181016162300-f8f6d4d2b643
- github.com/minio/highwayhash@v1.0.1
- github.com/mitchellh/mapstructure@v1.3.3
- github.com/mtibben/percent@v0.2.1
- github.com/pelletier/go-toml@v1.8.1
- github.com/pkg/errors@v0.9.1
- github.com/prometheus/client_golang@v1.10.0
- github.com/prometheus/client_model@v0.2.0
- github.com/prometheus/common@v0.23.0
- github.com/prometheus/procfs@v0.6.0
- github.com/rakyll/statik@v0.1.7
- github.com/rcrowley/go-metrics@v0.0.0-20200313005456-10cdbea86bc0
- github.com/regen-network/cosmos-proto@v0.3.1
- github.com/rs/cors@v1.7.0
- github.com/rs/zerolog@v1.21.0
- github.com/spf13/afero@v1.6.0
- github.com/spf13/cast@v1.3.1
- github.com/spf13/cobra@v1.1.3
- github.com/spf13/jwalterweatherman@v1.1.0
- github.com/spf13/pflag@v1.0.5
- github.com/spf13/viper@v1.7.1
- github.com/subosito/gotenv@v1.2.0
- github.com/syndtr/goleveldb@v1.0.1-0.20200815110645-5c35d600f0ca
- github.com/tendermint/btcd@v0.1.1
- github.com/tendermint/crypto@v0.0.0-20191022145703-50d29ede1e15
- github.com/tendermint/go-amino@v0.16.0
- github.com/tendermint/tendermint@v0.34.14
- github.com/tendermint/tm-db@v0.6.4
- github.com/zondax/hid@v0.9.0
- golang.org/x/crypto@v0.0.0-20220214200702-86341886e292
- golang.org/x/net@v0.0.0-20211112202133-69e39bad7dc2
- golang.org/x/sys@v0.0.0-20210917161153-d61c044b1678
- golang.org/x/term@v0.0.0-20201126162022-7de9c90e9dd1
- golang.org/x/text@v0.3.7
- google.golang.org/genproto@v0.0.0-20220228195345-15d65a4533f7
- google.golang.org/grpc@v1.44.0 => google.golang.org/grpc@v1.33.2
- google.golang.org/protobuf@v1.27.1
- gopkg.in/ini.v1@v1.61.0
- gopkg.in/yaml.v2@v2.4.0
cosmos_sdk_version: v0.42.11
```

### Development

The included Mainchain [DevNet](local-devnet.md) network can be used for development and testing of new features and 
bug fixes locally. To build the binaries for testing without installing, run:

```bash
make build
```

This will output the binaries to the `./build` directory located in the repository root.

See the [DevNet](local-devnet.md) docs for more information about running **DevNet**.

## CLI Help

Both the `und` and `und` commands can have the `--help` flag passed
to output details on what commands are available, and flags enabled for that
command:

```bash
und --help
```

Likewise, the `--help` flag can be passed to subcommands, for example:

```bash
und query wrkchain --help
```

#### Next

Running [Devnet](../networks/local-devnet.md), joining [a network](../networks/join-network.md)
