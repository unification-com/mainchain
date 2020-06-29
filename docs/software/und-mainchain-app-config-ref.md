# `.und_mainchain/config/app.toml` Reference

The `$HOME/.und_mainchain/config/app.toml` file contains all the configuration options for the `und` server binary. Below is a reference for the file.

#### Contents

[[toc]]

## Main base config options

### minimum-gas-prices

The minimum gas prices a validator is willing to accept for processing a transaction. A transaction's fees must meet the minimum of any denomination specified in this config (e.g. `0.25nund`).

Example

```toml
minimum-gas-prices = "0.25nund"
```

### halt-height

HaltHeight contains a non-zero block height at which a node will gracefully halt and shutdown that can be used to assist upgrades and testing.

::: tip Note
Commitment of state will be attempted on the corresponding block.
:::

Example

```toml
halt-height = 12345
```

### halt-time

HaltTime contains a non-zero minimum block time (in Unix seconds) at which a node will gracefully halt and shutdown that can be used to assist upgrades and testing.

::: tip Note
Commitment of state will be attempted on the corresponding block.
:::

Example

```toml
halt-time = 1587565852
```

### inter-block-cache

InterBlockCache enables inter-block caching.

Example

```toml
inter-block-cache = true
```

### pruning

Pruning sets the pruning strategy: `syncable`, `nothing`, `everything`

::: danger IMPORTANT
There is a known issue with the `syncable` pruning option in the Cosmos SDK. Since `pruning = "syncable"` is the default value when `und init` is run, it is recommended to set the value to either `pruning = "everything"` or `pruning = "nothing"`. Note that setting to `pruning = "nothing"` will increase storage usage considerably.
:::

- `syncable`: only those states not needed for state syncing will be deleted (keeps last 100 + every 10000th)
- `nothing`: all historic states will be saved, nothing will be deleted (i.e. archiving node)
- `everything`: all saved states will be deleted, storing only the current state

Example

```toml
pruning = "syncable"
```
