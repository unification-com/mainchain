# `und` Command Reference

The `und` binary is the Mainchain server-side software, used to run a full-node and validator.

#### Contents

[[toc]]

## und

Unification Mainchain Daemon (server)

Usage:
```bash
  und [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[init](#und-init)|Initialise private validator, p2p, genesis, and application configuration files|
|[validate-genesis](#und-validate-genesis)|validates the genesis file at the default location or at the location passed as an arg|
|[debug](#und-debug)|Tool for helping with debugging your application|
|[start](#und-start)|Run the full node|
|[unsafe-reset-all](#und-unsafe-reset-all)|Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state|
|[tendermint](#und-tendermint)|Tendermint subcommands|
|[export](#und-export)|Export state to JSON|
|[version](#und-version)|Print the app version|
|[help](#und-help)|Help about any command|

**Global Flags**

::: tip
Global flags can be passed to any of the commands and sub-commands below
:::

| Flag | Type | Description |
|------|------|-------------|
|`--home`|`string`|directory for config and data (default "/home/username/.und_mainchain")|
|`--inv-check-period`|`uint`|Assert registered invariants every N blocks |
|`--log_level`|`string`|Log level (default `"main:info,state:info,*:error"`)|  
|`--trace`||print out full stack trace on errors|

## und init

Initialise a validator and node's configuration files.

Usage:
```bash
  und init [moniker] [flags]
```

Example:
```bash
  und init MyAwesomeNode --chain-id="UND-Mainchain-TestNet"
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|genesis file chain-id, if left blank will be randomly created|
|`-h`, `--help`||help for init|
|`-o`, `--overwrite`||overwrite the genesis.json file|

## und start

Run the full node application with Tendermint in or out of process. By
default, the application will run with Tendermint in process.

Pruning options can be provided via the `--pruning` flag. The pruning options are as follows:

`syncable`: only those states not needed for state syncing will be deleted (flushes every 100th to disk and keeps every 10000th)  
`nothing`: all historic states will be saved, nothing will be deleted (i.e. archiving node)  
`everything`: all saved states will be deleted, storing only the current state

Node halting configurations exist in the form of two flags: `--halt-height` and `--halt-time`. During the ABCI Commit phase, the node will check if the current block height is greater than or equal to the `halt-height` or if the current block time is greater than or equal to the `halt-time`. If so, the node will attempt to gracefully shutdown and the block will not be committed. In addition, the node will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the
`--cpu-profile` flag which accepts a path for the resulting `pprof` file.

Usage:
```bash
  und start [flags]
```

Flags:

Additionally, see `$HOME/.und_mainchain/config/config.toml` and `$HOME/.und_mainchain/config/app.toml` where many of these values are set by default.

| Flag | Type | Description |
|------|------|-------------|
|`--abci`|`string`|Specify abci transport (`socket` \| `grpc`) (default "socket")|
|`--address`|`string`|Listen address (default "tcp://0.0.0.0:26658")|
|`--consensus.create_empty_blocks`||Set this to false to only produce blocks when there are txs or when the AppHash changes (default true)|
|`--consensus.create_empty_blocks_interval`|`string`|The possible interval between empty blocks (default "0s")|
|`--cpu-profile`|`string`|Enable CPU profiling and write to the provided file|
|`--db_backend`|`string`|Database backend: goleveldb \| cleveldb \| boltdb \| rocksdb (default "goleveldb")|
|`--db_dir`|`string`|Database directory (default "data")|
|`--fast_sync`||Fast blockchain syncing (default true)|
|`--genesis_hash`|`bytesHex`|Optional SHA-256 hash of the genesis file|
|`--halt-height`|`uint`|Block height at which to gracefully halt the chain and shutdown the node|
|`--halt-time`|`uint`|Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node|
|`-h`, `--help`||help for start|
|`--inter-block-cache`||Enable inter-block caching (default true)|
|`--minimum-gas-prices`|`string`|Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01nund;0.0001nund)|
|`--moniker`|`string`|Node Name|
|`--p2p.laddr`|`string`|Node listen address. (0.0.0.0:0 means any interface, any port) (default "tcp://0.0.0.0:26656")|
|`--p2p.persistent_peers`|`string`|Comma-delimited ID@host:port persistent peers|
|`--p2p.pex`||Enable/disable Peer-Exchange (default true)|
|`--p2p.private_peer_ids`|`string`|Comma-delimited private peer IDs|
|`--p2p.seed_mode`||Enable/disable seed mode|
|`--p2p.seeds`|`string`|Comma-delimited ID@host:port seed nodes|
|`--p2p.unconditional_peer_ids`|`string`|Comma-delimited IDs of unconditional peers|
|`--p2p.upnp`||Enable/disable UPNP port forwarding|
|`--priv_validator_laddr`|`string`|Socket address to listen on for connections from external priv_validator process|
|`--proxy_app`|`string`|Proxy app address, or one of: 'kvstore', 'persistent_kvstore', 'counter', 'counter_serial' or 'noop' for local testing. (default "tcp://127.0.0.1:26658")|
|`--pruning`|`string`|Pruning strategy: syncable, nothing, everything (default "syncable")|
|`--rpc.grpc_laddr`|`string`|GRPC listen address (BroadcastTx only). Port required|
|`--rpc.laddr`|`string`|RPC listen address. Port required (default "tcp://127.0.0.1:26657")|
|`--rpc.unsafe`||Enabled unsafe rpc methods|
|`--trace-store`|`string`|Enable KVStore tracing to an output file|
|`--with-tendermint`||Run abci app embedded in-process with tendermint (default true)|

## und unsafe-reset-all

Resets the blockchain database, removes address book files, and resets `priv_validator.json` to the genesis state

Usage:
```bash
  und unsafe-reset-all [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for unsafe-reset-all|

## und tendermint

Tendermint subcommands

Usage:
```bash
  und tendermint [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[show-node-id](#und-tendermint-show-node-id)|Show this node's ID|
|[show-validator](#und-tendermint-show-validator)|Show this node's tendermint validator info|
|[show-address](#und-tendermint-show-address)|Shows this node's tendermint validator consensus address|
|[version](#und-tendermint-version)|Print tendermint libraries' version|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for tendermint|

Use "`und tendermint [command] --help`" for more information about a command.

## und tendermint show-node-id

Show this node's ID

Usage:
```bash
  und tendermint show-node-id [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for show-node-id|

## und tendermint show-validator

Show this node's tendermint validator info

Usage:
```bash
  und tendermint show-validator [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for show-validator|
|`-o`, `--output`|`string`|Output format (text\|json) (default "text")|

## und tendermint show-address

Shows this node's tendermint validator consensus address

Usage:
```bash
  und tendermint show-address [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for show-address|
|`-o`, `--output`|`string`|Output format (text\|json) (default "text")|

## und tendermint version

Print protocols' and libraries' version numbers
against which this app has been compiled.

Usage:
```bash
  und tendermint version [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for version|

## und export

Export state to JSON

Usage:
```bash
  und export [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--for-zero-height`||Export state to start at height zero (perform preproccessing)|
|`--height`|`int`|Export state from a particular height (-1 means latest height) (default -1)|
|`-h`, `--help`||help for export|
|`--jail-whitelist`|`strings`|List of validators to not jail state export|

## und version

Print the und version

Usage:
```bash
  und version [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for version|
|`--long`||Print long version information|

## und debug

Tool for helping with debugging your application

Usage:
```bash
  und debug [flags]
  und debug [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[pubkey](#und-debug-pubkey)|Decode a ED25519 pubkey from hex, base64, or bech32|
|[addr](#und-debug-addr)|Convert an address between hex and bech32|
|[raw-bytes](#und-debug-raw-bytes)|Convert raw bytes output (eg. [10 21 13 255]) to hex|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for version|

## und debug pubkey

Decode a pubkey from hex, base64, or bech32.

Usage:
```bash
  und debug pubkey [pubkey] [flags]
```

Example:
```bash
  $ und debug pubkey TWFuIGlzIGRpc3Rpbmd1aXNoZWQsIG5vdCBvbmx5IGJ5IGhpcyByZWFzb24sIGJ1dCBieSB0aGlz
  $ und debug pubkey und1hp2km26czxlvesn8nmwswdd90umvcm5gxwpk98
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for pubkey|

## und debug addr

Convert an address between hex encoding and bech32.

Usage:
```bash
  und debug addr [address] [flags]
```

Example:
```bash
  und debug addr und1hp2km26czxlvesn8nmwswdd90umvcm5gxwpk98
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for addr|

## und debug raw-bytes

Convert raw-bytes to hex.

Usage:
```bash
  und debug raw-bytes [raw-bytes] [flags]
```

Example:
```bash
  und debug raw-bytes "[72 101 108 108 111 44 32 112 108 97 121 103 114 111 117 110 100]"
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for raw-bytes|

## und validate-genesis

validates the genesis file at the default location or at the location passed as an argument

Usage:
```bash
  und validate-genesis [file] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for validate-genesis|

## Ports overview

The major ports used by `und` are to handle P2P communication between peer nodes, and to handle any RPC requests.

### P2P: 26656

By default, port `26656` is used by `und` for P2P communication. P2P communication occurs between peer nodes for example when passing blocks etc. to each other.

P2P configuration can be set in the `[p2p]` section of the `$HOME/.und_mainchain/config/config.toml` configuration file.

### RPC: 26657

By default, port `26657` is used by `und` for handling incoming RPC requests, including for example, Tx broadcasts from `undcli` and any chain queries. The default configuration restricts RPC access to `undcli` running on the same host (localhost).

RPC configuration can be set in the `[rpc]` section of the `$HOME/.und_mainchain/config/config.toml` configuration file.
