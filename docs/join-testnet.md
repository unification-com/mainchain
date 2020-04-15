# Join the Public TestNet

Once you have [installed](installation.md) the required software, you can run a full node, join the public TestNet and try out [becoming a TestNet validator](become-testnet-validator.md).

::: warning IMPORTANT
Whenever you use `undcli` to send Txs or query the chain ensure you pass the correct data to the `--chain-id` and if necessary `--node=` flags so that you connect to the correct network!
:::

#### Contents

[[toc]]

## Prerequisites

Before continuing, ensure you have gone through the following docs:

1. [Installing the software](installation.md)

## Initialising a New Node

Once installed, you will need to initialise your node:

```bash
und init [your_node_moniker]
```

`[your_node_moniker]` can be any identifier you like, but are limited to ASCII characters. For example:

```bash
und init MyAwesomeNode
```

Once initialised, you can edit your configuration in `$HOME/.und_mainchain/config/config.toml`

::: tip NOTE
the default directory used by `und` is `$HOME/.und_mainchain`. This can be changed by passing the global `--home=` flag to the `und` command, for example `und start --home=$HOME/.und_mainchain_TestNet`.
:::

## Genesis

The latest TestNet genesis can always be found at [https://github.com/unification-com/testnet/latest](https://github.com/unification-com/testnet/latest)

### Download the latest Genesis

To spin up your new TestNet node, download the latest `genesis.json`:

```bash
curl https://raw.githubusercontent.com/unification-com/testnet/master/latest/genesis.json > $HOME/.und_mainchain/config/genesis.json
```

::: warning IMPORTANT
remember to change the output directory if you are using something other than the default `$HOME/.und_mainchain`
:::

### Validate Genesis

Optionally, you can validate the downloaded genesis file by running:

```bash
und validate-genesis $HOME/.und_mainchain/config/genesis.json
```

which should result in:

```bash
validating genesis file at $HOME/.und_mainchain/config/genesis.json
File at $HOME/.und_mainchain/config/genesis.json is a valid genesis file
```

### Get the current TestNet chain ID

The Chain ID will need to be passed to all `undcli` commands via the `--chain-id` flag. The current TestNet Chain ID can easily be found by running:

```bash
jq --raw-output '.chain_id' $HOME/.und_mainchain/config/genesis.json
```

This will output, for example:

```
UND-Mainchain-TestNet-v3
```

which can then be passed to your `undcli` commands:

```bash
undcli query tx TX_HASH --chain-id UND-Mainchain-TestNet-v3
```

## Seed Node Peers

Your node will need to know at least one seed node in order to join the network
and begin P2P communication with other nodes in the network. The latest seed information will always be available at [https://github.com/unification-com/testnet/blob/master/latest/seed_nodes.md](https://github.com/unification-com/testnet/blob/master/latest/seed_nodes.md)

Edit `$HOME/.und_mainchain/config/config.toml`, and set the `persistent_peers` value with a comma separated list of one or more peers. For example, a TestNet seed node:

```toml
persistent_peers = "dcff5de69dcc170b28b6628a1336d420f7eb60c0@seed1-testnet.unification.io:26656"
```

::: warning IMPORTANT
always check the latest TestNet seed node in the repository - the above example may not always match the actual current seed node!
:::

## Minimum Gas

In order to protect your full node from spam transactions, it is good practice to set the `minimum-gas-prices` value in `$HOME/.und_mainchain/config/app.toml`. This should be set as a decimal value, and the recommended value for **TestNet** is currently `0.025nund` to `0.25nund`.

## Running your node

Now that you have `genesis`, and some seed nodes, you can run your full node:

```bash
und start
```

You should see that your node connects to some peers, and after a few seconds begins syncing with the network.

Running:

```bash
undcli status
```

in a separate terminal should output show that the node is running and connected.

By default, any transactions you send via the `undcli` command will be
sent via your local node (which was started using the `und start` command, and whose RPC is on `tcp://localhost:26656`).

::: tip
You can use the `--node` flag with the `undcli` command to have it send to a different node instead.
:::

## Invariance checking

You don't need to become a validator to take part in the network - just running a full node as a p2p peer is very useful. Another method to help the network, is invariance checking. This can help the network by periodically checking blocks for invariances which could potentially cause issues.

::: tip NOTE
Invariance checking is resource intensive, so should not be invoked on validator nodes!
:::

Start a full node with the `--inv-check-period` flag. Value of 1 will
check every block for invariances:

```
und start --inv-check-period 1
```

Invariance Tx can be sent using:

```
undcli tx crisis invariant-broken enterprise module-account --from wrktest
```

## TestNet Faucet

Our public TestNet has a faucet which can be used to obtain Test UND for
use exclusively on the TestNet network. You will need an [account](accounts-wallets.md) and its associated address in order to be able to claim Test UND.

See [https://faucet-testnet.unification.io](https://faucet-testnet.unification.io)

::: tip NOTE
You will need an account setting up before requesting Test UND.
See [accounts and wallets](accounts-wallets.md) for more details
:::

## TestNet Explorer

Our public TestNet explorer can be found at [https://explorer-testnet.unification.io](https://explorer-testnet.unification.io)

#### Next

Creating and importing [accounts and wallets](accounts-wallets.md), [sending transactions](examples/transactions.md) and [becoming a TestNet validator](become-testnet-validator.md)
