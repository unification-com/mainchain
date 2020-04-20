# Installing the Mainchain Software

This documentation outlines how to install the UND Mainchain software, in
order to participate and interact with any of the Mainchain networks.

#### Contents

[[toc]]

## Installing the latest release binaries

The latest pre-compiled binaries are available from [https://github.com/unification-com/mainchain/releases](https://github.com/unification-com/mainchain/releases).

- The `undcli` binary has been compiled for Linux, OSX and Windows.
- The `und` binary has been compiled for Linux only.

Simply download the archives for your OS, extract them and copy the binaries to a suitable location - preferably a location in your `$PATH` environment variable, for example `/usr/local/bin`, `/opt`, etc.

Once installed, verify:

```bash
$ und version --long
$ undcli version --long
```

The output should match the latest release version tag.

## Building from Source

The Mainchain binaries can also be built from source.

### Prerequisites

`git`, `curl` and `make` are required to build the binaries. `jq` is also useful for quickly looking up values in `genesis.json`

These can all be installed via your package manager:

```bash
sudo apt-get install git curl make jq
```

or

```bash
sudo yum install git curl make jq
```

### Install Go

**Go 1.13+** is required to build the Mainchain binaries

Install `go` by following the [official docs](https://golang.org/doc/install).
Once Go is installed, set your `$PATH` environment variable:

```bash
$ mkdir -p $HOME/go/bin
$ echo "export PATH=$PATH:$(go env GOPATH)/bin" >> ~/.bash_profile
$ source ~/.bash_profile
```

### Build and install the binaries

::: warning IMPORTANT
unless you are contributing to Mainchain development, it is recommended you checkout and build from the latest release tag and **not** the `master` branch if you intend to connect to a live, public network (e.g. TestNet/MainNet).
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

This will install the two binaries `und` and `undcli` into your `$HOME/go/bin`

### Verify the installation

Run the following commands:

```bash
$ und version --long
$ undcli version --long
```

If they have installed correctly, you should see output similar to the following:

```bash
name: UndMainchain
server_name: und
client_name: undcli
version: 1.3.4
commit: 6913d4e349ef99aef9be0dfe3c03e8381dae0d81
build_tags: netgo ledger
go: go version go1.13.3 linux/amd64
```

### Development

The included Mainchain [DevNet](local-devnet.md) network can be used for development and testing of new features and bug fixes locally. To build the binaries for testing without installing, run:

```bash
make build
```

This will output the binaries to the `./build` directory located in the repository root.

See the [DevNet](local-devnet.md) docs for more information about running **DevNet**.

## CLI Help

Both the `und` and `undcli` commands can have the `--help` flag passed
to output details on what commands are available, and flags enabled for that
command:

```bash
und --help
undcli --help
```

Likewise, the `--help` flag can be passed to subcommands, for example:

```bash
undcli query wrkchain --help
```

#### Next

Running [Devnet](networks/local-devnet.md), joining [Testnet](networks/join-testnet.md) or [MainNet](networks/join-mainnet.md)
