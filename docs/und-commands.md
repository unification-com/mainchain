# `und` & command overview

The `und` binary is the Mainchain server-side software, used to run a full-node and validator.

[[toc]]

## Commands

## Flags

## Ports overview

The two main ports used by `und` are to handle P2P communication between peer nodes, and to handle any RPC requests.

### P2P: 26656

By default, port `26656` is used by `und` for P2P communication. P2P communication occurs between peer nodes for example when passing blocks etc. to each other.

P2P configuration can be set in the `[p2p]` section of the `$HOME/.und_mainchain/config/config.toml` configuration file.

### RPC: 26657

By default, port `26657` is used by `und` for handling incoming RPC requests, including for example, Tx broadcasts from `undcli` and any chain queries. The default configuration restricts RPC access to `undcli` running on the same host (localhost).

RPC configuration can be set in the `[rpc]` section of the `$HOME/.und_mainchain/config/config.toml` configuration file.
