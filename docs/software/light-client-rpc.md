# Light Client/REST

The `undcli` binary comes packaged with a full REST capable server, from which the majority of the `undcli query ...` and `undcli tx ...` commands can also be served.

The REST server is generally useful for third party services such as [wallets](https://github.com/unification-com/web-wallet) and [block explorers](https://github.com/unification-com/mainchain-explorer). It interacts with, and can be used alongside the `und` RPC interface.

#### Contents

[[toc]]

## Prerequisites

Before continuing, ensure you have gone through the following docs:

1. [Installing the software](installation.md)
2. Either [join TestNet](../networks/join-testnet.md), [join MainNet](../networks/join-mainnet.md) or [run DevNet](../networks/local-devnet.md)

## Running a light client

The Light Client can be started using the following command:

```bash
undcli rest-server --laddr=[tcp://ip:port] --node [tcp://ip:port] --chain-id=[chain_id]
```

For example:

```bash
undcli rest-server --laddr=tcp://localhost:1317 --node tcp://11.22.33.44:26657 --chain-id=UND-Mainchain-TestNet-v4
```

This will start the light client on your local host listening on `localhost:1317`, and use the node hosted at `11.22.33.44:26657` to source its data and interface with the `UND-Mainchain-TestNet-v4` chain (e.g. broadcast any transactions).

::: tip
setting the listen address IP to `0.0.0.0`, e.g. `--laddr=tcp://0.0.0.0:1317` will allow any host to connect to your REST server.
:::

Once running, you can visit [http://localhost:1317/swagger-ui/](http://localhost:1317/swagger-ui/) to view all of the REST endpoints available.

The full `undcli rest-server` command specification can be found [here](undcli-commands.md#undcli-rest-server).

## Running an Archive RPC node

Light Clients are more effective when interfacing with full nodes running in "archive" mode. Nodes running in archive mode do not prune any sync data, and keep a complete transaction event history.

The quickest way to get up and running with an archive node is to configure the pruning option in `$HOME/.und_mainchain/config/app.toml`:

```toml
pruning = "nothing"
```

Then, start the full node as usual using:

```bash
und start
```

Your light client can then be configured to connect to it via the `--node` flag by passing `tcp://127.0.0.1:26657` to it.
