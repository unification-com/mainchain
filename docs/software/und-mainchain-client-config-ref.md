# `.und_mainchain/config/client.toml` Reference

The `$HOME/.und_mainchain/config/client.toml` file contains all the configuration options for the `und` client
configuration. Below is a reference for the file.

#### Contents

[[toc]]

## Main base config options

### chain-id

The network chain ID

Example

```toml
chain-id = "FUND-TestNet-2"
```

### keyring-backend

The keyring's backend, where the keys are stored (`os|file|kwallet|pass|test|memory`)

Example

```toml
keyring-backend = "os"
```

### output

CLI output format (`text|json`)

Example

```toml
output = "text"
```

### node

`<host>:<port>` to Tendermint RPC interface for this chain

Example

```toml
node = "https://rpc-testnet.unification.io:443"
```

### broadcast-mode

Transaction broadcasting mode (`sync|async|block`)

Example

```toml
broadcast-mode = "sync"
```
