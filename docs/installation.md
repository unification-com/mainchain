# Installing the Mainchain Software

This documentation outlines how to install the UND Mainchain software, in
order to participate and interact with the Mainchain network.

## Install Go

**Go 1.13+** is required to build the Mainchain binaries

Install `go` by following the [official docs](https://golang.org/doc/install).
Once Go is installed, set your `$PATH` environment variable:

```bash
mkdir -p $HOME/go/bin
echo "export PATH=$PATH:$(go env GOPATH)/bin" >> ~/.bash_profile
source ~/.bash_profile
```

## Install the binaries

Once Go is installed, download the latest Mainchain release, and run:

```bash
cd mainchain
make install
```

This will install the two binaries `und` and `undcli` into your `$HOME/go/bin`

### Install with Ledger support

To install with Ledger support, run the following:

```bash
export LEDGER_ENABLED=true && make install
```

## Verify the installation

Run the following commands:

```bash
$ und version --long
$ undcli version --long
```

If they have installed correctly, you should see output similar to the following:

```json
{
  "name":"UndMainchain",
  "server_name":"und",
  "client_name":"undcli",
  "version":"1.0.0",
  "commit":"5797b5061b4035ec9d6818fef6a1a7967b4e2fba",
  "build_tags":"netgo ledger",
  "go":"go version go1.13.3 linux/amd64"
}
```

### Build Tags

Build tags indicate special features that have been enabled in the binary.

| Build Tag | Description                                     |
| --------- | ----------------------------------------------- |
| netgo     | Name resolution will use pure Go code           |
| ledger    | Ledger devices are supported (hardware wallets) |


## Development

The included Mainchain [Devnet](local-devnet.md) can be used for development
and testing of new features and bug fixes. For local binaries, use:

```bash
make clean && make build
```

### Build with Ledger support

To build with Ledger support, run the following:

```bash
export LEDGER_ENABLED=true && make clean && make build
```

This will output the binaries to the `./build` directory located in the repository 
root.

## Next

Running [Devnet](local-devnet.md), joining Testnet or Mainnet
