# `undcli` command overview

The `undcli` binary is the primary CLI client tool used for interacting with a full `und` node. The `und` node can be running locally, or one being run as a public service. By default, `undcli` will assume the `und` node is running locally and attempt to connect via RPC to `tcp://localhost:26657`.

[[toc]]

## Commands

### Query account

```
undcli query account und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
```

## Flags

#### Next

More detailed examples for specific modules can be found in [Transaction Examples](examples/transactions.md), [WRKChain Examples](examples/wrkchain.md), [BEACON examples](examples/beacon.md) and [Enterprise UND Examples](examples/enterprise-und.md).
