# Running a `und` Full Node

Once you have [installed](../software/installation.md) the required software, you can run a full node

#### Contents

[[toc]]

## Prerequisites

Before continuing, ensure you have gone through the following docs:

1. [Installing the software](../software/installation.md)

## Initialising a New Node

Once installed, you will need to initialise your node:

```bash
und init [your_node_moniker]
```

`[your_node_moniker]` can be any identifier you like, but are limited to ASCII characters. For example:

```bash
und init MyAwesomeNode
```

Once initialised, you can edit your configuration in `$HOME/.und_mainchain/config/config.toml`. See [configuration reference](../software/und-mainchain-config-ref.md) for more details on the config file.

## Genesis

Depending on the network you are joining, you will need the appropriate genesis file to replace the default `$HOME/.und_mainchain/config/genesis.json`
