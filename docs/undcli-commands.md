# `undcli` Command Reference

The `undcli` binary is the primary CLI client tool used for interacting with a full `und` node. The `und` node can be running locally, or one being run as a public service. By default, `undcli` will assume the `und` node is running locally and attempt to connect via RPC to `tcp://localhost:26657`.

#### Contents

[[toc]]

## undcli

Unification Mainchain CLI for interacting with Mainchain

Usage:
```bash
  undcli [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[config](#undcli-config)|Create or query an application CLI configuration file|
|[convert](#undcli-convert)|convert between nund<->UND denominations|
|[keys](#undcli-keys)|Add or view local private keys|
|[query](#undcli-query)|Querying subcommands|
|[rest-server](#undcli-rest-server)|Start LCD (light-client daemon), a local REST server|
|[status](#undcli-status)|Query remote node for status|
|[tx](#undcli-tx)|Transactions subcommands|
|[version](#undcli-version)|Print the app version|

**Global Flags**

::: tip
Global flags can be passed to any of the commands and sub-commands below
:::

| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of UND Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Use "`undcli [command] --help`" for more information about a command.

## undcli config

Create or query an application CLI configuration file

::: tip Note
`undcli` configuration is stored in `$HOME/.und_cli/config/config.toml`
:::

Usage:
```bash
  undcli config <key> [value] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--get`||print configuration value or its default if unset|
|`-h`, `--help`||help for config|

## undcli convert

convert between UND denominations

Usage:
```bash
  undcli convert [amount] [from_denom] [to_denom] [flags]
```

Example:
```bash
$ undcli convert 24 und nund
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for version|

## undcli keys

Keys allows you to manage your local keystore for tendermint.

These keys may be in any format supported by go-crypto and can be
used by light-clients, full nodes, or any other application that
needs to sign with a private key.

Usage:
```bash
  undcli keys [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[mnemonic](#undcli-keys-mnemonic)|Compute the bip39 mnemonic for some input entropy|
|[add](#undcli-keys-add)|Add an encrypted private key (either newly generated or recovered), encrypt it, and save to disk|
|[export](#undcli-keys-export)|Export private keys|
|[import](#undcli-keys-import)|Import private keys into the local keybase|
|[list](#undcli-keys-list)|List all keys|
|[show](#undcli-keys-show)|Show key info for the given name|
|[delete](#undcli-keys-delete)|Delete the given keys|
|[parse](#undcli-keys-parse)|Parse address from hex to bech32 and vice versa|
|[migrate](#undcli-keys-migrate)|Migrate keys from the legacy (db-based) Keybase|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for keys|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|

::: warning Note
Any `--ledger` and Ledger related flags below are currently unused by `undcli`. Ledger support for UND will be available in a future version.
:::

## undcli keys mnemonic

Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. To pass your own entropy, use --unsafe-entropy

Usage:
```bash
  undcli keys mnemonic [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for mnemonic|
|`--unsafe-entropy`||Prompt the user to supply their own entropy, instead of relying on the system|

## undcli keys add

Derive a new private key and encrypt to disk.
Optionally specify a BIP39 mnemonic, a BIP39 passphrase to further secure the mnemonic, and a bip32 HD path to derive a specific account. The key will be stored under the given name and encrypted with the given password. The only input that is required is the encryption password.

If run with `-i`, it will prompt the user for BIP44 path, BIP39 mnemonic, and passphrase.

The flag `--recover` allows one to recover a key from a seed passphrase.

If run with `--dry-run`, a key would be generated (or recovered) but not stored to the local keystore.

Use the `--pubkey` flag to add arbitrary public keys to the keystore for constructing multisig transactions.

You can add a multisig key by passing the list of key names you want the public
key to be composed of to the `--multisig` flag and the minimum number of signatures required through `--multisig-threshold`.

The keys are sorted by address, unless the flag `--nosort` is set.

Usage:
```bash
  undcli keys add <name> [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--account`|`uint32`|Account number for HD derivation|
|`--algo`|`string`|Key signing algorithm to generate keys for (default "`secp256k1`")|
|`--dry-run`||Perform action, but don't add key to local keystore|
|`--hd-path`|`string`|Manual HD Path derivation (overrides BIP44 config)|
|`-h`, `--help`||help for add|
|`--indent`||Add indent to JSON response|
|`--index`|`uint32`|Address index number for HD derivation|
|`-i`, `--interactive`||Interactively prompt user for BIP39 passphrase and mnemonic|
|`--ledger`||Store a local reference to a private key on a Ledger device|
|`--multisig`|`strings`|Construct and store a multisig public key (implies `--pubkey`)|
|`--multisig-threshold`|`uint`|K out of N required signatures. For use in conjunction with `--multisig` (default 1)|
|`--no-backup`||Don't print out seed phrase (if others are watching the terminal)|
|`--nosort`||Keys passed to `--multisig` are taken in the order they're supplied|
|`--pubkey`|`string`|Parse a public key in bech32 format and save it to disk|
|`--recover`||Provide seed phrase to recover existing key instead of creating|

## undcli keys export

Export a private key from the local keybase in ASCII-armored encrypted format.

Usage:
```bash
  undcli keys export <name> [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for export


## undcli keys import

Import a ASCII armored private key into the local keybase.

Usage:
```bash
  undcli keys import <name> <keyfile> [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for import|


## undcli keys list

Return a list of all public keys stored by this key manager
along with their associated name and address.

Usage:
```bash
  undcli keys list [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for list|
|`--indent`||Add indent to JSON response|
|`-n`, `--list-names`||List names only|


## undcli keys show

Return public details of a single local key. If multiple names are
provided, then an ephemeral multisig key will be created under the name "multi"
consisting of all the keys provided by name and multisig threshold.

Usage:
```bash
  undcli keys show [name] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-a`, `--address`||Output the address only (overrides `--output`)|
|`--bech`|`string`|The Bech32 prefix encoding for a key (`acc`\|`val`\|`cons`) (default "`acc`")|
|`-d`, `--device`||Output the address in a ledger device|
|`-h`, `--help`||help for show|
|`--indent`||Add indent to JSON response|
|`--multisig-threshold`|`uint`|K out of N required signatures (default 1)|
|`-p`, `--pubkey`||Output the public key only (overrides `--output`)|

## undcli keys delete

Delete keys from the Keybase backend.

Note that removing offline or ledger keys will remove
only the public key references stored locally, i.e.
private keys stored in a ledger device cannot be deleted with the CLI.

Usage:
```bash
  undcli keys delete <name>... [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-f`, `--force`||Remove the key unconditionally without asking for the passphrase. Deprecated.|
|`-h`, `--help`||help for delete|
|`-y`, `--yes`||Skip confirmation prompt when deleting offline or ledger key references|

## undcli keys parse

Convert and print to stdout key addresses and fingerprints from
hexadecimal into bech32 und prefixed format and vice versa.

Usage:
```bash
  undcli keys parse <hex-or-bech32-address> [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for parse|
|`--indent`||Indent JSON output|

## undcli keys migrate

Migrate key information from the legacy (db-based) Keybase to the new keyring-based Keybase.
For each key material entry, the command will prompt if the key should be skipped or not. If the key
is not to be skipped, the passphrase must be entered. The key will only be migrated if the passphrase
is correct. Otherwise, the command will exit and migration must be repeated.

It is recommended to run in 'dry-run' mode first to verify all key migration material.

Usage:
```bash
  undcli keys migrate [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--dry-run`||Run migration without actually persisting any changes to the new Keybase|
|`-h`, `--help`||help for migrate|

## undcli query

Querying subcommands. Each module has its own query sub-commands to query data on the chain.

Usage:
```bash
  undcli query [command]
```

Aliases:
  `query`, `q`

Available Commands:
| Command | Description |
|---------|-------------|
|[account](#undcli-query-account)|Query account information|
|[auth](#undcli-query-auth)|Querying commands for the auth module|
|[beacon](#undcli-query-beacon)|Querying commands for the beacon module|
|[block](#undcli-query-block)|Get verified data for a the block at given height|
|[distribution](#undcli-query-distribution)|Querying commands for the distribution module|
|[enterprise](#undcli-query-enterprise)|Querying commands for the enterprise module|
|[evidence](#undcli-query-evidence)|Query for evidence by hash or for all (paginated) submitted evidence|
|[gov](#undcli-query-gov)|Querying commands for the governance module|
|[slashing](#undcli-query-slashing)|Querying commands for the slashing module|
|[staking](#undcli-query-staking)|Querying commands for the staking module|
|[supply](#undcli-query-supply)|Query total supply including locked enterprise UND|
|[tendermint-validator-set](#undcli-query-tendermint-validator-set)|Get the full tendermint validator set at given height|
|[tx](#undcli-query-tx)|Query for a transaction by hash in a committed block|
|[txs](#undcli-query-txs)|Query for paginated transactions that match a set of events|
|[wrkchain](#undcli-query-wrkchain)|Querying commands for the wrkchain module|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for query|

Use "`undcli query [command] --help`" for more information about a command.

## undcli query account

Query account information

Usage:
```bash
  undcli query account [address] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for account|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query auth

Querying commands for the auth module

Usage:
```bash
  undcli query auth [flags]
  undcli query auth [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[account](#undcli-query-auth-account)|Query account balance|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for auth|

## undcli query auth account

Query account balance

Usage:
```bash
  undcli query auth account [address] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for account|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query beacon

Querying commands for the beacon module

Usage:
```bash
  undcli query beacon [flags]
  undcli query beacon [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[params](#undcli-query-beacon-params)|Query the current Beacon parameters|
|[beacon](#undcli-query-beacon-beacon)|Query a BEACON for given ID|
|[timestamp](#undcli-query-beacon-timestamp)|Query a BEACON for given ID and timestamp ID to retrieve recorded timestamp|
|[search](#undcli-query-beacon-search)|Query all BEACONs with optional filters|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for beacon|

## undcli query beacon params

Query the current Beacon parameters

Usage:
```bash
  undcli query beacon params [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for params|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query beacon beacon

Query a BEACON for given ID

Usage:
```bash
  undcli query beacon beacon [beacon id] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for beacon|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query beacon timestamp

Query a BEACON for given ID and timestamp ID to retrieve recorded timestamp

Usage:
```bash
  undcli query beacon timestamp [beacon id] [timestamp id] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for timestamp|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query beacon search

Query for all paginated BEACONs that match optional filters:

Usage:
```bash
  undcli query beacon search [flags]
```

Example:
```bash
$ undcli query beacon search --moniker beacon1
$ undcli query beacon search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ undcli query beacon search --page=2 --limit=100
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for search|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--limit`|`int`|pagination limit of beacons to query for (default 100)|
|`--moniker`|`string`|(optional) filter beacons by name|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--owner`|`string`|(optional) filter beacons by owner address|
|`--page`|`int`|pagination page of beacons to to query for (default 1)|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query block

Get verified data for a the block at given height

Usage:
```bash
  undcli query block [height] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for block|
|`-n`, `--node`|`string`|Node to connect to (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query distribution

Querying commands for the distribution module

Usage:
```bash
  undcli query distribution [flags]
  undcli query distribution [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[params](#undcli-query-distribution-params)|Query distribution params|
|[validator-outstanding-rewards](#undcli-query-distribution-validator-outstanding-rewards)|Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations|
|[commission](#undcli-query-distribution-commission)|Query distribution validator commission|
|[slashes](#undcli-query-distribution-slashes)|Query distribution validator slashes|
|[rewards](#undcli-query-distribution-rewards)|Query all distribution delegator rewards or rewards from a particular validator|
|[community-pool](#undcli-query-distribution-community-pool)|Query the amount of coins in the community pool|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for distribution|

## undcli query distribution params

Query distribution params

Usage:
```bash
  undcli query distribution params [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for params|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query distribution validator-outstanding-rewards

Query distribution outstanding (un-withdrawn) rewards
for a validator and all their delegations.

Usage:
```bash
  undcli query distribution validator-outstanding-rewards [validator] [flags]
```

Example:
```bash
  undcli query distribution validator-outstanding-rewards undvaloper1lwjmdnks33xwnmfayc64ycprww49n33mtm92ne
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for validator-outstanding-rewards|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query distribution commission

Query validator commission rewards from delegators to that validator.

Usage:
```bash
  undcli query distribution commission [validator] [flags]
```

Example:
```bash
  undcli query distribution commission undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for commission|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query distribution slashes

Query all slashes of a validator for a given block range.

Usage:
```bash
  undcli query distribution slashes [validator] [start-height] [end-height] [flags]
```

Example:
```bash
  undcli query distribution slashes undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 0 100
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for slashes|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query distribution rewards

Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

Usage:
```bash
  undcli query distribution rewards [delegator-addr] [<validator-addr>] [flags]
```
Example:
```bash
  $ undcli query distribution rewards und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
  $ undcli query distribution rewards und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for rewards|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query distribution community-pool

Query all coins in the community pool which is under Governance control.

Usage:
```bash
  undcli query distribution community-pool [flags]
```

Example:
```bash
  undcli query distribution community-pool
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for community-pool|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query enterprise

Querying commands for the enterprise module

Usage:
```bash
  undcli query enterprise [flags]
  undcli query enterprise [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[params](#undcli-query-enterprise-params)|Query the current enterprise UND parameters|
|[orders](#undcli-query-enterprise-orders)|Query Enterprise UND purchase orders with optional filters|
|[order](#undcli-query-enterprise-order)|get a purchase order by ID|
|[locked](#undcli-query-enterprise-locked)|get locked UND for an address|
|[total-locked](#undcli-query-enterprise-total-locked)|Query the current total locked enterprise UND|
|[total-unlocked](#undcli-query-enterprise-total-unlocked)|Query the current total unlocked und in circulation|
|[whitelist](#undcli-query-enterprise-whitelist)|get addresses whitelisted for raising enterprise purchase orders|
|[whitelisted](#undcli-query-enterprise-whitelisted)|check if given address is whitelested for purchase orders|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for enterprise|

## undcli query enterprise params

Query the current enterprise UND parameters

Usage:
```bash
  undcli query enterprise params [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query enterprise orders

Query for a all paginated Enterprise UND purchase orders that match optional filters:

Usage:
```bash
  undcli query enterprise orders [flags]
```

Example:
```bash
  $ undcli query enterprise orders --status (raised|accept|reject|complete)
  $ undcli query enterprise orders --purchaser und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
  $ undcli query enterprise orders --page=2 --limit=100
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

|`--limit`|`int`|pagination limit to query for (default 100)|
|`--page`|`int`|pagination page to query for|
|`--purchaser`|`string`|(optional) filter purchase orders raised by address|
|`--status`|`string`|(optional) filter purchase orders by status, status: `raised`/`accept`/`reject`/`complete`|

## undcli query enterprise order

get a purchase order by ID

Usage:
```bash
  undcli query enterprise order [purchase_order_id] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query enterprise locked

get locked UND for an address

Usage:
```bash
  undcli query enterprise locked [address] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query enterprise total-locked

Query the current total locked enterprise UND

Usage:
```bash
  undcli query enterprise total-locked [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query enterprise total-unlocked

Query the current total unlocked und in circulation

Usage:
```bash
  undcli query enterprise total-unlocked [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query enterprise whitelist

get addresses whitelisted for raising enterprise purchase orders

Usage:
```bash
  undcli query enterprise whitelist [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query enterprise whitelisted

check if given address is whitelested for purchase orders

Usage:
```bash
  undcli query enterprise whitelisted [address] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query evidence

Query for specific submitted evidence by hash or query for all (paginated) evidence.

Usage:
```bash
  undcli query evidence [flags]
  undcli query evidence [command]
```

Example:
```bash
  $ undcli query evidence DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660
  $ undcli query evidence --page=2 --limit=50
```

Available Commands:
| Command | Description |
|---------|-------------|
|[params](#undcli-query-evidence-params)|Query the current evidence parameters|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for evidence|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--limit`|`int`|pagination limit of evidence to query for (default 100)|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--page`|`int`|pagination page of evidence to to query for (default 1)|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query evidence params

Query the current evidence parameters:

Usage:
```bash
  undcli query evidence params [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for params|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query gov

Querying commands for the governance module

Usage:
```bash
  undcli query gov [flags]
  undcli query gov [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[proposal](#undcli-query-gov-proposal)|Query details of a single proposal|
|[proposals](#undcli-query-gov-proposals)|Query proposals with optional filters|
|[vote](#undcli-query-gov-vote)|Query details of a single vote|
|[votes](#undcli-query-gov-votes)|Query votes on a proposal|
|[param](#undcli-query-gov-param)|Query the parameters (voting|tallying|deposit) of the governance process|
|[params](#undcli-query-gov-params)|Query the parameters of the governance process|
|[proposer](#undcli-query-gov-proposer)|Query the proposer of a governance proposal|
|[deposit](#undcli-query-gov-deposit)|Query details of a deposit|
|[deposits](#undcli-query-gov-deposits)|Query deposits on a proposal|
|[tally](#undcli-query-gov-tally)|Get the tally of a proposal vote|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for gov|

## undcli query gov proposal

Query details for a proposal. You can find the
proposal-id by running "undcli query gov proposals".

Usage:
```bash
  undcli query gov proposal [proposal-id] [flags]
```

Example:
```bash
  undcli query gov proposal 1
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for proposal|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query gov proposals

Query for a all paginated proposals that match optional filters:

Usage:
```bash
  undcli query gov proposals [flags]
```

Example:
```bash
  $ undcli query gov proposals --depositor und1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
  $ undcli query gov proposals --voter und1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
  $ undcli query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected)
  $ undcli query gov proposals --page=2 --limit=100
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|
|`--limit`|`int`|pagination limit to query for (default 100)|
|`--page`|`int`|pagination page to query for|
|`--depositor`|`string`|(optional) filter by proposals deposited on by depositor
 (default 1)|
|`--status`|`string`|(optional) filter proposals by proposal status, status: `deposit_period`/`voting_period`/`passed`/`rejected`|
|`--voter`|`string`|(optional) filter by proposals voted on by voted|

## undcli query gov vote

Query details for a single vote on a proposal given its identifier.

Usage:
```bash
  undcli query gov vote [proposal-id] [voter-addr] [flags]
```

Example:
```bash
  undcli query gov vote 1 cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for vote|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query gov votes

Query vote details for a single proposal by its identifier.

Usage:
```bash
  undcli query gov votes [proposal-id] [flags]
```

Example:
```bash
  $ undcli query gov votes 1
  $ undcli query gov votes 1 --page=2 --limit=100
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|
|`--limit`|`int`|pagination limit to query for (default 100)|
|`--page`|`int`|pagination page to query for|

## undcli query gov param

Query the all the parameters for the governance process.

Usage:
```bash
  undcli query gov param [param-type] [flags]
```

Example:
```bash
  $ undcli query gov param voting
  $ undcli query gov param tallying
  $ undcli query gov param deposit
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for param|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query gov params

Query the all the parameters for the governance process.

Usage:
```bash
  undcli query gov params [flags]
```

Example:
```bash
  undcli query gov params
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for params|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query gov proposer

Query which address proposed a proposal with a given ID.

Usage:
```bash
  undcli query gov proposer [proposal-id] [flags]
```

Example:
```bash
  undcli query gov proposer 1
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for proposer|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query gov deposit

Query details for a single proposal deposit on a proposal by its identifier.

Usage:
```bash
  undcli query gov deposit [proposal-id] [depositer-addr] [flags]
```

Example:
```bash
  undcli query gov deposit 1 und1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for deposit|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query gov deposits

Query details for all deposits on a proposal.
You can find the proposal-id by running "undcli query gov proposals".

Usage:
```bash
  undcli query gov deposits [proposal-id] [flags]
```

Example:
```bash
  undcli query gov deposits 1
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for deposits|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query gov tally

Query tally of votes on a proposal. You can find
the proposal-id by running "undcli query gov proposals".

Usage:
```bash
  undcli query gov tally [proposal-id] [flags]
```

Example:
```bash
  undcli query gov tally 1
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for tally|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query slashing

Querying commands for the slashing module

Usage:
```bash
  undcli query slashing [flags]
  undcli query slashing [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[signing-info](#undcli-query-slashing-signing-info)|Query a validator's signing information|
|[params](#undcli-query-slashing-params)|Query the current slashing parameters|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for slashing|

## undcli query slashing signing-info

Use a validators' consensus public key to find the signing-info for that validator:

Usage:
```bash
  undcli query slashing signing-info [validator-conspub] [flags]
```

Example:
```bash
  undcli query slashing signing-info undvalconspub1zcjduepqfhvwcmt7p06fvdgexxhmz0l8c7sgswl7ulv7aulk364x4g5xsw7sr0k2g5
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for signing-info|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query slashing params

Query parameters for the slashing module:

Usage:
```bash
  undcli query slashing params [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for params|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking

Querying commands for the staking module

Usage:
```bash
  undcli query staking [flags]
  undcli query staking [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[delegation](#undcli-query-staking-delegation)|Query a delegation based on address and validator address|
|[delegations](#undcli-query-staking-delegations)|Query all delegations made by one delegator|
|[unbonding-delegation](#undcli-query-staking-unbonding-delegation)|Query an unbonding-delegation record based on delegator and validator address|
|[unbonding-delegations](#undcli-query-staking-unbonding-delegations)|Query all unbonding-delegations records for one delegator|
|[redelegation](#undcli-query-staking-redelegation)|Query a redelegation record based on delegator and a source and destination validator address|
|[redelegations](#undcli-query-staking-redelegations)|Query all redelegations records for one delegator|
|[validator](#undcli-query-staking-validator)|Query a validator|
|[validators](#undcli-query-staking-validators)|Query for all validators|
|[delegations-to](#undcli-query-staking-delegations-to)|Query all delegations made to one validator|
|[unbonding-delegations-from](#undcli-query-staking-unbonding-delegations-from)|Query all unbonding delegatations from a validator|
|[redelegations-from](#undcli-query-staking-redelegations-from)|Query all outgoing redelegatations from a validator|
|[historical-info](#undcli-query-staking-historical-info)|Query historical info at given height|
|[params](#undcli-query-staking-params)|Query the current staking parameters information|
|[pool](#undcli-query-staking-pool)|Query the current staking pool values|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for staking|

## undcli query staking delegation

Query delegations for an individual delegator on an individual validator.

Usage:
```bash
  undcli query staking delegation [delegator-addr] [validator-addr] [flags]
```

Example:
```bash
  undcli query staking delegation und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for delegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking delegations

Query delegations for an individual delegator on all validators.

Usage:
```bash
  undcli query staking delegations [delegator-addr] [flags]
```

Example:
```bash
  undcli query staking delegations und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for delegations|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking unbonding-delegation

Query unbonding delegations for an individual delegator on an individual validator.

Usage:
```bash
  undcli query staking unbonding-delegation [delegator-addr] [validator-addr] [flags]
```

Example:
```bash
  undcli query staking unbonding-delegation und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for unbonding-delegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking unbonding-delegations

Query unbonding delegations for an individual delegator.

Usage:
```bash
  undcli query staking unbonding-delegations [delegator-addr] [flags]
```

Example:
```bash
  undcli query staking unbonding-delegation und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for unbonding-delegations|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking redelegation

Query a redelegation record for an individual delegator between a source and destination validator.

Usage:
```bash
  undcli query staking redelegation [delegator-addr] [src-validator-addr] [dst-validator-addr] [flags]
```

Example:
```bash
  undcli query staking redelegation und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p undvaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking redelegations

Query all redelegation records for an individual delegator.

Usage:
```bash
  undcli query staking redelegations [delegator-addr] [flags]
```

Example:
```bash
  undcli query staking redelegation und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegations|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking validator

Query details about an individual validator.

Usage:
```bash
  undcli query staking validator [validator-addr] [flags]
```

Example:
```bash
  undcli query staking validator undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for validator|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking validators

Query details about all validators on a network.

Usage:
```bash
  undcli query staking validators [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for validators|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking delegations-to

Query delegations on an individual validator.

Usage:
```bash
  undcli query staking delegations-to [validator-addr] [flags]
```

Example:
```bash
  undcli query staking delegations-to undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for delegations-to|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking unbonding-delegations-from

Query delegations that are unbonding _from_ a validator.

Usage:
```bash
  undcli query staking unbonding-delegations-from [validator-addr] [flags]
```

Example:
```bash
  undcli query staking unbonding-delegations-from undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for unbonding-delegations-from|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking redelegations-from

Query delegations that are redelegating _from_ a validator.

Usage:
```bash
  undcli query staking redelegations-from [validator-addr] [flags]
```

Example:
```bash
  undcli query staking redelegations-from undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegations-from|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking historical-info

Query historical info at given height.

Usage:
```bash
  undcli query staking historical-info [height] [flags]
```

Example:
```bash
  undcli query staking historical-info 5
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for historical-info|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking params

Query values set as staking parameters.

Usage:
```bash
  undcli query staking params [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for params|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query staking pool

Query values for amounts stored in the staking pool.

Usage:
```bash
  undcli query staking pool [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for pool|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query supply

Query total UND supply, including locked and unlocked

Returns three values:

locked:
total UND locked through Enterprise purchases.
This UND is only available to pay WRKChain/BEACON fees
and cannot be used for transfers or staking/delegation

amount:
Liquid UND in active circulation, which can be used for
transfers, staking etc. It is the
LOCKED amount subtracted from TOTAL_SUPPLY

total_supply:
The total amount of UND currently on the chain, including locked UND

Usage:
```bash
  undcli query supply [flags]
```
Example:
```bash
  undcli query supply
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for supply|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query tendermint-validator-set

Get the full tendermint validator set at given height

Usage:
```bash
  undcli query tendermint-validator-set [height] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for tendermint-validator-set|
|`--indent`||indent JSON response|
|`--limit`|`int`|Query number of results returned per page (default 100)|
|`-n`, `--node`|`string`|Node to connect to (default "tcp://localhost:26657")|
|`--page`|`int`|Query a specific page of paginated results|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query tx

Query for a transaction by hash in a committed block

Usage:
```bash
  undcli query tx [hash] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for tx|
|`-n`, `--node`|string|Node to connect to (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|


## undcli query txs

Search for transactions that match the exact given events where results are paginated.

Each event takes the form of '`{eventType}.{eventAttribute}={value}`'. Please refer to each module's documentation for the full set of events to query for. Each module documents its respective events under 'xx_events.md'.

Usage:
```bash
  undcli query txs [flags]
```

Example:
```bash
undcli query txs --events 'message.sender=und1hp2km26czxlvesn8nmwswdd90umvcm5gxwpk98&message.action=withdraw_delegator_reward' --page 1 --limit 30
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--events`|`string`|list of transaction events in the form of `{eventType}.{eventAttribute}={value}`|
|`-h`, `--help`||help for txs|
|`--limit`|`uint32`|Query number of transactions results per page returned (default 30)|
|`-n`, `--node`|`string`|Node to connect to (default "tcp://localhost:26657")|
|`--page`|`uint32`|Query a specific page of paginated results (default 1)|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query wrkchain

Querying commands for the wrkchain module

Usage:
```bash
  undcli query wrkchain [flags]
  undcli query wrkchain [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[params](#undcli-query-wrkchain-params)|Query the current WRKChain parameters|
|[wrkchain](#undcli-query-wrkchain-wrkchain)|Query a WRKChain for given ID|
|[search](#undcli-query-wrkchain-search)|Query all WRKChains with optional filters|
|[block](#undcli-query-wrkchain-block)|Query a WRKChain for given ID and block height to retrieve recorded hashes for that block|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for wrkchain|

## undcli query wrkchain params

Query the current WRKChain parameters

Usage:
```bash
  undcli query wrkchain params [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query wrkchain wrkchain

Query a WRKChain for given ID

Usage:
```bash
  undcli query wrkchain wrkchain [wrkchain id] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli query wrkchain search

Query for all paginated WRKChains that match optional filters:

Usage:
```bash
  undcli query wrkchain search [flags]
```

Example:
```bash
  $ undcli query wrkchain search --moniker wrkchain1
  $ undcli query wrkchain search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
  $ undcli query wrkchain search --page=2 --limit=100
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|
|`--limit`|`int`|pagination limit to query for (default 100)|
|`--page`|`int`|pagination page to query for|
|`--moniker`|`string`|(optional) filter wrkchains by moniker|
|`--owner`|`string`|(optional) filter wrkchains by owner address|

## undcli query wrkchain block

Query a WRKChain for given ID and block height to retrieve recorded hashes for that block

Usage:
```bash
  undcli query wrkchain block [wrkchain id] [height] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

## undcli rest-server

Start LCD (light-client daemon), a local REST server. The REST server will serve all endpoints made available by the API. See the [Swagger](https://github.com/unification-com/mainchain/blob/master/client/lcd/swagger-ui/swagger.yaml) definition for more information.

Usage:
```bash
  undcli rest-server [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|--height|int|Use a specific height to query state at (this can error if the node is pruning state)|
|-h, --help||help for rest-server|
|--indent||Add indent to JSON response|
|--laddr|string|The address for the server to listen on (default "tcp://localhost:1317")|
|--ledger||Use a connected Ledger device|
|--max-open|int|The number of maximum open connections (default 1000)|
|--node|string|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|--read-timeout|uint|The RPC read timeout (in seconds) (default 10)|
|--trust-node||Trust connected full node (don't verify proofs for responses)|
|--write-timeout|uint|The RPC write timeout (in seconds) (default 10)|

## undcli status

Query remote node for status

Usage:
```bash
  undcli status [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for status|
|`--indent`||Add indent to JSON response|
|`-n`, `--node`|`string`|Node to connect to (default "tcp://localhost:26657")|

## undcli tx

Transactions subcommands, for generating, signing and broadcasting Txs to the chain.

Usage:
```bash
  undcli tx [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[broadcast](#undcli-tx-broadcast)|Broadcast transactions generated offline|
|[encode](#undcli-tx-encode)|Encode transactions generated offline|
|[decode](#undcli-tx-decode)|Decode an amino-encoded transaction string.|
|[multisign](#undcli-tx-multisign)|Generate multisig signatures for transactions generated offline|
|[send](#undcli-tx-send)|Create and sign a send tx|
|[sign](#undcli-tx-sign)|Sign transactions generated offline|

Module Specific Sub-Commands:
| Command | Description |
|---------|-------------|
|[auth](#undcli-tx-auth)|Auth transaction subcommands|
|[bank](#undcli-tx-bank)|Bank transaction subcommands|
|[beacon](#undcli-tx-beacon)|Beacon transaction subcommands|
|[crisis](#undcli-tx-crisis)|Crisis transactions subcommands|
|[distribution](#undcli-tx-distribution)|Distribution transactions subcommands|
|[enterprise](#undcli-tx-enterprise)|Enterprise UND transaction subcommands|
|[evidence](#undcli-tx-evidence)|Evidence transaction subcommands|
|[gov](#undcli-tx-gov)|Governance transactions subcommands|
|[slashing](#undcli-tx-slashing)|Slashing transactions subcommands|
|[staking](#undcli-tx-staking)|Staking transaction subcommands|
|[wrkchain](#undcli-tx-wrkchain)|WRKChain transaction subcommands|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for tx|

## undcli tx broadcast

Broadcast transactions created with the --generate-only
flag and signed with the sign command. Read a transaction from `[file_path]` and
broadcast it to a node. If you supply a dash (-) argument in place of an input
filename, the command reads from standard input.

Usage:
```bash
  undcli tx broadcast [file_path] [flags]
```

Example:
```bash
  undcli tx broadcast ./mytxn.json
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-a`, `--account-number`|`uint`|The account number of the signing account (offline mode only)|
|`-b`, `--broadcast-mode`|`string`|Transaction broadcasting mode (`sync`\|`async`\|`block`) (default "`sync`")|
|`--dry-run`||ignore the `--gas` flag and perform a simulation of a transaction, but don't broadcast it|
|`--fees`|`string`|Fees to pay along with transaction; eg: `10000nund`|
|`--from`|`string`|Name or address of private key with which to sign|
|`--gas`|`string`|gas limit to set per-transaction; set to "auto" to calculate required gas automatically (default 200000) (default "200000")|
|`--gas-adjustment`|`float`|adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)|
|`--gas-prices`|`string`|Gas prices to determine the transaction fee (e.g. `0.25nund`)|
|`--generate-only`||Build an unsigned transaction and write it to `STDOUT` (when enabled, the local Keybase is not accessible and the node operates offline)|
|`-h`, `--help`||help for broadcast|
|`--indent`||Add indent to JSON response|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`--ledger`||Use a connected Ledger device|
|`--memo`|`string`|Memo to send along with transaction|
|`--node`|`string`|\<host\>:\<port\> to tendermint rpc interface for this chain (default "tcp://localhost:26657")|
|`-s`, `--sequence`|`uint`|The sequence number of the signing account (offline mode only)|
|`--trust-node`||Trust connected full node (don't verify proofs for responses) (default true)|
|`-y`, `--yes`||Skip tx broadcasting prompt confirmation|

## undcli tx encode


## undcli tx decode


## undcli tx multisign


## undcli tx send


## undcli tx sign

## undcli tx auth
## undcli tx bank
## undcli tx beacon
## undcli tx crisis
## undcli tx distribution
## undcli tx enterprise
## undcli tx evidence
## undcli tx gov
## undcli tx slashing
## undcli tx staking
## undcli tx wrkchain

## undcli version

Print the app version

Usage:
```bash
  undcli version [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for version|
|`--long`||Print long version information|
