# Light Client/REST

The `und` binary comes packaged with a full REST capable server, from which the majority of 
the `und query ...` and `und tx ...` commands can also be served.

The REST server is generally useful for third party services such as 
[wallets](https://github.com/unification-com/web-wallet) and 
[block explorers](https://github.com/unification-com/mainchain-explorer). It interacts with, and can be used 
alongside the `und` RPC interface.

#### Contents

[[toc]]

## Prerequisites

Before continuing, ensure you have gone through the following docs:

1. [Installing the software](installation.md)
2. [join a Network](../networks/join-network.md), or [run DevNet](../networks/local-devnet.md)

## Running a light client

The Light Client can be started by setting the configuring options in `app.toml` as follows:

```toml
[api]
enable = true
swagger = true
address = "tcp://0.0.0.0:1317"
```

Then running the `und start` command as normal

```bash
und start
```

This will start the light client on your local host listening on `localhost:1317`, and use the node 
hosted at `11.22.33.44:26657` to source its data and interface with the `FUND-Mainchain-TestNet-v8` chain 
(e.g. broadcast any transactions).

::: tip
setting the listen address IP to `0.0.0.0`, e.g. `tcp://0.0.0.0:1317` will allow any host to connect to your REST server.
:::

Once running, you can visit [http://localhost:1317/swagger/](http://localhost:1317/swagger/) to view all the REST 
endpoints available.

## Running an Archive RPC node

Light Clients are more effective when interfacing with full nodes running in "archive" mode. Nodes running in archive 
mode do not prune any sync data, and keep a complete transaction event history.

The quickest way to get up and running with an archive node is to configure the pruning option in 
`$HOME/.und_mainchain/config/app.toml`:

```toml
pruning = "nothing"
```

Then, start the full node as usual using:

```bash
und start
```

Your light client can then be configured to connect to it via the `--node` flag by passing `tcp://127.0.0.1:26657` to it.
