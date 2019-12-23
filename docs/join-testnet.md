# Join a Public Testnet

Once you have [installed](installation.md) the required software, you can join one of the 
public Testnets.

## Initialising a New Node

Once installed, you will need to initialise the config files for your new node:

```bash
und init [your_node_tag]
```

`[your_node_tag]` can be any identifier you like, but are limited to ASCII characters.

Once initialised, you can edit your configuration in `~/.und_mainchain/config/config.toml`

## Genesis

The latest Testnet genesis can always be found at https://github.com/unification-com/testnet/latest

To spin up your new Testnet node, download the latest `genesis.json`:

```bash
mkdir -p ~/.und_mainchain/config
curl https://raw.githubusercontent.com/unification-com/testnet/master/latest/genesis.json > ~/.und_mainchain/config/genesis.json
```

### Validate Genesis

You can validate the downloaded genesis file by running:

```bash
und validate-genesis ~/.und_mainchain/config/genesis.json
```

which should result in:

```bash
validating genesis file at ~/.und_mainchain/config/genesis.json
File at ~/.und_mainchain/config/genesis.json is a valid genesis file
```

## Seed Node Peers

Your node will need to know at least one seed node in order to join the network
and begin P2P communication. The latest seed information will always be available 
at https://github.com/unification-com/testnet/blob/master/latest/seed_nodes.md

Edit `~/.und_mainchain/config/config.toml`, and set the `persistent_peers` value with
a comma separated list of one or more peers. For example, the DevNet nodes
use the following:

```toml
persistent_peers = "53e857acc2df7127d5ef33b0dd98c55e7068ae06@172.25.0.4:26656,33a49c1eae31ce82ffab25ed821e8cec7f8bbd00@172.25.0.5:26656"
```

## Minimum Gas

In order to protect your full node from spam transactions, it is good practice to 
set the `minimum-gas-prices` value in `~/.und_mainchain/config/app.toml`. This should be
set as a decimal value, and the recommended value is currently `0.025nund`

## Running your node

Now that you have `genesis`, and some seed nodes, you can run your full node:

```bash
und start
```

Running:

```bash
undcli status
```

should output show that the node is running and connected.

By default, any transactions you send via the `undcli` command will be
sent via your local node (which was started using the `und start` command).
You can use the `--node` flag with the `undcli` command to have it send
to a public node instead.

## TestNet Faucet

Our public TestNet has a faucet which can be used to obtain Test UND for
use on the TestNet network. You will need an [account](accounts-wallets.md) and its associated 
address in order to be able to claim Test UND.

See https://faucet-testnet.unification.io

**Note**: You will need an account setting up before requesting Test UND.
See [accounts and wallets](accounts-wallets.md) for more details

## Next

Creating and importing [accounts and wallets](accounts-wallets.md), [sending transactions](transactions.md)
