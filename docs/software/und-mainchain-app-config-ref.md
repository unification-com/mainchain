# `.und_mainchain/config/app.toml` Reference

The `$HOME/.und_mainchain/config/app.toml` file contains all the configuration options for the `und` server binary. Below is a reference for the file.

#### Contents

[[toc]]

## Main base config options

### minimum-gas-prices

The minimum gas prices a validator is willing to accept for processing a
transaction. A transaction's fees must meet the minimum of any denomination
specified in this config (e.g. 25.0nund).

Example

```toml
minimum-gas-prices = "25.0nund"
```

### pruning

default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Example

```toml
pruning = "default"
```

#### pruning-keep-recent
#### pruning-keep-every
#### pruning-interval

These are applied if and only if the pruning strategy is custom.

Example

```toml
pruning-keep-recent = "200"
pruning-keep-every = "1000"
pruning-interval = "100"
```

### halt-height

HaltHeight contains a non-zero block height at which a node will gracefully
halt and shutdown that can be used to assist upgrades and testing.

Note: Commitment of state will be attempted on the corresponding block.

Example

```toml
halt-height = 123456789
```

### halt-time

HaltTime contains a non-zero minimum block time (in Unix seconds) at which
a node will gracefully halt and shutdown that can be used to assist upgrades
and testing.

Note: Commitment of state will be attempted on the corresponding block.

Example

```toml
halt-time = 1654686000
```

### min-retain-blocks

MinRetainBlocks defines the minimum block height offset from the current
block being committed, such that all blocks past this offset are pruned
from Tendermint. It is used as part of the process of determining the
ResponseCommit.RetainHeight value during ABCI Commit. A value of 0 indicates
that no blocks should be pruned.

This configuration value is only responsible for pruning Tendermint blocks.
It has no bearing on application state pruning which is determined by the
"pruning-*" configurations.

Note: Tendermint block pruning is dependant on this parameter in conunction
with the unbonding (safety threshold) period, state pruning and state sync
snapshot parameters to determine the correct minimum value of
ResponseCommit.RetainHeight.

Example

```toml
min-retain-blocks = 0
```

### inter-block-cache

InterBlockCache enables inter-block caching.

Example

```toml
inter-block-cache = true
```

### index-events

IndexEvents defines the set of events in the form {eventType}.{attributeKey},
which informs Tendermint what to index. If empty, all events will be indexed.

Example

```toml
index-events = ["message.sender", "message.recipient"]
```

## telemetry

Configuration options in the `[telemetry]` section of the `app.toml` file.

### telemetry.service-name

Prefixed with keys to separate services.

Example

```toml
service-name = ""
```

### telemetry.enabled

Enabled enables the application telemetry functionality. When enabled,
an in-memory sink is also enabled by default. Operators may also enabled
other sinks such as Prometheus.

Example

```toml
enabled = false
```

### telemetry.enable-hostname

Enable prefixing gauge values with hostname.

Example

```toml
enable-hostname = false
```

### telemetry.enable-hostname-label

Enable adding hostname to labels.

Example

```toml
enable-hostname-label = false
```

### telemetry.enable-service-label

Enable adding service to labels.

Example

```toml
enable-service-label = false
```

### telemetry.prometheus-retention-time

PrometheusRetentionTime, when positive, enables a Prometheus metrics sink.

Example

```toml
prometheus-retention-time = 0
```

### telemetry.global-labels

GlobalLabels defines a global set of name/value label tuples applied to all
metrics emitted using the wrapper functions defined in telemetry package.

Example

```toml
global-labels = [["chain_id", "FUND-MainNet-2"]]
```

## api

Configuration options in the `[api]` section of the `app.toml` file.

### api.enable

Enable defines if the API server should be enabled.

Example

```toml
enable = false
```

### api.swagger

Swagger defines if swagger documentation should automatically be registered.

Example

```toml
swagger = false
```

### api.address

Address defines the API server to listen on.

Example

```toml
address = "tcp://0.0.0.0:1317"
```

### api.max-open-connections

MaxOpenConnections defines the number of maximum open connections.

Example

```toml
max-open-connections = 1000
```

### api.rpc-read-timeout

RPCReadTimeout defines the Tendermint RPC read timeout (in seconds).

Example

```toml
rpc-read-timeout = 10
```

### api.rpc-write-timeout

RPCWriteTimeout defines the Tendermint RPC write timeout (in seconds).

Example

```toml
rpc-write-timeout = 0
```

### api.rpc-max-body-bytes

RPCMaxBodyBytes defines the Tendermint maximum response body (in bytes).

Example

```toml
rpc-max-body-bytes = 1000000
```

### api.enabled-unsafe-cors

EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk).

Example

```toml
enabled-unsafe-cors = false
```

## grpc

Configuration options in the `[grpc]` section of the `app.toml` file.

### grpc.enable

Enable defines if the gRPC server should be enabled.

Example

```toml
enable = true
```

### grpc.address

Address defines the gRPC server address to bind to.

Example

```toml
address = "0.0.0.0:9090"
```

## state-sync

Configuration options in the `[state-sync]` section of the `app.toml` file.

State sync snapshots allow other nodes to rapidly join the network without replaying historical
blocks, instead downloading and applying a snapshot of the application state at a given height.

### state-sync.snapshot-interval

snapshot-interval specifies the block interval at which local state sync snapshots are
taken (0 to disable). Must be a multiple of pruning-keep-every.

Example

```toml
snapshot-interval = 500
```

### state-sync.snapshot-keep-recent

snapshot-keep-recent specifies the number of recent snapshots to keep and serve (0 to keep all).

Example

```toml
snapshot-keep-recent = 3
```
