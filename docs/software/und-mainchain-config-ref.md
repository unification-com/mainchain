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

TCP or UNIX socket address of the ABCI application,
or the name of an ABCI application compiled in with the Tendermint binary

Example

```toml
proxy_app = "tcp://127.0.0.1:26658"
```

### moniker

A custom human readable name for this node

Example

```toml
moniker = "node-moniker"
```

### fast_sync

If this node is many blocks behind the tip of the chain, FastSync
allows them to catchup quickly by downloading blocks in parallel
and verifying their commits

Example

```toml
fast_sync = true
```

### db_backend

Database backend: goleveldb | cleveldb | boltdb | rocksdb | badgerdb
* goleveldb (github.com/syndtr/goleveldb - most popular implementation)
  - pure go
  - stable
* cleveldb (uses levigo wrapper)
  - fast
  - requires gcc
  - use cleveldb build tag (go build -tags cleveldb)
* boltdb (uses etcd's fork of bolt - github.com/etcd-io/bbolt)
  - EXPERIMENTAL
  - may be faster is some use-cases (random reads - indexer)
  - use boltdb build tag (go build -tags boltdb)
* rocksdb (uses github.com/tecbot/gorocksdb)
  - EXPERIMENTAL
  - requires gcc
  - use rocksdb build tag (go build -tags rocksdb)
* badgerdb (uses github.com/dgraph-io/badger)
  - EXPERIMENTAL
  - use badgerdb build tag (go build -tags badgerdb)

Example

```toml
db_backend = "goleveldb"
```

### db_dir

Database directory

Example

```toml
db_dir = "data"
```

### log_level

Output level for logging, including package level options

Example

```toml
log_level = "info"
```

### log_format

Output format: 'plain' (colored text) or 'json'

Example

```toml
log_format = "plain"
```

### genesis_file

Path to the JSON file containing the initial validator set and other meta data

Example

```toml
genesis_file = "config/genesis.json"
```

### priv_validator_key_file

Path to the JSON file containing the private key to use as a validator in the consensus protocol

Example

```toml
priv_validator_key_file = "config/priv_validator_key.json"
```

### priv_validator_state_file

Path to the JSON file containing the last sign state of a validator

Example

```toml
priv_validator_state_file = "data/priv_validator_state.json"
```

### priv_validator_laddr

TCP or UNIX socket address for Tendermint to listen on for
connections from an external PrivValidator process

Example

```toml
priv_validator_laddr = ""
```

### node_key_file

Path to the JSON file containing the private key to use for node authentication in the p2p protocol

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

### filter_peers

If true, query the ABCI app on connecting to a new peer
so the app can decide if we should keep the connection or not

Example

```toml
filter_peers = false
```

## rpc

RPC Server Configuration Options found in `[rpc]`

### rpc.laddr

TCP or UNIX socket address for the RPC server to listen on

Example

```toml
laddr = "tcp://127.0.0.1:26657"
```

### rpc.cors_allowed_origins

A list of origins a cross-domain request can be executed from
Default value '[]' disables cors support
Use '["*"]' to allow any origin

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
NOTE: This server only supports /broadcast_tx_commit

Example

```toml
grpc_laddr = ""
```

### rpc.grpc_max_open_connections

Maximum number of simultaneous connections.
Does not include RPC (HTTP&WebSocket) connections. See max_open_connections
If you want to accept a larger number than the default, make sure
you increase your OS limits.
0 - unlimited.
Should be < {ulimit -Sn} - {MaxNumInboundPeers} - {MaxNumOutboundPeers} - {N of wal, db and other open files}
1024 - 40 - 10 - 50 = 924 = ~900

Example

```toml
grpc_max_open_connections = 900
```

### rpc.unsafe

Activate unsafe RPC commands like /dial_seeds and /unsafe_flush_mempool

Example

```toml
unsafe = false
```

### rpc.max_open_connections

Maximum number of simultaneous connections (including WebSocket).
Does not include gRPC connections. See grpc_max_open_connections
If you want to accept a larger number than the default, make sure
you increase your OS limits.
0 - unlimited.
Should be < {ulimit -Sn} - {MaxNumInboundPeers} - {MaxNumOutboundPeers} - {N of wal, db and other open files}
1024 - 40 - 10 - 50 = 924 = ~900

Example

```toml
max_open_connections = 900
```

### rpc.max_subscription_clients

Maximum number of unique clientIDs that can /subscribe
If you're using /broadcast_tx_commit, set to the estimated maximum number
of broadcast_tx_commit calls per block.

Example

```toml
max_subscription_clients = 100
```

### rpc.max_subscriptions_per_client

Maximum number of unique queries a given client can /subscribe to
If you're using GRPC (or Local RPC client) and /broadcast_tx_commit, set to
the estimated # maximum number of broadcast_tx_commit calls per block.

Example

```toml
max_subscriptions_per_client = 5
```

### rpc.timeout_broadcast_tx_commit

How long to wait for a tx to be committed during /broadcast_tx_commit.
WARNING: Using a value larger than 10s will result in increasing the
global HTTP write timeout, which applies to all connections and endpoints.
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
Might be either absolute path or path related to Tendermint's config directory.
If the certificate is signed by a certificate authority,
the certFile should be the concatenation of the server's certificate, any intermediates,
and the CA's certificate.
NOTE: both tls_cert_file and tls_key_file must be present for Tendermint to create HTTPS server.
Otherwise, HTTP server is run.

Example

```toml
tls_cert_file = ""
```

### rpc.tls_key_file

The path to a file containing matching private key that is used to create the HTTPS server.
Might be either absolute path or path related to Tendermint's config directory.
NOTE: both tls-cert-file and tls-key-file must be present for Tendermint to create HTTPS server.
Otherwise, HTTP server is run.

Example

```toml
tls_key_file = ""
```

### rpc.pprof_laddr

pprof listen address (https://golang.org/pkg/net/http/pprof)

Example

```toml
pprof_laddr = "localhost:6060"
```


## p2p

P2P Configuration Options in `[p2p]` section

### p2p.laddr

Address to listen for incoming connections

Example

```toml
laddr = "tcp://0.0.0.0:26656"
```

### p2p.external_address

Address to advertise to peers for them to dial
If empty, will use the same port as the laddr,
and will introspect on the listener or use UPnP
to figure out the address. ip and port are required
example: 159.89.10.97:26656

Example

```toml
external_address = ""
```

### p2p.seeds

Comma separated list of seed nodes to connect to

Example

```toml
seeds = ""
```

### p2p.persistent_peers

Comma separated list of nodes to keep persistent connections to

Example

```toml
persistent_peers = ""
```

### p2p.upnp

UPNP port forwarding

Example

```toml
upnp = false
```

### p2p.addr_book_file

Path to address book

Example

```toml
addr_book_file = "config/addrbook.json"
```

### p2p.addr_book_strict

Set true for strict address routability rules
Set false for private or local networks

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
unconditional_peer_ids = ""
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

Seed mode, in which node constantly crawls the network and looks for
peers. If another node asks it for addresses, it responds and disconnects.

Does not work if the peer-exchange reactor is disabled.

Example

```toml
seed_mode = false
```

### p2p.private_peer_ids

Comma separated list of peer IDs to keep private (will not be gossiped to other peers)

Example

```toml
private_peer_ids = ""
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

## mempool

Mempool Configuration Option `[mempool]` section

### mempool.recheck
### mempool.broadcast
### mempool.wal_dir

Example

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
This only accounts for raw transactions (e.g. given 1MB transactions and
max_txs_bytes=5MB, mempool will only accept 5 transactions).

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

### mempool.keep-invalid-txs-in-cache

Do not remove invalid transactions from the cache (default: false)
Set to true if it's not possible for any invalid transaction to become valid
again in the future.

Example

```toml
keep-invalid-txs-in-cache = false
```

### mempool.max_tx_bytes

Maximum size of a single transaction.
NOTE: the max size of a tx transmitted over the network is {max_tx_bytes}.

Example

```toml
max_tx_bytes = 1048576
```

### mempool.max_batch_bytes

Maximum size of a batch of transactions to send to a peer
Including space needed by encoding (one varint per transaction).
XXX: Unused due to https://github.com/tendermint/tendermint/issues/5796

Example

```toml
max_batch_bytes = 0
```

## statesync

State Sync Configuration Options `[statesync]`

State sync rapidly bootstraps a new node by discovering, fetching, and restoring a state machine
snapshot from peers instead of fetching and replaying historical blocks. Requires some peers in
the network to take and serve state machine snapshots. State sync is not attempted if the node
has any local state (LastBlockHeight > 0). The node will have a truncated block history,
starting from the height of the snapshot.

### statesync.enable

Example

```toml
enable = false
```

### statesync.rpc_servers
### statesync.trust_height
### statesync.trust_hash
### statesync.trust_period

RPC servers (comma-separated) for light client verification of the synced state machine and
retrieval of state data for node bootstrapping. Also needs a trusted height and corresponding
header hash obtained from a trusted source, and a period during which validators can be trusted.

For Cosmos SDK-based chains, trust_period should usually be about 2/3 of the unbonding time (~2
weeks) during which they can be financially punished (slashed) for misbehavior.

Example

```toml
rpc_servers = ""
trust_height = 0
trust_hash = ""
trust_period = "168h0m0s"
```

### statesync.discovery_time

Time to spend discovering snapshots before initiating a restore.

Example

```toml
discovery_time = "15s"
```

### statesync.temp_dir

Temporary directory for state sync snapshot chunks, defaults to the OS tempdir (typically /tmp).
Will create a new, randomly named directory within, and remove it when done.

Example

```toml
temp_dir = ""
```

### statesync.chunk_request_timeout

The timeout duration before re-requesting a chunk, possibly from a different
peer (default: 1 minute).

Example

```toml
chunk_request_timeout = "30s"
```

### statesync.chunk_fetchers

The number of concurrent chunk fetchers to run (default: 1).

Example

```toml
chunk_fetchers = "4"
```

## fastsync

Fast Sync Configuration Connections in '[fastsync]'

### fastsync.version

Fast Sync version to use:
  1) "v0" (default) - the legacy fast sync implementation
  2) "v1" - refactor of v0 version for better testability
  2) "v2" - complete redesign of v0, optimized for testability & readability

Example

```toml
version = "v0"
```

## consensus

Consensus Configuration Options in `[consensus]`

### consensus.wal_file

Example

```toml
wal_file = "data/cs.wal/wal"
```

### consensus.timeout_propose

How long we wait for a proposal block before prevoting nil

Example

```toml
timeout_propose = "3s"
```

### consensus.timeout_propose_delta

How much timeout_propose increases with each round

Example

```toml
timeout_propose_delta = "500ms"
```

### consensus.timeout_prevote

How long we wait after receiving +2/3 prevotes for “anything” (ie. not a single block or nil)

Example

```toml
timeout_prevote = "1s"
```

### consensus.timeout_prevote_delta

How much the timeout_prevote increases with each round

Example

```toml
timeout_prevote_delta = "500ms"
```

### consensus.timeout_precommit

How long we wait after receiving +2/3 precommits for “anything” (ie. not a single block or nil)

Example

```toml
timeout_precommit = "1s"
```

### consensus.timeout_precommit_delta

How much the timeout_precommit increases with each round

Example

```toml
timeout_precommit_delta = "500ms"
```

### consensus.timeout_commit

How long we wait after committing a block, before starting on the new
height (this gives us a chance to receive some more precommits, even
though we already have +2/3).

Example

```toml
timeout_commit = "5s"
```

### consensus.double_sign_check_height

How many blocks to look back to check existence of the node's consensus votes before joining consensus
When non-zero, the node will panic upon restart
if the same consensus key was used to sign {double_sign_check_height} last blocks.
So, validators should stop the state machine, wait for some blocks, and then restart the state machine to avoid panic.

Example

```toml
double_sign_check_height = 0
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

## tx_index

Transaction Indexer Configuration Options in `[tx_index]`

### tx_index.indexer

What indexer to use for transactions

The application will set which txs to index. In some cases a node operator will be able
to decide which txs to index based on configuration set in the application.

Options:
  1) "null"
  2) "kv" (default) - the simplest possible indexer, backed by key-value storage (defaults to levelDB; see DBBackend).
		- When "kv" is chosen "tx.height" and "tx.hash" will always be indexed.

Example

```toml
indexer = "kv"
```

## instrumentation

Instrumentation Configuration Options in `[instrumentation]`

### instrumentation.prometheus

When true, Prometheus metrics are served under /metrics on
PrometheusListenAddr.
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
If you want to accept a larger number than the default, make sure
you increase your OS limits.
0 - unlimited.

Example

```toml
max_open_connections = 3
```

### instrumentation.namespace

Instrumentation namespace

Example

```toml
namespace = "tendermint"
```
