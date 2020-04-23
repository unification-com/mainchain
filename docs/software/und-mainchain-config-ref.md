# `.und_mainchain/config/config.toml` Reference

The `$HOME/.und_mainchain/config/config.toml` file contains all the configuration options for the `und` server binary. Below is a reference for the file.

::: tip
Any path below can be absolute (e.g. "`/var/myawesomeapp/data`") or
relative to the home directory (e.g. "`data`"). The home directory is
"`$HOME/.und_mainchain`" by default, but could be changed via or `--home` und flag.
:::

#### Contents

[[toc]]

## Main base config options

### proxy_app

TCP or UNIX socket address of the ABCI application, or the name of an ABCI application compiled in with the Tendermint binary

Example

```toml
proxy_app = "tcp://127.0.0.1:26658"
```

### moniker

A custom human readable name for this node

Example

```toml
moniker = "test_add_acc"
```

### fast_sync

If this node is many blocks behind the tip of the chain, `FastSync` allows them to catchup quickly by downloading blocks in parallel and verifying their commits.

Example

```toml
fast_sync = true
```

### db_backend

Database backend: `goleveldb` | `cleveldb` | `boltdb` | `rocksdb`

* `goleveldb` ([github.com/syndtr/goleveldb](https://github.com/syndtr/goleveldb) - most popular implementation)
  - pure go
  - stable
* `cleveldb` (uses `levigo` wrapper)
  - fast
  - requires `gcc`
  - use `cleveldb` build tag (`go build -tags cleveldb`)
* `boltdb` (uses `etcd`'s fork of bolt - [github.com/etcd-io/bbolt](https://github.com/etcd-io/bbolt))
  - EXPERIMENTAL
  - may be faster is some use-cases (random reads - indexer)
  - use `boltdb` build tag (`go build -tags boltdb`)
* `rocksdb` (uses [github.com/tecbot/gorocksdb](https://github.com/tecbot/gorocksdb))
  - EXPERIMENTAL
  - requires `gcc`
  - use `rocksdb` build tag (`go build -tags rocksdb`)

Example

```toml
db_backend = "goleveldb"
```

### db_dir

Database directory - i.e. `$HOME/.und_mainchain/[db_dir]`

Example

```toml
db_dir = "data"
```

### log_level

Output level for logging, including package level options

Example

```toml
log_level = "main:info,state:info,*:error"
```

### log_format

Output format: '`plain`' (coloured text) or '`json`'

Example

```toml
log_format = "plain"
```

## Additional base config options

### genesis_file

Path to the JSON file containing the initial validator set and other meta data, i.e. `$HOME/.und_mainchain/[genesis_file]`

Example

```toml
genesis_file = "config/genesis.json"
```

### priv_validator_key_file

Path to the JSON file containing the private key to use as a validator in the consensus protocol, i.e. `$HOME/.und_mainchain/[priv_validator_key_file]`

Example

```toml
priv_validator_key_file = "config/priv_validator_key.json"
```

### priv_validator_state_file

Path to the JSON file containing the last sign state of a validator, i.e. `$HOME/.und_mainchain/[priv_validator_state_file]`

Example

```toml
priv_validator_state_file = "data/priv_validator_state.json"
```

### priv_validator_laddr

TCP or UNIX socket address for Tendermint to listen on for connections from an external PrivValidator process

Example

```toml
priv_validator_laddr = ""
```

### node_key_file

Path to the JSON file containing the private key to use for node authentication in the p2p protocol, i.e. `$HOME/.und_mainchain/[node_key_file]`

Example

```toml
node_key_file = "config/node_key.json"
```

### abci

Mechanism to connect to the ABCI application: `socket` | `grpc`

Example

```toml
abci = "socket"
```

### prof_laddr

TCP or UNIX socket address for the profiling server to listen on

Example

```toml
prof_laddr = "localhost:6060"
```

### filter_peers

If true, query the ABCI app on connecting to a new peer so the app can decide if we should keep the connection or not

Example

```toml
filter_peers = false
```

## Advanced configuration options

## rpc server configuration options

Configuration options in the `[rpc]` section of the `config.toml` file.

### rpc.laddr

TCP or UNIX socket address for the RPC server to listen on

Example

```toml
laddr = "tcp://127.0.0.1:26657"
```

### rpc.cors_allowed_origins

A list of origins a cross-domain request can be executed from Default value '`[]`' disables cors support  

Use '`["*"]`' to allow any origin

Example

```toml
cors_allowed_origins = []
```

### rpc.cors_allowed_methods

A list of methods the client is allowed to use with cross-domain requests

Example

```toml
cors_allowed_methods = ["HEAD", "GET", "POST", ]
```

### rpc.cors_allowed_headers

A list of non simple headers the client is allowed to use with cross-domain requests

Example

```toml
cors_allowed_headers = ["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time", ]
```

### rpc.grpc_laddr

TCP or UNIX socket address for the gRPC server to listen on

::: tip NOTE
The gRPC server only supports` /broadcast_tx_commit`
:::

Example

```toml
grpc_laddr = ""
```

### rpc.grpc_max_open_connections

Maximum number of simultaneous gRPC connections.

Does not include RPC (HTTP & WebSocket) connections. See `max_open_connections`

If you want to accept a larger number than the default, make sure you increase your OS limits.

0 - unlimited.

Should be < `{ulimit -Sn} - {MaxNumInboundPeers} - {MaxNumOutboundPeers} - {N of wal, db and other open files}`

E.g. 1024 - 40 - 10 - 50 = 924 = ~900

Example

```toml
grpc_max_open_connections = 900
```

### rpc.unsafe

Activate unsafe RPC commands like `/dial_seeds` and `/unsafe_flush_mempool`

Example

```toml
unsafe = false
```

### rpc.max_open_connections

Maximum number of simultaneous connections (including WebSocket).

Does not include gRPC connections. See `grpc_max_open_connections`

If you want to accept a larger number than the default, make sure you increase your OS limits.

0 - unlimited.

Should be < `{ulimit -Sn} - {MaxNumInboundPeers} - {MaxNumOutboundPeers} - {N of wal, db and other open files}`

E.g. 1024 - 40 - 10 - 50 = 924 = ~900

Example

```toml
max_open_connections = 900
```

### rpc. max_subscription_clients

Maximum number of unique clientIDs that can `/subscribe`

If you're using `/broadcast_tx_commit`, set to the estimated maximum number
of `broadcast_tx_commit` calls per block.

Example

```toml
max_subscription_clients = 100
```

### rpc.max_subscriptions_per_client

Maximum number of unique queries a given client can `/subscribe` to

If you're using GRPC (or Local RPC client) and `/broadcast_tx_commit`, set to
the estimated maximum number of `broadcast_tx_commit` calls per block.

Example

```toml
max_subscriptions_per_client = 5
```

### rpc.timeout_broadcast_tx_commit

How long to wait for a tx to be committed during `/broadcast_tx_commit`.

::: warning
Using a value larger than 10s will result in increasing the
global HTTP write timeout, which applies to all connections and endpoints.
:::

See https://github.com/tendermint/tendermint/issues/3435

Example

```toml
timeout_broadcast_tx_commit = "10s"
```

### rpc.max_body_bytes

Maximum size of request body, in bytes

Example

```toml
max_body_bytes = 1000000
```

### rpc.max_header_bytes

Maximum size of request header, in bytes

Example

```toml
max_header_bytes = 1048576
```

### rpc.tls_cert_file

The path to a file containing certificate that is used to create the HTTPS server.

Might be either absolute path or path related to tendermint's config directory.

If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.

::: tip Note
both `tls_cert_file` and `tls_key_file` must be present for Tendermint to create HTTPS server.
:::

Otherwise, HTTP server is run.

Example

```toml
tls_cert_file = ""
```

### rpc.tls_key_file

The path to a file containing matching private key that is used to create the HTTPS server.

Might be either absolute path or path related to tendermint's config directory.

::: tip Note
both `tls_cert_file` and `tls_key_file` must be present for Tendermint to create HTTPS server.
:::

Otherwise, HTTP server is run.

Example

```toml
tls_key_file = ""
```

## peer to peer (p2p) server configuration options

Configuration options in the `[p2p]` section of the `config.toml` file.

### p2p.laddr

Address to listen for incoming connections

Example

```toml
laddr = "tcp://0.0.0.0:26656"
```

### p2p.external_address

Address to advertise to peers for them to dial

If empty, will use the same port as the laddr, and will introspect on the listener or use UPnP to figure out the address.

Example

```toml
external_address = "11.22.33.44:26656"
```

### p2p.seeds

Comma separated list of seed nodes to connect to

Example

```toml
seeds = "dcff5de69dcc170b28b6628a1336d420f7eb60c0@seed1-testnet.unification.io:26656"
```

### p2p.persistent_peers

Comma separated list of nodes to keep persistent connections to

Example

```toml
persistent_peers = "3da95f113600fc324ecf759915993c13d701ed80@172.25.0.3:26656,53e857acc2df7127d5ef33b0dd98c55e7068ae06@172.25.0.4:26656"
```

### p2p.upnp

UPNP port forwarding

Example

```toml
upnp = false
```

### p2p.addr_book_file

Path to address book, i.e. `$HOME/.und_mainchain/[addr_book_file]`

Example

```toml
addr_book_file = "config/addrbook.json"
```

### p2p.addr_book_strict

Set `true` for strict address routability rules

Set `false` for private or local networks, or when adding non-routable (e.g. private subnet) IPs as peers etc.

Example

```toml
addr_book_strict = true
```

### p2p.max_num_inbound_peers

Maximum number of inbound peers

Example

```toml
max_num_inbound_peers = 40
```

### p2p.max_num_outbound_peers

Maximum number of outbound peers to connect to, excluding persistent peers

Example

```toml
max_num_outbound_peers = 10
```

### p2p.unconditional_peer_ids

List of node IDs, to which a connection will be (re)established ignoring any existing limits

Example

```toml
unconditional_peer_ids = "3da95f113600fc324ecf759915993c13d701ed80"
```

### p2p.persistent_peers_max_dial_period

Maximum pause when redialing a persistent peer (if zero, exponential backoff is used)

Example

```toml
persistent_peers_max_dial_period = "0s"
```

### p2p.flush_throttle_timeout

Time to wait before flushing messages out on the connection

Example

```toml
flush_throttle_timeout = "100ms"
```

### p2p.max_packet_msg_payload_size

Maximum size of a message packet payload, in bytes

Example

```toml
max_packet_msg_payload_size = 1024
```

### p2p.send_rate

Rate at which packets can be sent, in bytes/second

Example

```toml
send_rate = 5120000
```

### p2p.recv_rate

Rate at which packets can be received, in bytes/second

Example

```toml
recv_rate = 5120000
```

### p2p.pex

Set true to enable the peer-exchange reactor

Example

```toml
pex = true
```

### p2p.seed_mode

Seed mode, in which node constantly crawls the network and looks for peers. If another node asks it for addresses, it responds and disconnects.

Does not work if the peer-exchange reactor is disabled.

Example

```toml
seed_mode = false
```

### p2p.private_peer_ids

Comma separated list of peer IDs to keep private (will not be gossiped to other peers)

Example

```toml
private_peer_ids = "3da95f113600fc324ecf759915993c13d701ed80"
```

### p2p.allow_duplicate_ip

Toggle to disable guard against peers connecting from the same ip.

Example

```toml
allow_duplicate_ip = false
```

### p2p.handshake_timeout
### p2p.dial_timeout

Peer connection configuration.

Example
```toml
handshake_timeout = "20s"
dial_timeout = "3s"
```

## mempool configuration options

Configuration options in the `[mempool]` section of the `config.toml` file.

```toml
recheck = true
broadcast = true
wal_dir = ""
```

### mempool.size

Maximum number of transactions in the mempool

Example

```toml
size = 5000
```

### mempool.max_txs_bytes

Limit the total size of all txs in the mempool.

This only accounts for raw transactions (e.g. given 1MB transactions and `max_txs_bytes=5MB`, mempool will only accept 5 transactions).

Example

```toml
max_txs_bytes = 1073741824
```

### mempool.cache_size

Size of the cache (used to filter transactions we saw earlier) in transactions

Example

```toml
cache_size = 10000
```

### mempool.max_tx_bytes

Maximum size of a single transaction.

::: tip NOTE
the max size of a tx transmitted over the network is `{max_tx_bytes} + {amino overhead}`.
:::

Example

```toml
max_tx_bytes = 1048576
```

## fast sync configuration options

Configuration options in the `[fastsync]` section of the `config.toml` file.

### fastsync.version

Fast Sync version to use:
 1) "v0" (default) - the legacy fast sync implementation
 2) "v1" - refactor of v0 version for better testability

Example

```toml
version = "v0"
```

## consensus configuration options

Configuration options in the `[consensus]` section of the `config.toml` file.

```toml
wal_file = "data/cs.wal/wal"
```

```toml
timeout_propose = "3s"
timeout_propose_delta = "500ms"
timeout_prevote = "1s"
timeout_prevote_delta = "500ms"
timeout_precommit = "1s"
timeout_precommit_delta = "500ms"
timeout_commit = "5s"
```

### consensus.skip_timeout_commit

Make progress as soon as we have all the precommits (as if TimeoutCommit = 0)

Example

```toml
skip_timeout_commit = false
```

### consensus.create_empty_blocks
### consensus.create_empty_blocks_interval

EmptyBlocks mode and possible interval between empty blocks

Example

```toml
create_empty_blocks = true
create_empty_blocks_interval = "0s"
```

### consensus.peer_gossip_sleep_duration
### consensus.peer_query_maj23_sleep_duration

Reactor sleep duration parameters

Example

```toml
peer_gossip_sleep_duration = "100ms"
peer_query_maj23_sleep_duration = "2s"
```

## transactions indexer configuration options

Configuration options in the `[tx_index]` section of the `config.toml` file.

### tx_index.indexer

What indexer to use for transactions
Options:
 1) "null"
 2) "kv" (default) - the simplest possible indexer, backed by key-value storage (defaults to levelDB; see DBBackend).

Example

```toml
indexer = "kv"
```

### tx_indexer.index_keys

Comma-separated list of `compositeKeys` to index (by default the only key is "`tx.hash`")

Remember that Event has the following structure: `type.key`

```
type: [
  key: value,
  ...
]
```

You can also index transactions by height by adding "`tx.height`" key here.

It's recommended to index only a subset of keys due to possible memory
bloat. This is, of course, depends on the indexer's DB and the volume of
transactions.

Example

```toml
index_keys = ""
```

### tx_indexer.index_all_keys

When set to true, tells indexer to index all `compositeKeys` (predefined keys: "`tx.hash`", "`tx.height`" and all keys from `DeliverTx` responses).

:::tip Note
this may be not desirable (see the comment above). `IndexKeys` has a
precedence over `IndexAllKeys` (i.e. when given both, `IndexKeys` will be
indexed).
:::

Example

```toml
index_all_keys = true
```

## instrumentation configuration options

Configuration options in the `[instrumentation]` section of the `config.toml` file.

### instrumentation.prometheus

When true, Prometheus metrics are served under `/metrics` on `PrometheusListenAddr`.
Check out the documentation for the list of available metrics.

Example

```toml
prometheus = false
```

### instrumentation.prometheus_listen_addr

Address to listen for Prometheus collector(s) connections

Example

```toml
prometheus_listen_addr = ":26660"
```

### instrumentation.max_open_connections

Maximum number of simultaneous connections.

If you want to accept a larger number than the default, make sure you increase your OS limits.

0 - unlimited.

Example

```toml
max_open_connections = 3
```

### instumentation.namespace

Instrumentation namespace

Example

```toml
namespace = "tendermint"
```
