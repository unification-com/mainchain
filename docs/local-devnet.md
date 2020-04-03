# Deploying a Local DevNet

>**IMPORTANT**: Whenever you use `undcli` to send Txs or query the chain ensure you pass the correct data to the `--chain-id` and if necessary `--node=` flags so that you connect to the correct network!

The repository contains a ready to deploy Docker composition for local
development and testing. DevNet can be run with the following command:

```bash
make devnet
```

To bring DevNet down cleanly, use <kbd>Ctrl</kbd>+<kbd>C</kbd>, followed by:

```bash
make devnet-down
```

## DevNet Chain ID

**Important**: DevNet's Chain ID is `UND-Mainchain-DevNet`. Any `und` or `undcli` commands
intended for DevNet should use the flag `--chain-id UND-Mainchain-DevNet`

## DevNet RPC Nodes

By default `undcli` will attempt to broadcast transactions to tcp://localhost:26656. However, any of the DevNet nodes can be used to send transactions via `undcli` using the `--node=` flag, for example:

```bash
undcli query tx TX_HASH --chain-id UND-Mainchain-DevNet --node=tcp://172.25.0.3:26661
```

See below for each node's RPC IPs and Ports.

## DevNet Docker containers

The DevNet composition will spin up three full nodes, one light REST client, and a proxy server in the following Docker containers:

- `node1` - Full validation node, RPC on 172.25.0.3:26661
- `node2` - Full validation node, RPC on 172.25.0.4:26662
- `node3` - Full validation node, RPC on 172.25.0.5:26663
- `rest-server` - Light Client for REST interaction on 172.25.0.6:1317
- `proxy` - a small proxy server allowing CORS queries to the `rest-server` via 172.25.0.7:1318

## DevNet test accounts, wallets and keys

DevNet is deployed with a pre-defined [genesis.json](../Docker/assets/node1/config/genesis.json), containing several test accounts loaded with UND and pre-defined validators with self delegation.

See <a href="https://github.com/unification-com/mainchain/blob/master/Docker/README.md" target="_blank">https://github.com/unification-com/mainchain/blob/master/Docker/README.md</a> for the mnemonic phrases and keys used by the above nodes, and for test accounts included in DevNet's genesis.

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

## Next

Creating and importing [accounts and wallets](accounts-wallets.md), [sending transactions](transactions.md)
