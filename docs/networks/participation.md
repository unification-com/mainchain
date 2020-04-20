# Non-Validator Participation

There are many ways in which users can contribute to the running of the network without the need to become a full validator.

## Local full-nodes

One of the simplest ways to participate in the network is to run your own full node in which to send your own transactions to. Your node, by default, will only accept transactions sent to it from `localhost` before broadcasting them to its network peers, and is therefore a secure method by which to broadcast your transactions to the network. Further, by running your own local full node, transactions and their messages can be validated locally prior to being broadcast to network peers.

By default, your node's `P2P` address is not advertised, and its `RPC` is only accessible to the `localhost`.

This scenario is ideal for users who wish to infrequently spin up their node for the purpose of sending transactions to the network, or for users who are for example running a WRKChain, and wish to use the local node for broadcasting their WRKChain transactions.

## Seed Nodes

Seed nodes are used by full nodes to bootstrap their address book, by keeping a record of permanently connected nodes, and broadcasting their addresses on request. Seed nodes don't accept or broadcast transactions, and immediately disconnect from a peer once it has sent its address book to the connected peer.

### Configuring a seed node

In `$HOME/.und_mainchain/config/config.toml`:

```toml
p2p.persistent_peers = "[node_id_1]@[ip_1]:[port],[node_id_2]@[ip_2]:[port]" # List of peers known to the seed that are permanently available
p2p.seed_mode = true
```

## Archive Nodes

Archive nodes are full nodes that keep a complete history of the chain state, by not pruning any sync data (i.e. `pruning` is set to `nothing` in `$HOME/.und_mainchain/app.toml`). They are used as the data source for third party applications such as block explorers and wallet apps, since they keep a complete event and transaction history.

### Configuring an archive node

In `$HOME/.und_mainchain/config/config.toml`, the following suggested settings offer a basic starting point. Users are encouraged to explore the settings and configure as required to increase security and availability:

```toml
p2p.external_address = "[ip]:26656"
rpc.laddr = "tcp://0.0.0.0:26657"
tx_index.index_keys = ""
tx_index.index_all_keys = true
```

In `$HOME/.und_mainchain/config/app.toml`:

```toml
minimum-gas-prices = "0.25nund"
pruning = "nothing"
```

Ensure any firewall rules allow incoming requests to ports `26656` and `26657` from `0.0.0.0/0`.

### REST server

A REST (light-client) server can be run alongside an Archive Node to offer deeper querying capabilities to the network. Third party applications such as block explorers and wallet apps rely on REST servers along with Archive nodes for their data.

## Relay Nodes

A relay node is simply a full node that is always online and advertises its `P2P` network address to other peers. Relay nodes can potentially help reduce overall P2P network latency.

### Configuring a Relay Node

A relay node should have high availability (for example, be running on a Cloud VM, or other host that is always online with a static IP address).

In `$HOME/.und_mainchain/config/config.toml` set the `p2p.external_address` value to `[ip]:26656`, and ensure any firewall rules allow incoming requests from `0.0.0.0/0`

In `$HOME/.und_mainchin/config/app.toml` set the `minimum-gas-prices` value, to for example `0.25nund`.

## Invariance checking

Invariance checking can be done by any full-node and can help identify any malicious/odd activity on the chain db, and the overall data integrity. To do so, simply run the `und` node and pass the `--inv-check-period` flag:

```bash
und start --inv-check-period=1
```

Setting `--inv-check-period` to 1 will check every block.

Invariance checking is a resource intensive process.
