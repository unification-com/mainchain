# Deploying a Local DevNet

::: warning IMPORTANT
Whenever you use `undcli` to send Txs or query the chain ensure you pass the correct data to the `--chain-id` and if necessary `--node=` flags so that you connect to the correct network!
:::

The repository contains a ready to deploy Docker composition for local
development and testing. DevNet comes in two flavours - `local` and `upstream`.

#### Contents

[[toc]]

## Local build

The local build copies the current local codebase to the Docker containers, and is used during development to test changes before committing to the repository.

```
docker-compose -f Docker/docker-compose.local.yml up --build
docker-compose -f Docker/docker-compose.local.yml down --remove-orphans
```

or using the `make` target:

```bash
make devnet
```

To bring DevNet down cleanly, use <kbd>Ctrl</kbd>+<kbd>C</kbd>, followed by:

```bash
make devnet-down
```

## Pure Upstream build

Pure upstream downloads the `master` branch on GitHub to build the binaries, and is useful for testing the latest code committed to `master`, for example for pre-release testing.

```
docker-compose -f Docker/docker-compose.upstream.yml up --build
docker-compose -f Docker/docker-compose.upstream.yml down --remove-orphans
```

or using the `make` target:

```bash
make devnet-pristine
```

To bring DevNet down cleanly, use <kbd>Ctrl</kbd>+<kbd>C</kbd>, followed by:

```bash
make devnet-pristine-down
```

## DevNet Chain ID

::: warning IMPORTANT
DevNet's Chain ID is `FUND-Mainchain-DevNet`. Any `und` or `undcli` commands
intended for DevNet should use the flag `--chain-id FUND-Mainchain-DevNet`
:::

## DevNet RPC Nodes

By default `undcli` will attempt to broadcast transactions to tcp://localhost:26656. However, any of the DevNet nodes can be used to send transactions via `undcli` using the `--node=` flag, for example:

```bash
undcli query tx TX_HASH --chain-id FUND-Mainchain-DevNet --node=tcp://172.25.0.3:26661
```

See below for each node's RPC IPs and Ports.

## DevNet Docker containers

The DevNet composition will spin up three full nodes, one light REST client, and a proxy server in the following Docker containers:

- `node1` - Full validation node, RPC on 172.25.0.3:26661, P2P on 172.25.0.3:26651
- `node2` - Full validation node, RPC on 172.25.0.4:26662, P2P on 172.25.0.4:26652
- `node3` - Full validation node, RPC on 172.25.0.5:26663, P2P on 172.25.0.5:26653
- `rest-server` - Light Client for REST interaction on 172.25.0.6:1317
- `proxy` - a small proxy server allowing CORS queries to the `rest-server` via 172.25.0.7:1318

::: tip NOTE
The DevNet nodes:  
P2P ports set to 26651, 26652 and 26653 respectively, and not the default 26656.  
RPC ports set to 26661, 26662 and 26663 respectively, and not the default 26657.
:::

## DevNet test accounts, wallets and keys

DevNet is deployed with a pre-defined [genesis.json](https://raw.githubusercontent.com/unification-com/mainchain/master/Docker/assets/node1/config/genesis.json), containing several test accounts loaded with FUND and pre-defined validators with self delegation.

See [https://github.com/unification-com/mainchain/blob/master/Docker/README.md](https://github.com/unification-com/mainchain/blob/master/Docker/README.md) for the mnemonic phrases and keys used by the above nodes, and for test accounts included in DevNet's genesis.

### Importing the DevNet keys

The DevNet accounts can be imported as follows. First, build the `und` and
`undcli` binaries:

```bash
make build
```

Then, for each account run the following command:

```bash
./build/undcli keys add node1 --recover
```

You will be prompted to enter the mnemonic phrase, and a password for your OS's keyring. Change `node1` to an appropriate moniker for each imported account.

### Useful DevNet Defaults for `undcli`

`undcli` defaults for DevNet can be set as follows. This will set the corresponding values in `$HOME/.und_cli/config/config.toml`

```
undcli config chain-id FUND-Mainchain-DevNet
undcli config node tcp://localhost:26661
```

### REST API Endpoints

With DevNet up, the REST API endpoints can be seen via [http://localhost:1318/swagger-ui/](http://localhost:1318/swagger-ui/)

#### Next

Creating and importing [accounts and wallets](accounts-wallets.md), [sending transactions](transactions.md)
