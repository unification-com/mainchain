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
|[convert](#undcli-convert)|convert between nund<->FUND denominations|
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
|`--chain-id`|`string`|Chain ID of und Mainchain node|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli convert

convert between FUND denominations

Usage:
```bash
  undcli convert [amount] [from_denom] [to_denom] [flags]
```

Example:
```bash
$ undcli convert 24 fund nund
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
Any `--ledger` and Ledger related flags below are currently unused by `undcli`. Ledger support for und will be available in a future version.
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli keys export

Export a private key from the local keybase in ASCII-armored encrypted format.

Usage:
```bash
  undcli keys export <name> [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for export|

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|


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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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

Global Flags:
| Flag | Type | Description |

|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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
|[supply](#undcli-query-supply)|Query total supply including locked enterprise FUND|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|


Example:
```bash
  undcli query account und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
```

`--output=json` Result:
```json
{
  "account": {
    "type": "cosmos-sdk/Account",
    "value": {
      "address": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs",
      "coins": [
        {
          "denom": "nund",
          "amount": "975269999995000"
        }
      ],
      "public_key": "undpub1addwnpepqvcner5ngj4tqxadx3wfpsruzcv533xnpjq9arqrkunw30sjy53tus776nr",
      "account_number": 36,
      "sequence": 34732
    }
  },
  "enterprise": {
    "locked": {
      "denom": "nund",
      "amount": "99999000000000"
    },
    "available_for_wrkchain": [
      {
        "denom": "nund",
        "amount": "1075268999995000"
      }
    ]
  }
}
```

`--output=text` Result:
```yaml
account: |
  address: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
  coins:
  - denom: nund
    amount: "975269999995000"
  public_key: undpub1addwnpepqvcner5ngj4tqxadx3wfpsruzcv533xnpjq9arqrkunw30sjy53tus776nr
  account_number: 36
  sequence: 34733
enterprise:
  locked:
    denom: nund
    amount: "99999000000000"
  available:
  - denom: nund
    amount: "1075268999995000"

```

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

Query account balance.

::: warning Note
The `undcli query auth account` command will only return the liquid FUND available for the account, and does not include any Enterprise UND information.

To obtain full account information, including Enterprise FUND, use `undcli query account`
:::

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
undcli query auth account und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
```

`--output=json` Result:
```json
{
  "type": "cosmos-sdk/Account",
  "value": {
    "address": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs",
    "coins": [
      {
        "denom": "nund",
        "amount": "975269999995000"
      }
    ],
    "public_key": "undpub1addwnpepqvcner5ngj4tqxadx3wfpsruzcv533xnpjq9arqrkunw30sjy53tus776nr",
    "account_number": 36,
    "sequence": 34736
  }
}
```

`--output=text` Result:
```yaml
|
  address: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
  coins:
  - denom: nund
    amount: "975269999995000"
  public_key: undpub1addwnpepqvcner5ngj4tqxadx3wfpsruzcv533xnpjq9arqrkunw30sjy53tus776nr
  account_number: 36
  sequence: 34740
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query beacon params
```

`--output=json` Result:
```json
{
  "fee_register": "1000000000",
  "fee_record": "1000000000",
  "denom": "nund"
}
```

`--output=text` Result:
```yaml
fee_register: 1000000000
fee_record: 1000000000
denom: nund
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query beacon beacon 1
```

`--output=json` Result:
```json
{
  "beacon_id": "1",
  "moniker": "MyBeacon",
  "name": "My BEACON",
  "last_timestamp_id": "4885",
  "owner": "und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly"
}

```

`--output=text` Result:
```yaml
beaconid: 1
moniker: MyBeacon
name: My BEACON
lasttimestampid: 4885
owner: und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query beacon timestamp 1 1
```

`--output=json` Result:
```json
{
  "beacon_id": "1",
  "timestamp_id": "1",
  "submit_time": "1",
  "hash": "123",
  "owner": "und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly"
}
```

`--output=text` Result:
```yaml
text
beaconid: 1
timestampid: 1
submittime: 1
hash: "123"
owner: und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly
```

## undcli query beacon search

Query for all paginated BEACONs that match optional filters:

Usage:
```bash
  undcli query beacon search [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
$ undcli query beacon search --moniker MyBeacon
$ undcli query beacon search --owner und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly
$ undcli query beacon search --page=1 --limit=100
```

`--output=json` Result:
```json
[
  {
    "beacon_id": "1",
    "moniker": "MyBeacon",
    "name": "My BEACON",
    "last_timestamp_id": "4885",
    "owner": "und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly"
  }
]
```

`--output=text` Result:
```yaml
- beaconid: 1
  moniker: MyBeacon
  name: My BEACON
  lasttimestampid: 4885
  owner: und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query block 1234
```

`--output=json` Result:
```json
{
  "block_id": {
    "hash": "A0544753D01D09963A24B9DC67E24AEA3F073F8F82E67985E82958BFDBF41C04",
    "parts": {
      "total": "1",
      "hash": "5549CDD5511EF7EC573981F69F4C511057DB69DE1A19BE43AB205A5956CC4270"
    }
  },
  "block": {
    "header": {
      "version": {
        "block": "10",
        "app": "0"
      },
      "chain_id": "UND-Mainchain-TestNet-v4",
      "height": "1234",
      "time": "2020-04-09T18:48:00.246726475Z",
      "last_block_id": {
        "hash": "066E23184AC1AE19D25C986FBE3422132EF03D934A77A010D33D00136591332A",
        "parts": {
          "total": "1",
          "hash": "9E5EE87EDBA06C026456C73F2286E1609860109E384A9A3E615C75463C183BE6"
        }
      },
      "last_commit_hash": "98CCEB6DEBD4E04C2E485218F699D0F77CAB2A40C447BCF72A1E07AA7BC45178",
      "data_hash": "",
      "validators_hash": "75AAC0A199E23518AAF2CB2BBAF0B03772289AA461A8E820B70ED0B213F55CA3",
      "next_validators_hash": "75AAC0A199E23518AAF2CB2BBAF0B03772289AA461A8E820B70ED0B213F55CA3",
      "consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
      "app_hash": "722B2650285AFE318554F6FBD4C2564349B14C1813A9B2AB2A15BC38ECC56751",
      "last_results_hash": "",
      "evidence_hash": "",
      "proposer_address": "9309670C2C2DC4020F06E9D4F5F184A0121B17B5"
    },
    "data": {
      "txs": null
    },
    "evidence": {
      "evidence": null
    },
    "last_commit": {
      "height": "1233",
      "round": "0",
      "block_id": {
        "hash": "066E23184AC1AE19D25C986FBE3422132EF03D934A77A010D33D00136591332A",
        "parts": {
          "total": "1",
          "hash": "9E5EE87EDBA06C026456C73F2286E1609860109E384A9A3E615C75463C183BE6"
        }
      },
      "signatures": [
        {
          "block_id_flag": 2,
          "validator_address": "176E25F7B0C100B63C9B3A2D886D93EF11FF3AE8",
          "timestamp": "2020-04-09T18:48:00.359553211Z",
          "signature": "oYbC675btbp7y8xh9UjB7GzHh6xdCsAJhpEmuKcCpHu4VJD1pbhcszuLpYWCDv/YF4DnEb0VSfZKkhAMLpcbDg=="
        },
        {
          "block_id_flag": 2,
          "validator_address": "6CF1CBC687CDFB037745BBBEED920EA34025FF52",
          "timestamp": "2020-04-09T18:48:00.167018827Z",
          "signature": "kJyFAvh7rP7J6zVaLoWsgJSKH/8tGHDuNArvTaSf7DH7vr7/yA0b7wJB2/FunURazuX0MCzOQgkXZsvaM4P6Ag=="
        },
        {
          "block_id_flag": 2,
          "validator_address": "771340B06E0FDBEF42C89316A7C9461C75E2636F",
          "timestamp": "2020-04-09T18:48:00.246726475Z",
          "signature": "GtW0eJ1YZNPycYw7Szl9GiEdXvXv+49T5ciVoS5pFhQJkLirCdQk+oGoHe6EiEgoVUsDRgX+1RHadRT9FXyFAA=="
        },
        {
          "block_id_flag": 2,
          "validator_address": "9309670C2C2DC4020F06E9D4F5F184A0121B17B5",
          "timestamp": "2020-04-09T18:48:00.175695327Z",
          "signature": "cGuYan8WqUF+UFdkeyfRYrMk/OSZ5Vhvh/JtyeR/VvrEt+ij/gzCF/8JPktsMUMCL59C1x+/7KEtlnXQBpvGDQ=="
        },
        {
          "block_id_flag": 2,
          "validator_address": "BCEE21D0C9BB7DC6454966670A140DE9D67D488E",
          "timestamp": "2020-04-09T18:47:59.559232798Z",
          "signature": "KdDpXrrOn+e5frneN8JB4XkOUDdPR0PJeqHZayV0gPZ528VSB9hwYf1LjJnYHR/7tH6EUsJ/1M9fky8kio8DAQ=="
        }
      ]
    }
  }
}
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query distribution params
```

`--output=json` Result:
```json
{
  "community_tax": "0.000000000000000000",
  "base_proposer_reward": "0.010000000000000000",
  "bonus_proposer_reward": "0.040000000000000000",
  "withdraw_addr_enabled": true
}
```

`--output=text` Result:
```yaml
community_tax: "0.000000000000000000"
base_proposer_reward: "0.010000000000000000"
bonus_proposer_reward: "0.040000000000000000"
withdraw_addr_enabled: true
```

## undcli query distribution validator-outstanding-rewards

Query distribution outstanding (un-withdrawn) rewards
for a validator and all their delegations.

Usage:
```bash
  undcli query distribution validator-outstanding-rewards [validator] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query distribution validator-outstanding-rewards undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
```

`--output=json` Result:
```json
[
  {
    "denom": "nund",
    "amount": "1238029579197.350778828167866215"
  }
]
```

`--output=text` Result:
```yaml
- denom: nund
  amount: "1238029579197.350778828167866215"
```

## undcli query distribution commission

Query validator commission rewards from delegators to that validator.

Usage:
```bash
  undcli query distribution commission [validator] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query distribution commission undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
```

`--output=json` Result:
```json
[
  {
    "denom": "nund",
    "amount": "173834351436.585261128618628279"
  }
]
```

`--output=text` Result:
```yaml
- denom: nund
  amount: "173834351436.585261128618628279"
```

## undcli query distribution slashes

Query all slashes of a validator for a given block range.

Usage:
```bash
  undcli query distribution slashes [validator] [start-height] [end-height] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query distribution slashes undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps 0 100000
```

## undcli query distribution rewards

Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

Usage:
```bash
  undcli query distribution rewards [delegator-addr] [<validator-addr>] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  $ undcli query distribution rewards und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
  $ undcli query distribution rewards und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
```

`--output=json` Result:
```json
{
  "rewards": [
    {
      "validator_address": "undvaloper1duyhqzcgrzjy9y2yvueur2h3e2yqxhjl6yvpze",
      "reward": [
        {
          "denom": "nund",
          "amount": "83489748666.746008599999999999"
        }
      ]
    },
    {
      "validator_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
      "reward": [
        {
          "denom": "nund",
          "amount": "76284149221.457684407698415710"
        }
      ]
    },
    {
      "validator_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
      "reward": [
        {
          "denom": "nund",
          "amount": "1065310354090.116372019242500000"
        }
      ]
    }
  ],
  "total": [
    {
      "denom": "nund",
      "amount": "1225084251978.320065026940915709"
    }
  ]
}
```

`--output=text` Result:
```yaml
- validator_address: undvaloper1duyhqzcgrzjy9y2yvueur2h3e2yqxhjl6yvpze
  reward:
  - denom: nund
    amount: "83489748666.746008599999999999"
- validator_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
  reward:
  - denom: nund
    amount: "76284149221.457684407698415710"
- validator_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
  reward:
  - denom: nund
    amount: "1065310354090.116372019242500000"
total:
- denom: nund
  amount: "1225084251978.320065026940915709"
```

## undcli query distribution community-pool

Query all coins in the community pool which is under Governance control.

Usage:
```bash
  undcli query distribution community-pool [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query distribution community-pool
```

`--output=json` Result:
```json
[
  {
    "denom": "nund",
    "amount": "25686887481640.752897315286145048"
  }
]
```

`--output=text` Result:
```yaml
- denom: nund
  amount: "25686887481640.752897315286145048"
```

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
|[params](#undcli-query-enterprise-params)|Query the current enterprise FUND parameters|
|[orders](#undcli-query-enterprise-orders)|Query Enterprise FUND purchase orders with optional filters|
|[order](#undcli-query-enterprise-order)|get a purchase order by ID|
|[locked](#undcli-query-enterprise-locked)|get locked FUND for an address|
|[total-locked](#undcli-query-enterprise-total-locked)|Query the current total locked enterprise UND|
|[total-unlocked](#undcli-query-enterprise-total-unlocked)|Query the current total unlocked und in circulation|
|[whitelist](#undcli-query-enterprise-whitelist)|get addresses whitelisted for raising enterprise purchase orders|
|[whitelisted](#undcli-query-enterprise-whitelisted)|check if given address is whitelested for purchase orders|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for enterprise|

## undcli query enterprise params

Query the current enterprise FUND parameters

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query enterprise params
```

`--output=json` Result:
```json
{
  "ent_signers": "und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn",
  "denom": "nund",
  "min_Accepts": "1",
  "decision_time_limit": "84600"
}
```

`--output=text` Result:
```yaml
ent_signers: und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn
denom: nund
min_Accepts: 1
decision_time_limit: 84600
```

## undcli query enterprise orders

Query for a all paginated Enterprise FUND purchase orders that match optional filters:

Usage:
```bash
  undcli query enterprise orders [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  $ undcli query enterprise orders --status (raised|accept|reject|complete)
  $ undcli query enterprise orders --purchaser und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
  $ undcli query enterprise orders --page=1 --limit=100
```

`--output=json` Result:
```json
[
  {
    "id": "7",
    "purchaser": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs",
    "amount": {
      "denom": "nund",
      "amount": "10000000000000"
    },
    "status": "complete",
    "raise_time": "1585559505",
    "decisions": [
      {
        "signer": "und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn",
        "decision": "accept",
        "decision_time": "1585559536"
      }
    ],
    "completion_time": "1585559541"
  },
  {
    "id": "8",
    "purchaser": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs",
    "amount": {
      "denom": "nund",
      "amount": "100000000000000"
    },
    "status": "complete",
    "raise_time": "1587117676",
    "decisions": [
      {
        "signer": "und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn",
        "decision": "accept",
        "decision_time": "1587117836"
      }
    ],
    "completion_time": "1587117841"
  }
]
```

`--output=text` Result:
```yaml
- purchaseorderid: 7
  purchaser: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
  amount:
    denom: nund
    amount: "10000000000000"
  status: 4
  raisedtime: 1585559505
  decisions:
  - signer: und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn
    decision: 2
    decisiontime: 1585559536
  completiontime: 1585559541
- purchaseorderid: 8
  purchaser: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
  amount:
    denom: nund
    amount: "100000000000000"
  status: 4
  raisedtime: 1587117676
  decisions:
  - signer: und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn
    decision: 2
    decisiontime: 1587117836
  completiontime: 1587117841
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query enterprise order 8
```

`--output=json` Result:
```json
{
  "id": "8",
  "purchaser": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs",
  "amount": {
    "denom": "nund",
    "amount": "100000000000000"
  },
  "status": "complete",
  "raise_time": "1587117676",
  "decisions": [
    {
      "signer": "und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn",
      "decision": "accept",
      "decision_time": "1587117836"
    }
  ],
  "completion_time": "1587117841"
}
```

`--output=text` Result:
```yaml
purchaseorderid: 8
purchaser: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
amount:
  denom: nund
  amount: "100000000000000"
status: 4
raisedtime: 1587117676
decisions:
- signer: und1dz0llvwrhg4ln6ngwl69v9nvvsvj6y52h96ndn
  decision: 2
  decisiontime: 1587117836
completiontime: 1587117841
```

## undcli query enterprise locked

get locked FUND for an address

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query enterprise locked und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
```

`--output=json` Result:
```json
{
  "owner": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs",
  "amount": {
    "denom": "nund",
    "amount": "99966000000000"
  }
}
```

`--output=text` Result:
```yaml
amount:
  denom: nund
  amount: "99965000000000"
```

## undcli query enterprise total-locked

Query the current total locked enterprise FUND

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query enterprise total-locked
```

`--output=json` Result:
```json
{
  "denom": "nund",
  "amount": "101520000000000"
}
```

`--output=text` Result:
```yaml
denom: nund
amount: "101519000000000"
```

## undcli query enterprise total-unlocked

Query the current total unlocked (liquid) FUND in circulation

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query enterprise total-unlocked
```

`--output=json` Result:
```json
{
  "denom": "nund",
  "amount": "1002737496556850989"
}
```

`--output=text` Result:
```yaml
denom: nund
amount: "1002737498556850989"
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query enterprise whitelist
```

`--output=json` Result:
```json
[
  "und1q43fg7x7yn6wv3zxwxjknj7xv7tfqd0ahnc0mv",
  "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs",
  "und15s4ec3s97tu4pstk8tq86l5ues4dxnmadqmrjl",
  "und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly"
]
```

`--output=text` Result:
```yaml
- und1q43fg7x7yn6wv3zxwxjknj7xv7tfqd0ahnc0mv
- und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
- und15s4ec3s97tu4pstk8tq86l5ues4dxnmadqmrjl
- und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query enterprise whitelisted und1mtxp3jh5ytygjfpfyx4ell495zc2m4k8ft8uly
```

`--output=json` Result:
```json
true
```

`--output=text` Result:
```yaml
true
```

## undcli query evidence

Query for specific submitted evidence by hash or query for all (paginated) evidence.

Usage:
```bash
  undcli query evidence [flags]
  undcli query evidence [command]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query evidence params
```

`--output=json` Result:
```json
{
  "max_evidence_age": "120000000000"
}
```

`--output=text` Result:
```yaml
max_evidence_age: 2m0s
```

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

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for proposal|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query gov proposal 9
```

`--output=json` Result:
```json
{
  "content": {
    "type": "cosmos-sdk/ParameterChangeProposal",
    "value": {
      "title": "Slashing parameters",
      "description": "change the signed blocks window from 100 to 10,000, and minimum signed requirement from 50% to 5%",
      "changes": [
        {
          "subspace": "slashing",
          "key": "SignedBlocksWindow",
          "value": "\"10000\""
        },
        {
          "subspace": "slashing",
          "key": "MinSignedPerWindow",
          "value": "\"0.050000000000000000\""
        }
      ]
    }
  },
  "id": "9",
  "proposal_status": "Passed",
  "final_tally_result": {
    "yes": "3714680003497661",
    "abstain": "0",
    "no": "0",
    "no_with_veto": "0"
  },
  "submit_time": "2020-04-09T09:58:59.281924052Z",
  "deposit_end_time": "2020-04-11T09:58:59.281924052Z",
  "total_deposit": [
    {
      "denom": "nund",
      "amount": "1000000000000"
    }
  ],
  "voting_start_time": "2020-04-09T09:58:59.281924052Z",
  "voting_end_time": "2020-04-11T09:58:59.281924052Z"
}
```

`--output=text` Result:
```yaml
content:
  title: Slashing parameters
  description: change the signed blocks window from 100 to 10,000, and minimum signed
    requirement from 50% to 5%
  changes:
  - subspace: slashing
    key: SignedBlocksWindow
    value: '"10000"'
  - subspace: slashing
    key: MinSignedPerWindow
    value: '"0.050000000000000000"'
id: 9
proposal_status: 3
final_tally_result:
  "yes": "3714680003497661"
  abstain: "0"
  "no": "0"
  no_with_veto: "0"
submit_time: 2020-04-09T09:58:59.281924052Z
deposit_end_time: 2020-04-11T09:58:59.281924052Z
total_deposit:
- denom: nund
  amount: "1000000000000"
voting_start_time: 2020-04-09T09:58:59.281924052Z
voting_end_time: 2020-04-11T09:58:59.281924052Z
```
## undcli query gov proposals

Query for a all paginated proposals that match optional filters:

Usage:
```bash
  undcli query gov proposals [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  $ undcli query gov proposals --depositor und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce
  $ undcli query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected)
  $ undcli query gov proposals --page=1 --limit=100
```

`--output=json` Result:
```json
[
  {
    "content": {
      "type": "cosmos-sdk/CommunityPoolSpendProposal",
      "value": {
        "title": "Community Pool Spend",
        "description": "Send community pool FUND to the TN Faucet",
        "recipient": "und17jv7rerc2e3undqumpf32a3xs9jc0kjk4z2car",
        "amount": [
          {
            "denom": "nund",
            "amount": "25000000000000"
          }
        ]
      }
    },
    "id": "7",
    "proposal_status": "Rejected",
    "final_tally_result": {
      "yes": "0",
      "abstain": "0",
      "no": "0",
      "no_with_veto": "0"
    },
    "submit_time": "2020-03-20T12:22:18.684400772Z",
    "deposit_end_time": "2020-03-22T12:22:18.684400772Z",
    "total_deposit": [
      {
        "denom": "nund",
        "amount": "1000000000000"
      }
    ],
    "voting_start_time": "2020-03-20T12:22:18.684400772Z",
    "voting_end_time": "2020-03-22T12:22:18.684400772Z"
  }
]
```

`--output=text` Result:
```yaml
- content:
    title: Community Pool Spend
    description: Send community pool FUND to the TN Faucet
    recipient: und17jv7rerc2e3undqumpf32a3xs9jc0kjk4z2car
    amount:
    - denom: nund
      amount: "25000000000000"
  id: 7
  proposal_status: 4
  final_tally_result:
    "yes": "0"
    abstain: "0"
    "no": "0"
    no_with_veto: "0"
  submit_time: 2020-03-20T12:22:18.684400772Z
  deposit_end_time: 2020-03-22T12:22:18.684400772Z
  total_deposit:
  - denom: nund
    amount: "1000000000000"
  voting_start_time: 2020-03-20T12:22:18.684400772Z
  voting_end_time: 2020-03-22T12:22:18.684400772Z
```

## undcli query gov vote

Query details for a single vote on a proposal given its identifier.

Usage:
```bash
  undcli query gov vote [proposal-id] [voter-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query gov vote 9 und1nkh2dteta8drxntqp646sr6vz74lt9w9yc60pd
```

`--output=json` Result:
```json
{
  "proposal_id": "9",
  "voter": "und1nkh2dteta8drxntqp646sr6vz74lt9w9yc60pd",
  "option": "Yes"
}
```

`--output=text` Result:
```yaml
proposal_id: 9
voter: und1nkh2dteta8drxntqp646sr6vz74lt9w9yc60pd
option: 1
```

## undcli query gov votes

Query vote details for a single proposal by its identifier.

Usage:
```bash
  undcli query gov votes [proposal-id] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  $ undcli query gov votes 9
  $ undcli query gov votes 9 --page=1 --limit=100
```

`--output=json` Result:
```json
[
  {
    "proposal_id": "9",
    "voter": "und1nkh2dteta8drxntqp646sr6vz74lt9w9yc60pd",
    "option": "Yes"
  }
]
```

`--output=text` Result:
```yaml
- proposal_id: 9
  voter: und1nkh2dteta8drxntqp646sr6vz74lt9w9yc60pd
  option: 1
```

## undcli query gov param

Query the all the parameters for the governance process.

Usage:
```bash
  undcli query gov param [param-type] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  $ undcli query gov param voting
  $ undcli query gov param tallying
  $ undcli query gov param deposit
```

`--output=json` Result:
```json
{
  "voting_period": "172800000000000"
}
```

`--output=text` Result:
```yaml
voting_period: 48h0m0s
```

## undcli query gov params

Query the all the parameters for the governance process.

Usage:
```bash
  undcli query gov params [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query gov params
```

`--output=json` Result:
```json
{
  "voting_params": {
    "voting_period": "172800000000000"
  },
  "tally_params": {
    "quorum": "0.334000000000000000",
    "threshold": "0.500000000000000000",
    "veto": "0.334000000000000000"
  },
  "deposit_params": {
    "min_deposit": [
      {
        "denom": "nund",
        "amount": "1000000000000"
      }
    ],
    "max_deposit_period": "172800000000000"
  }
}
```

`--output=text` Result:
```yaml
voting_params:
  voting_period: 48h0m0s
tally_params:
  quorum: "0.334000000000000000"
  threshold: "0.500000000000000000"
  veto: "0.334000000000000000"
deposit_parmas:
  min_deposit:
  - denom: nund
    amount: "1000000000000"
  max_deposit_period: 48h0m0s
```

## undcli query gov proposer

Query which address proposed a proposal with a given ID.

Usage:
```bash
  undcli query gov proposer [proposal-id] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query gov proposer 10
```

`--output=json` Result:
```json
{
  "proposal_id": "10",
  "proposer": "und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce"
}
```

`--output=text` Result:
```yaml
proposal_id: 10
proposer: und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce
```

## undcli query gov deposit

Query details for a single proposal deposit on a proposal by its identifier.

Usage:
```bash
  undcli query gov deposit [proposal-id] [depositer-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query gov deposit 10 und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce
```

`--output=json` Result:
```json
{
  "proposal_id": "10",
  "depositor": "und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce",
  "amount": [
    {
      "denom": "nund",
      "amount": "1000000000000"
    }
  ]
}
```

`--output=text` Result:
```yaml
proposal_id: 10
depositor: und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce
amount:
- denom: nund
  amount: "1000000000000"
```

## undcli query gov deposits

Query details for all deposits on a proposal.
You can find the proposal-id by running "undcli query gov proposals".

Usage:
```bash
  undcli query gov deposits [proposal-id] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query gov deposits 10
```

`--output=json` Result:
```json
[
  {
    "proposal_id": "10",
    "depositor": "und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce",
    "amount": [
      {
        "denom": "nund",
        "amount": "1000000000000"
      }
    ]
  }
]
```

`--output=text` Result:
```yaml
- proposal_id: 10
  depositor: und1ru9zsaek6nh94e4tzyx49r23mexwwzqd8yasce
  amount:
  - denom: nund
    amount: "1000000000000"
```

## undcli query gov tally

Query tally of votes on a proposal. You can find
the proposal-id by running "undcli query gov proposals".

Usage:
```bash
  undcli query gov tally [proposal-id] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query gov tally 9
```

`--output=json` Result:
```json
{
  "yes": "3714680003497661",
  "abstain": "0",
  "no": "0",
  "no_with_veto": "0"
}
```

`--output=text` Result:
```yaml
"yes": "3714680003497661"
abstain: "0"
"no": "0"
no_with_veto: "0"
```

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

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for signing-info|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query slashing signing-info undvalconspub1addwnpepqddjl38kr5smtcvcgtwen6z373xm8hntdn62905j607l7h3cef5863cwhnf
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query slashing params
```

`--output=json` Result:
```json
{
  "signed_blocks_window": "10000",
  "min_signed_per_window": "0.050000000000000000",
  "downtime_jail_duration": "600000000000",
  "slash_fraction_double_sign": "0.050000000000000000",
  "slash_fraction_downtime": "0.010000000000000000"
}
```

`--output=text` Result:
```yaml
signed_blocks_window: 10000
min_signed_per_window: "0.050000000000000000"
downtime_jail_duration: 10m0s
slash_fraction_double_sign: "0.050000000000000000"
slash_fraction_downtime: "0.010000000000000000"
```

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

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for delegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking delegation und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
```

`--output=json` Result:
```json
{
  "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
  "validator_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
  "shares": "1315507652899786.664739529784222192",
  "balance": {
    "denom": "nund",
    "amount": "1263671402500000"
  }
}
```

`--output=text` Result:
```yaml
delegation:
  delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
  validator_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
  shares: "1315507652899786.664739529784222192"
balance:
  denom: nund
  amount: "1263671402500000"
```

## undcli query staking delegations

Query delegations for an individual delegator on all validators.

Usage:
```bash
  undcli query staking delegations [delegator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking delegations und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
```

`--output=json` Result:
```json
[
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_address": "undvaloper1duyhqzcgrzjy9y2yvueur2h3e2yqxhjl6yvpze",
    "shares": "107080910111213.141516171819202122",
    "balance": {
      "denom": "nund",
      "amount": "104950000000000"
    }
  },
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
    "shares": "106215728125200.313827203266684784",
    "balance": {
      "denom": "nund",
      "amount": "96059601007661"
    }
  },
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
    "shares": "1315507652899786.664739529784222192",
    "balance": {
      "denom": "nund",
      "amount": "1263671402500000"
    }
  }
]
```

`--output=text` Result:
```yaml
- delegation:
    delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_address: undvaloper1duyhqzcgrzjy9y2yvueur2h3e2yqxhjl6yvpze
    shares: "107080910111213.141516171819202122"
  balance:
    denom: nund
    amount: "104950000000000"
- delegation:
    delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
    shares: "106215728125200.313827203266684784"
  balance:
    denom: nund
    amount: "96059601007661"
- delegation:
    delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
    shares: "1315507652899786.664739529784222192"
  balance:
    denom: nund
    amount: "1263671402500000"
```

## undcli query staking unbonding-delegation

Query unbonding delegations for an individual delegator on an individual validator.

Usage:
```bash
  undcli query staking unbonding-delegation [delegator-addr] [validator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking unbonding-delegation und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
```

## undcli query staking unbonding-delegations

Query unbonding delegations for an individual delegator.

Usage:
```bash
  undcli query staking unbonding-delegations [delegator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking unbonding-delegations  und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
```

`--output=json` Result:
```json
[
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
    "entries": [
      {
        "creation_height": "117202",
        "completion_time": "2020-05-08T11:26:07.375057672Z",
        "initial_balance": "10000000000",
        "balance": "10000000000"
      }
    ]
  }
]
```

`--output=text` Result:
```yaml
- delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
  validator_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
  entries:
  - creation_height: 117202
    completion_time: 2020-05-08T11:26:07.375057672Z
    initial_balance: "10000000000"
    balance: "10000000000"
```

## undcli query staking redelegation

Query a redelegation record for an individual delegator between a source and destination validator.

Usage:
```bash
  undcli query staking redelegation [delegator-addr] [src-validator-addr] [dst-validator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking redelegation und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
```

`--output=json` Result:
```json
[
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_src_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
    "validator_dst_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
    "entries": [
      {
        "creation_height": 117259,
        "completion_time": "2020-05-08T11:31:30.934863875Z",
        "initial_balance": "10000000000",
        "shares_dst": "11057273506.344128034033086712",
        "balance": "10000000000"
      }
    ]
  }
]
```

`--output=text` Result:
```yaml
- redelegation:
    delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_src_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
    validator_dst_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
    entries: []
  entries:
  - redelegationentry:
      creation_height: 117259
      completion_time: 2020-05-08T11:31:30.934863875Z
      initial_balance: "10000000000"
      shares_dst: "11057273506.344128034033086712"
    balance: "10000000000"
```

## undcli query staking redelegations

Query all redelegation records for an individual delegator.

Usage:
```bash
  undcli query staking redelegations [delegator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking redelegations und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
```

`--output=json` Result:
```json
[
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_src_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
    "validator_dst_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
    "entries": [
      {
        "creation_height": 117259,
        "completion_time": "2020-05-08T11:31:30.934863875Z",
        "initial_balance": "10000000000",
        "shares_dst": "11057273506.344128034033086712",
        "balance": "10000000000"
      }
    ]
  }
]
```

`--output=text` Result:
```yaml
- redelegation:
    delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_src_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
    validator_dst_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
    entries: []
  entries:
  - redelegationentry:
      creation_height: 117259
      completion_time: 2020-05-08T11:31:30.934863875Z
      initial_balance: "10000000000"
      shares_dst: "11057273506.344128034033086712"
    balance: "10000000000"
```

## undcli query staking validator

Query details about an individual validator.

Usage:
```bash
  undcli query staking validator [validator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking validator undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
```

`--output=json` Result:
```json
{
  "operator_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
  "consensus_pubkey": "undvalconspub1zcjduepq6yq7drzefkavsrxhxk69cy63tj3r7trq4qgvksre47266gpfpevqz8n8h5",
  "jailed": false,
  "status": 2,
  "tokens": "1263661402500000",
  "delegator_shares": "1315497242696229.812572300815615504",
  "description": {
    "moniker": "SerenityTN",
    "identity": "",
    "website": "",
    "security_contact": "",
    "details": "Serenity TestNet"
  },
  "unbonding_height": "567798",
  "unbonding_time": "2020-04-21T02:42:15.384696258Z",
  "commission": {
    "commission_rates": {
      "rate": "0.050000000000000000",
      "max_rate": "0.100000000000000000",
      "max_change_rate": "0.010000000000000000"
    },
    "update_time": "2020-03-19T15:40:11.555272561Z"
  },
  "min_self_delegation": "1"
}
```

`--output=text` Result:
```yaml
|
  operatoraddress: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
  conspubkey: undvalconspub1zcjduepq6yq7drzefkavsrxhxk69cy63tj3r7trq4qgvksre47266gpfpevqz8n8h5
  jailed: false
  status: 2
  tokens: "1263661402500000"
  delegatorshares: "1315497242696229.812572300815615504"
  description:
    moniker: SerenityTN
    identity: ""
    website: ""
    security_contact: ""
    details: Serenity TestNet
  unbondingheight: 567798
  unbondingcompletiontime: 2020-04-21T02:42:15.384696258Z
  commission:
    commission_rates:
      rate: "0.050000000000000000"
      max_rate: "0.100000000000000000"
      max_change_rate: "0.010000000000000000"
    update_time: 2020-03-19T15:40:11.555272561Z
  minselfdelegation: "1"
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query staking validators
```

`--output=json` Result:
```json
[
  {
    "operator_address": "undvaloper1z575nxpemg4dtcrv0rd2u5pwctjxjzsyppgfj4",
    "consensus_pubkey": "undvalconspub1zcjduepqng6dka9glcps04vq67hv46hel3xklmpl00sknek67fdpy4z40hvs5krlz0",
    "jailed": false,
    "status": 2,
    "tokens": "348987245313000",
    "delegator_shares": "356073099983000.000000000000000000",
    "description": {
      "moniker": "集体计算机",
      "identity": "",
      "website": "",
      "security_contact": "",
      "details": ""
    },
    "unbonding_height": "97634",
    "unbonding_time": "2020-05-07T04:00:33.807902347Z",
    "commission": {
      "commission_rates": {
        "rate": "0.400000000000000000",
        "max_rate": "1.000000000000000000",
        "max_change_rate": "0.100000000000000000"
      },
      "update_time": "2020-03-30T06:18:32.22789609Z"
    },
    "min_self_delegation": "1"
  },
  {
    "operator_address": "undvaloper19enk9tmm98sa6yzk4pwarxamkkc0nth9p84vny",
    "consensus_pubkey": "undvalconspub1zcjduepqger0ht3yc4aqylwjlatmmur629wkunqg8jhy9zml2nhe3spp0ajqsar4p0",
    "jailed": false,
    "status": 2,
    "tokens": "297974677770000",
    "delegator_shares": "300984523000000.000000000000000000",
    "description": {
      "moniker": "UNDEurope",
      "identity": "",
      "website": "",
      "security_contact": "",
      "details": ""
    },
    "unbonding_height": "591414",
    "unbonding_time": "2020-04-22T12:07:32.72068041Z",
    "commission": {
      "commission_rates": {
        "rate": "0.100000000000000000",
        "max_rate": "1.000000000000000000",
        "max_change_rate": "0.100000000000000000"
      },
      "update_time": "2020-01-25T05:39:47.228048938Z"
    },
    "min_self_delegation": "1"
  },
  ...
]
```

`--output=text` Result:
```yaml
- |
  operatoraddress: undvaloper1z575nxpemg4dtcrv0rd2u5pwctjxjzsyppgfj4
  conspubkey: undvalconspub1zcjduepqng6dka9glcps04vq67hv46hel3xklmpl00sknek67fdpy4z40hvs5krlz0
  jailed: false
  status: 2
  tokens: "348987245313000"
  delegatorshares: "356073099983000.000000000000000000"
  description:
    moniker: 集体计算机
    identity: ""
    website: ""
    security_contact: ""
    details: ""
  unbondingheight: 97634
  unbondingcompletiontime: 2020-05-07T04:00:33.807902347Z
  commission:
    commission_rates:
      rate: "0.400000000000000000"
      max_rate: "1.000000000000000000"
      max_change_rate: "0.100000000000000000"
    update_time: 2020-03-30T06:18:32.22789609Z
  minselfdelegation: "1"
- |
  operatoraddress: undvaloper19enk9tmm98sa6yzk4pwarxamkkc0nth9p84vny
  conspubkey: undvalconspub1zcjduepqger0ht3yc4aqylwjlatmmur629wkunqg8jhy9zml2nhe3spp0ajqsar4p0
  jailed: false
  status: 2
  tokens: "297974677770000"
  delegatorshares: "300984523000000.000000000000000000"
  description:
    moniker: UNDEurope
    identity: ""
    website: ""
    security_contact: ""
    details: ""
  unbondingheight: 591414
  unbondingcompletiontime: 2020-04-22T12:07:32.72068041Z
  commission:
    commission_rates:
      rate: "0.100000000000000000"
      max_rate: "1.000000000000000000"
      max_change_rate: "0.100000000000000000"
    update_time: 2020-01-25T05:39:47.228048938Z
  minselfdelegation: "1"
  ...
```

## undcli query staking delegations-to

Query delegations on an individual validator.

Usage:
```bash
  undcli query staking delegations-to [validator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking delegations-to undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
```

`--output=json` Result:
```json
[
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
    "shares": "1315497242696229.812572300815615504",
    "balance": {
      "denom": "nund",
      "amount": "1263661402500000"
    }
  }
]
```

`--output=text` Result:
```yaml
- delegation:
    delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
    shares: "1315497242696229.812572300815615504"
  balance:
    denom: nund
    amount: "1263661402500000"
```

## undcli query staking unbonding-delegations-from

Query delegations that are unbonding _from_ a validator.

Usage:
```bash
  undcli query staking unbonding-delegations-from [validator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking unbonding-delegations-from undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
```

`--output=json` Result:
```json
[
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
    "entries": [
      {
        "creation_height": "117202",
        "completion_time": "2020-05-08T11:26:07.375057672Z",
        "initial_balance": "10000000000",
        "balance": "10000000000"
      }
    ]
  }
]
```

`--output=text` Result:
```yaml
- delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
  validator_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
  entries:
  - creation_height: 117202
    completion_time: 2020-05-08T11:26:07.375057672Z
    initial_balance: "10000000000"
    balance: "10000000000"
```

## undcli query staking redelegations-from

Query delegations that are redelegating _from_ a validator.

Usage:
```bash
  undcli query staking redelegations-from [validator-addr] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking redelegations-from undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
```

`--output=json` Result:
```json
[
  {
    "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
    "validator_src_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
    "validator_dst_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
    "entries": [
      {
        "creation_height": 117259,
        "completion_time": "2020-05-08T11:31:30.934863875Z",
        "initial_balance": "10000000000",
        "shares_dst": "11057273506.344128034033086712",
        "balance": "10000000000"
      }
    ]
  }
]
```

`--output=text` Result:
```yaml
- redelegation:
    delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_src_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
    validator_dst_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
    entries: []
  entries:
  - redelegationentry:
      creation_height: 117259
      completion_time: 2020-05-08T11:31:30.934863875Z
      initial_balance: "10000000000"
      shares_dst: "11057273506.344128034033086712"
    balance: "10000000000"
```

## undcli query staking historical-info

Query historical info at given height.

Usage:
```bash
  undcli query staking historical-info [height] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query staking historical-info 5
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query staking params
```

`--output=json` Result:
```json
{
  "unbonding_time": "1814400000000000",
  "max_validators": 96,
  "max_entries": 7,
  "historical_entries": 3,
  "bond_denom": "nund"
}
```

`--output=text` Result:
```yaml
unbonding_time: 504h0m0s
max_validators: 96
max_entries: 7
historical_entries: 3
bond_denom: nund
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query staking pool
```

`--output=json` Result:
```json
{
  "not_bonded_tokens": "171571803499999",
  "bonded_tokens": "5397623673233001"
}
```

`--output=text` Result:
```yaml
not_bonded_tokens: "171571803499999"
bonded_tokens: "5397623673233001"
```

## undcli query supply

Query total FUND supply, including locked and unlocked

Returns three values:

1. **amount**: Liquid FUND in active circulation, and the actual circulating total supply which is available and can be used for FUND transfers, staking, Tx fees etc. It is the **locked** amount subtracted from **total**. _This is the important value when processing any calculations dependent on FUND circulation/total supply of FUND etc._
2. **locked**: Total FUND locked through Enterprise purchases. This FUND is only available specifically to pay WRKChain/BEACON fees and **cannot** be used for transfers, staking/delegation or any other transactions. _Locked FUND only enters the active circulation supply once it has been used to pay for WRKChain/BEACON fees. Until then, it is considered "dormant", and not part of the circulating total supply_
3. **total**: The total amount of FUND currently known on the chain, including any Enterprise **locked** FUND. This is for informational purposes only and should not be used for any "circulating/total supply" calculations.

The **amount** value is the important value regarding total supply _currently in active circulation_, and is the information that should be used to represent any "total supply/circulation" values for example in block explorers, wallets, exchanges etc.

Usage:
```bash
  undcli query supply [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  undcli query supply
```

`--output=json` Result:
```json
{
  "denom": "nund",
  "amount": "120010263000000000",
  "locked": "89737000000000",
  "total": "120100000000000000"
}
```

`--output=text` Result:
```yaml
denom: nund
amount: 120010263000000000
locked: 89737000000000
total: 120100000000000000
```

In the above example, the active circulating supply, usable for transfers and standard transactions etc. is 120,010,263 FUND. 89,737 FUND is currently locked, and can only be used for paying for WRKChain/BEACON fees - it is "dormant" and cannot be used for any other purpose until it has been used to pay for WRKChain/BEACON fees. Finally, the total amount of FUND known on the chain is 120,100,000 FUND, and is the equivalent of 120,010,263 + 89,737.

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query tendermint-validator-set
```

`--output=json` Result:
```json
{
  "block_height": "117406",
  "validators": [
    {
      "address": "undvalcons1zahztaascyqtv0ym8gkcsmvnaugl7whg8slsux",
      "pub_key": "undvalconspub1zcjduepqng6dka9glcps04vq67hv46hel3xklmpl00sknek67fdpy4z40hvs5krlz0",
      "proposer_priority": "-405676367",
      "voting_power": "348987245"
    },
    {
      "address": "undvalcons1dncuh358ehasxa69hwlwmysw5dqztl6jcukydq",
      "pub_key": "undvalconspub1zcjduepq6ftvp4cw0t89swuf6wmj7dh9j3khjjxja96qrmn85uc0metnh5zsyc4aas",
      "proposer_priority": "-573761026",
      "voting_power": "350004403"
    },
    {
      "address": "undvalcons1wuf5pvrwpld77skgjvt20j2xr367ycm0f55avf",
      "pub_key": "undvalconspub1zcjduepqf25xmd7q00x0mxekuq96uhrye6tpm2a8l9n5ehvlr2te2zkk7past6g4lu",
      "proposer_priority": "1658834369",
      "voting_power": "2502048332"
    },
    {
      "address": "undvalcons1jvykwrpv9hzqyrcxa820tuvy5qfpk9a47utec8",
      "pub_key": "undvalconspub1zcjduepq6yq7drzefkavsrxhxk69cy63tj3r7trq4qgvksre47266gpfpevqz8n8h5",
      "proposer_priority": "3531853998",
      "voting_power": "1263661402"
    },
    {
      "address": "undvalcons1n75a4hwv8cn3smm8r6e6ssgpks82700w9c763t",
      "pub_key": "undvalconspub1zcjduepqrcgsx6380dpcgdvns8mdcexd402dpe06ewvk4aphjvdt267tg8kqhu4umx",
      "proposer_priority": "880424921",
      "voting_power": "490050000"
    },
    {
      "address": "undvalcons1hnhzr5xfhd7uv32fvens59qda8t86jywd88prz",
      "pub_key": "undvalconspub1zcjduepqger0ht3yc4aqylwjlatmmur629wkunqg8jhy9zml2nhe3spp0ajqsar4p0",
      "proposer_priority": "-552115870",
      "voting_power": "297974677"
    },
    {
      "address": "undvalcons1ezskagtc4lmr6wsnp6tvgztkf2scq8rwvrhrez",
      "pub_key": "undvalconspub1zcjduepq2d87vxx5g7px5hwdx07mvfnj24rt9j52cs4j49glf388jlwyjcxs2h4fau",
      "proposer_priority": "-2185936894",
      "voting_power": "45211448"
    },
    {
      "address": "undvalcons170cp2v6pnwefvayxtrjh6u3kftprhd9ud5jy0c",
      "pub_key": "undvalconspub1zcjduepq7ayhnappxxzwm23l4gu3zl7s2tggdwaw4mxydf4uq5j5r7l35p3qzla3l5",
      "proposer_priority": "-2353623126",
      "voting_power": "99686164"
    }
  ]
}
```

`--output=text` Result:
```yaml
blockheight: 117412
validators:
- address: undvalcons1zahztaascyqtv0ym8gkcsmvnaugl7whg8slsux
  pubkey: undvalconspub1zcjduepqng6dka9glcps04vq67hv46hel3xklmpl00sknek67fdpy4z40hvs5krlz0
  proposerpriority: 2037234348
  votingpower: 348987245
- address: undvalcons1dncuh358ehasxa69hwlwmysw5dqztl6jcukydq
  pubkey: undvalconspub1zcjduepq6ftvp4cw0t89swuf6wmj7dh9j3khjjxja96qrmn85uc0metnh5zsyc4aas
  proposerpriority: 1876269795
  votingpower: 350004403
- address: undvalcons1wuf5pvrwpld77skgjvt20j2xr367ycm0f55avf
  pubkey: undvalconspub1zcjduepqf25xmd7q00x0mxekuq96uhrye6tpm2a8l9n5ehvlr2te2zkk7past6g4lu
  proposerpriority: -2417321991
  votingpower: 2502048332
- address: undvalcons1jvykwrpv9hzqyrcxa820tuvy5qfpk9a47utec8
  pubkey: undvalconspub1zcjduepq6yq7drzefkavsrxhxk69cy63tj3r7trq4qgvksre47266gpfpevqz8n8h5
  proposerpriority: 1582236470
  votingpower: 1263661402
- address: undvalcons1n75a4hwv8cn3smm8r6e6ssgpks82700w9c763t
  pubkey: undvalconspub1zcjduepqrcgsx6380dpcgdvns8mdcexd402dpe06ewvk4aphjvdt267tg8kqhu4umx
  proposerpriority: -1086848750
  votingpower: 490050000
- address: undvalcons1hnhzr5xfhd7uv32fvens59qda8t86jywd88prz
  pubkey: undvalconspub1zcjduepqger0ht3yc4aqylwjlatmmur629wkunqg8jhy9zml2nhe3spp0ajqsar4p0
  proposerpriority: 1533706869
  votingpower: 297974677
- address: undvalcons1ezskagtc4lmr6wsnp6tvgztkf2scq8rwvrhrez
  pubkey: undvalconspub1zcjduepq2d87vxx5g7px5hwdx07mvfnj24rt9j52cs4j49glf388jlwyjcxs2h4fau
  proposerpriority: -1869456758
  votingpower: 45211448
- address: undvalcons170cp2v6pnwefvayxtrjh6u3kftprhd9ud5jy0c
  pubkey: undvalconspub1zcjduepq7ayhnappxxzwm23l4gu3zl7s2tggdwaw4mxydf4uq5j5r7l35p3qzla3l5
  proposerpriority: -1655819978
  votingpower: 99686164
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query tx DE11EF61545D742B9EBB4D43B7A86EC36412B401CCF77A44403A49B67CCC4DF1
```

`--output=json` Result:
```json
{
  "height": "117259",
  "txhash": "DE11EF61545D742B9EBB4D43B7A86EC36412B401CCF77A44403A49B67CCC4DF1",
  "data": "0C089286D5F5051083C8E3BD03",
  "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"begin_redelegate\"},{\"key\":\"sender\",\"value\":\"und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e\"},{\"key\":\"sender\",\"value\":\"und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e\"},{\"key\":\"module\",\"value\":\"staking\"},{\"key\":\"sender\",\"value\":\"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk\"}]},{\"type\":\"redelegate\",\"attributes\":[{\"key\":\"source_validator\",\"value\":\"undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps\"},{\"key\":\"destination_validator\",\"value\":\"undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg\"},{\"key\":\"amount\",\"value\":\"10000000000\"},{\"key\":\"completion_time\",\"value\":\"2020-05-08T11:31:30Z\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk\"},{\"key\":\"amount\",\"value\":\"1078652753472nund\"},{\"key\":\"recipient\",\"value\":\"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk\"},{\"key\":\"amount\",\"value\":\"91319655nund\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "begin_redelegate"
            },
            {
              "key": "sender",
              "value": "und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e"
            },
            {
              "key": "sender",
              "value": "und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e"
            },
            {
              "key": "module",
              "value": "staking"
            },
            {
              "key": "sender",
              "value": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"
            }
          ]
        },
        {
          "type": "redelegate",
          "attributes": [
            {
              "key": "source_validator",
              "value": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps"
            },
            {
              "key": "destination_validator",
              "value": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg"
            },
            {
              "key": "amount",
              "value": "10000000000"
            },
            {
              "key": "completion_time",
              "value": "2020-05-08T11:31:30Z"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"
            },
            {
              "key": "amount",
              "value": "1078652753472nund"
            },
            {
              "key": "recipient",
              "value": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"
            },
            {
              "key": "amount",
              "value": "91319655nund"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "250000",
  "gas_used": "230308",
  "tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "cosmos-sdk/MsgBeginRedelegate",
          "value": {
            "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
            "validator_src_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps",
            "validator_dst_address": "undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg",
            "amount": {
              "denom": "nund",
              "amount": "10000000000"
            }
          }
        }
      ],
      "fee": {
        "amount": [
          {
            "denom": "nund",
            "amount": "15000"
          }
        ],
        "gas": "250000"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A1svxPYdIbXhmELdmehR9E2z3mts9KK+ktP9/144ymh9"
          },
          "signature": "GC1gcask7LY8SLbVYfhbBPA2pCoZQFtW6Yd60HCRKedHppAudBK4ikn/VLijMAeWutWStj1CptWU2tDEFjv4RQ=="
        }
      ],
      "memo": "sent from Unification Web Wallet"
    }
  },
  "timestamp": "2020-04-17T11:31:30Z"
}
```

`--output=text` Result:
```yaml
height: 117259
txhash: DE11EF61545D742B9EBB4D43B7A86EC36412B401CCF77A44403A49B67CCC4DF1
codespace: ""
code: 0
data: 0C089286D5F5051083C8E3BD03
rawlog: '[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"begin_redelegate"},{"key":"sender","value":"und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e"},{"key":"sender","value":"und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e"},{"key":"module","value":"staking"},{"key":"sender","value":"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"}]},{"type":"redelegate","attributes":[{"key":"source_validator","value":"undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps"},{"key":"destination_validator","value":"undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg"},{"key":"amount","value":"10000000000"},{"key":"completion_time","value":"2020-05-08T11:31:30Z"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"},{"key":"amount","value":"1078652753472nund"},{"key":"recipient","value":"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"},{"key":"amount","value":"91319655nund"}]}]}]'
logs:
- msgindex: 0
  log: ""
  events:
  - type: message
    attributes:
    - key: action
      value: begin_redelegate
    - key: sender
      value: und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e
    - key: sender
      value: und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e
    - key: module
      value: staking
    - key: sender
      value: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
  - type: redelegate
    attributes:
    - key: source_validator
      value: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
    - key: destination_validator
      value: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
    - key: amount
      value: "10000000000"
    - key: completion_time
      value: "2020-05-08T11:31:30Z"
  - type: transfer
    attributes:
    - key: recipient
      value: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    - key: amount
      value: 1078652753472nund
    - key: recipient
      value: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    - key: amount
      value: 91319655nund
info: ""
gaswanted: 250000
gasused: 230308
tx:
  msg:
  - delegator_address: und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk
    validator_src_address: undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps
    validator_dst_address: undvaloper1w2dlf0793gk3m5zk8e554stg97x7uw95dfx4kg
    amount:
      denom: nund
      amount: "10000000000"
  fee:
    amount:
    - denom: nund
      amount: "15000"
    gas: 250000
  signatures:
  - |
    pubkey: undpub1addwnpepqddjl38kr5smtcvcgtwen6z373xm8hntdn62905j607l7h3cef586k283s9
    signature: !!binary |
      GC1gcask7LY8SLbVYfhbBPA2pCoZQFtW6Yd60HCRKedHppAudBK4ikn/VLijMAeWutWStj
      1CptWU2tDEFjv4RQ==
  memo: sent from Unification Web Wallet
timestamp: "2020-04-17T11:31:30Z"
```


## undcli query txs

Search for transactions that match the exact given events where results are paginated.

Each event takes the form of '`{eventType}.{eventAttribute}={value}`'. Please refer to each module's documentation for the full set of events to query for. Each module documents its respective events under 'xx_events.md'.

Usage:
```bash
  undcli query txs [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
undcli query txs --events 'message.sender=und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk&message.action=withdraw_delegator_reward' --page 1 --limit 30
```


`--output=json` Result:
```json
{
  "total_count": "3",
  "count": "3",
  "page_number": "1",
  "page_total": "1",
  "limit": "30",
  "txs": [
    {
      "height": "12750",
      "txhash": "03CF3FD6780A7EC9E833D673B1FACAF92D8E4212397A12BFB41E2255BB909101",
      "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"withdraw_delegator_reward\"},{\"key\":\"sender\",\"value\":\"und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e\"},{\"key\":\"module\",\"value\":\"distribution\"},{\"key\":\"sender\",\"value\":\"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk\"},{\"key\":\"amount\",\"value\":\"291320378636nund\"}]},{\"type\":\"withdraw_rewards\",\"attributes\":[{\"key\":\"amount\",\"value\":\"291320378636nund\"},{\"key\":\"validator\",\"value\":\"undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps\"}]}]}]",
      "logs": [
        {
          "msg_index": 0,
          "log": "",
          "events": [
            {
              "type": "message",
              "attributes": [
                {
                  "key": "action",
                  "value": "withdraw_delegator_reward"
                },
                {
                  "key": "sender",
                  "value": "und1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8q6r40e"
                },
                {
                  "key": "module",
                  "value": "distribution"
                },
                {
                  "key": "sender",
                  "value": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"
                }
              ]
            },
            {
              "type": "transfer",
              "attributes": [
                {
                  "key": "recipient",
                  "value": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk"
                },
                {
                  "key": "amount",
                  "value": "291320378636nund"
                }
              ]
            },
            {
              "type": "withdraw_rewards",
              "attributes": [
                {
                  "key": "amount",
                  "value": "291320378636nund"
                },
                {
                  "key": "validator",
                  "value": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps"
                }
              ]
            }
          ]
        }
      ],
      "gas_wanted": "190000",
      "gas_used": "106818",
      "tx": {
        "type": "cosmos-sdk/StdTx",
        "value": {
          "msg": [
            {
              "type": "cosmos-sdk/MsgWithdrawDelegationReward",
              "value": {
                "delegator_address": "und16twxa6lyj7uhp56tukrcfz2p6q93mrxgqvrspk",
                "validator_address": "undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps"
              }
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "nund",
                "amount": "6000"
              }
            ],
            "gas": "190000"
          },
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A1svxPYdIbXhmELdmehR9E2z3mts9KK+ktP9/144ymh9"
              },
              "signature": "eNsxrhygORZHPNp2KkhfpgJTRyNrE2Tx7aDvdhjzSQ1SW5pz8UQO4OATASTgc0Ip3kgllFz/ar5kJPMP/wWELg=="
            }
          ],
          "memo": "sent from Unification Web Wallet"
        }
      },
      "timestamp": "2020-04-10T12:58:38Z"
    },
    ...
  ]
}

```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query wrkchain params
```

`--output=json` Result:
```json
{
  "fee_register": "1000000000",
  "fee_record": "1000000000",
  "denom": "nund"
}
```

`--output=text` Result:
```yaml
fee_register: 1000000000
fee_record: 1000000000
denom: nund
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query wrkchain wrkchain 1
```

`--output=json` Result:
```json
{
  "wrkchain_id": "1",
  "moniker": "finchain",
  "name": "Finchain TestNet",
  "genesis": "19e4596ad881f2fbf2d80c04c2228b06bafc75e2ae611b2e19f9164b2c901354",
  "type": "geth",
  "lastblock": "309",
  "num_blocks": "65",
  "reg_time": "1584004250",
  "owner": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs"
}
```

`--output=text` Result:
```yaml
wrkchainid: 1
moniker: finchain
name: Finchain TestNet
genesishash: 19e4596ad881f2fbf2d80c04c2228b06bafc75e2ae611b2e19f9164b2c901354
basetype: geth
lastblock: 309
numberblocks: 65
registertime: 1584004250
owner: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
```

## undcli query wrkchain search

Query for all paginated WRKChains that match optional filters:

Usage:
```bash
  undcli query wrkchain search [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--limit`|`int`|pagination limit to query for (default 100)|
|`--page`|`int`|pagination page to query for|
|`--moniker`|`string`|(optional) filter wrkchains by moniker|
|`--owner`|`string`|(optional) filter wrkchains by owner address|
|`--height`|`int`|Use a specific height to query state at (this can error if the node is pruning state)|
|`-h`, `--help`||help for redelegation|
|`--indent`||Add indent to JSON response|
|`--ledger`||Use a connected Ledger device|
|`--node`|`string`|\<host\>:\<port\> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")|
|`--trust-node`||Trust connected full node (don't verify proofs for responses)|

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:
```bash
  $ undcli query wrkchain search --moniker wrkchain1
  $ undcli query wrkchain search --owner und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
  $ undcli query wrkchain search --page=2 --limit=100
```

`--output=json` Result:
```json
[
  {
    "wrkchain_id": "1",
    "moniker": "finchain",
    "name": "Finchain TestNet",
    "genesis": "19e4596ad881f2fbf2d80c04c2228b06bafc75e2ae611b2e19f9164b2c901354",
    "type": "geth",
    "lastblock": "309",
    "num_blocks": "65",
    "reg_time": "1584004250",
    "owner": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs"
  },
  {
    "wrkchain_id": "2",
    "moniker": "finchain-tn",
    "name": "Finchain TestNet",
    "genesis": "19e4596ad881f2fbf2d80c04c2228b06bafc75e2ae611b2e19f9164b2c901354",
    "type": "geth",
    "lastblock": "206694",
    "num_blocks": "34771",
    "reg_time": "1584012884",
    "owner": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs"
  }
]
```

`--output=text` Result:
```yaml
- wrkchainid: 1
  moniker: finchain
  name: Finchain TestNet
  genesishash: 19e4596ad881f2fbf2d80c04c2228b06bafc75e2ae611b2e19f9164b2c901354
  basetype: geth
  lastblock: 309
  numberblocks: 65
  registertime: 1584004250
  owner: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
- wrkchainid: 2
  moniker: finchain-tn
  name: Finchain TestNet
  genesishash: 19e4596ad881f2fbf2d80c04c2228b06bafc75e2ae611b2e19f9164b2c901354
  basetype: geth
  lastblock: 206694
  numberblocks: 34771
  registertime: 1584012884
  owner: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

Example:

```bash
  undcli query wrkchain block 2 206650
```

`--output=json` Result:
```json
{
  "wrkchain_id": "2",
  "height": "206650",
  "blockhash": "0xcd843fef654343f1a2befa00aadd69d424cbbfe8732a005e0d8cf46b0f25c0a4",
  "parenthash": "0x253420699f7af31eaff0b0d3ff4c5c43f4e39cdb6081fbd7ec595fba2f5e62c4",
  "hash1": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  "hash2": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  "hash3": "0xfa190168b7c394c91006f3670aae52bf0dbdadafe3051ed8398e003c0effc9d4",
  "sub_time": "1587123683",
  "owner": "und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs"
}
```

`--output=text` Result:
```yaml
wrkchainid: 2
height: 206650
blockhash: 0xcd843fef654343f1a2befa00aadd69d424cbbfe8732a005e0d8cf46b0f25c0a4
parenthash: 0x253420699f7af31eaff0b0d3ff4c5c43f4e39cdb6081fbd7ec595fba2f5e62c4
hash1: 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
hash2: 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
hash3: 0xfa190168b7c394c91006f3670aae52bf0dbdadafe3051ed8398e003c0effc9d4
submittime: 1587123683
owner: und12zns8tfm0g2rskl4f9zg2hr9n53agkyvtftngs
```

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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
|[distribution](#undcli-tx-distribution)|Distribution transactions subcommands|
|[enterprise](#undcli-tx-enterprise)|Enterprise FUND transaction subcommands|
|[gov](#undcli-tx-gov)|Governance transactions subcommands|
|[slashing](#undcli-tx-slashing)|Slashing transactions subcommands|
|[staking](#undcli-tx-staking)|Staking transaction subcommands|
|[wrkchain](#undcli-tx-wrkchain)|WRKChain transaction subcommands|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for tx|

## undcli tx broadcast

Broadcast transactions created with the `--generate-only`
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx encode

Encode transactions created with the `--generate-only` flag and signed with the sign command.
Read a transaction from \<file\>, serialise it to the Amino wire protocol, and output it as base64.
If you supply a dash (-) argument in place of an input filename, the command reads from standard input.

Usage:
```bash
  undcli tx encode [file] [flags]
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|


## undcli tx decode

Decode an amino-encoded transaction string.

Usage:
```bash
  undcli tx decode [amino-byte-string] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-x`, `--hex`||Treat input as hexadecimal instead of base64|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx multisign

Alias of [undcli tx auth multisign](#undcli-tx-auth-multisign)

Sign transactions created with the --generate-only flag that require multisig signatures.

Read signature(s) from `[signature]` file(s), generate a multisig signature compliant to the multisig key `[name]`, and attach it to the transaction read from `[file]`.

Example:
```bash
  undcli multisign transaction.json k1k2k3 k1sig.json k2sig.json k3sig.json
```

If the flag `--signature-only` flag is on, it outputs a JSON representation
of the generated signature only.

The `--offline` flag makes sure that the client will not reach out to an external node. Thus account number or sequence number lookups will not be performed and it is recommended to set such parameters manually.

Usage:
```bash
  undcli tx multisign [file] [name] [[signature]...] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--offline`||Offline mode. Do not query a full node|
|`--output-document`|`string`|The document will be written to the given file instead of `STDOUT`|
|`--signature-only`||Print only the generated signature, then exit|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|


## undcli tx send

Alias of [undcli tx bank send](#undcli-tx-bank-send)

Create and sign a send tx. Amount to send is in `nund`, e.g. `1000000000nund`.

Usage:
```bash
  undcli tx send [from_key_or_address] [to_address] [amount] [flags]
```

Example:
```bash
  undcli tx send my-wallet und1hp2km26czxlvesn8nmwswdd90umvcm5gxwpk98 1000000000nund
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx sign

Alias of [undcli tx auth sign](#undcli-tx-auth-sign)

Sign transactions created with the `--generate-only` flag.
It will read a transaction from `[file]`, sign it, and print its JSON encoding.

If the flag `--signature-only` flag is set, it will output a JSON representation
of the generated signature only.

If the flag `--validate-signatures` is set, then the command would check whether all required signers have signed the transactions, whether the signatures were collected in the right order, and if the signature is valid over the given transaction.

If the `--offline` flag is also set, signature validation over the transaction will be not be performed as that will require RPC communication with a full node.

The `--offline` flag makes sure that the client will not reach out to full node.
As a result, the account and sequence number queries will not be performed and
it is required to set such parameters manually. Note, invalid values will cause
the transaction to fail.

The `--multisig=<multisig_key>` flag generates a signature on behalf of a multisig account key. It implies `--signature-only`. Full multisig signed transactions may eventually be generated via the '`multisign`' command.

Usage:
```bash
  undcli tx sign [file] [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--append`||Append the signature to the existing ones. If disabled, old signatures would be overwritten. Ignored if` --multisig` is on (default true)
|`--multisig`|`string`|Address of the multisig account on behalf of which the transaction shall be signed|
|`--offline`||Offline mode; Do not query a full node. `--account` and `--sequence` options would be required if offline is set|
|`--output-document`|`string`|The document will be written to the given file instead of `STDOUT`|
|`--signature-only`||Print only the generated signature, then exit|
|`--validate-signatures`||Print the addresses that must sign the transaction, those who have already signed it, and make sure that signatures are in the correct order|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx auth

Auth transaction subcommands

Usage:
```bash
  undcli tx auth [flags]
  undcli tx auth [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[multisign](#undcli-tx-auth-multisign)|Generate multisig signatures for transactions generated offline|
|[sign](#undcli-tx-auth-sign)|Sign transactions generated offline|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for auth|

## undcli tx auth multisign

See [undcli tx multisign](#undcli-tx-multisign)

## undcli tx auth sign

See [undcli tx sign](#undcli-tx-sign)

## undcli tx bank

Bank transaction subcommands

Usage:
```bash
  undcli tx bank [flags]
  undcli tx bank [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[send](#undcli-tx-bank-send)|Create and sign a send tx|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for bank|

## undcli tx bank send

See [undcli tx send](#undcli-tx-send)

## undcli tx beacon

Beacon transaction subcommands

Usage:
```bash
  undcli tx beacon [flags]
  undcli tx beacon [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[register](#undcli-tx-beacon-register)|register a new BEACON|
|[record](#undcli-tx-beacon-record)|record a BEACON's timestamp hash|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for beacon|

## undcli tx beacon register

Register a new BEACON, to enable timestamp hash submissions.

The BEACON registration fees are automatically calculated and applied to the transaction. Fees can be queried using the `undcli query beacon params` command.

Usage:
```bash
  undcli tx beacon register [flags]
```

::: tip
The `--moniker` flag is required to register a BEACON.
:::

Example:
```bash
  undcli tx beacon register --moniker=MyBeacon --name="My BEACON" --from mykey
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--moniker`|`string`|BEACON's moniker|
|`--name`|`string`|(optional) BEACON's name|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx beacon record

Record a BEACON's' timestamp hash.

::: tip Note
The `--hash` flag is required to record a BEACON hash.  
If the `--subtime` is not set, the current UTC UNIX time will be used.
:::

Usage:
```bash
  undcli tx beacon record [beacon id] [flags]
```

Example:
```bash
  undcli tx beacon record 1 --hash=d04b98f48e8 --subtime=1234356 --from mykey
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--hash`|`string`|BEACON's timestamp hash|
|`--subtime`|`uint`|BEACON's timestamp submission time|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx distribution

Distribution transactions subcommands

Usage:
```bash
  undcli tx distribution [flags]
  undcli tx distribution [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[withdraw-rewards](#undcli-tx-distribution-withdraw-rewards)|Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator|
|[set-withdraw-addr](#undcli-tx-distribution-set-withdraw-addr)|change the default withdraw address for rewards associated with an address|
|[withdraw-all-rewards](#undcli-tx-distribution-withdraw-all-rewards)|withdraw all delegations rewards for a delegator|
|[fund-community-pool](#undcli-tx-distribution-fund-community-pool)|Funds the community pool with the specified amount|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for distribution|

## undcli tx distribution withdraw-rewards

Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator.

Usage:
```bash
  undcli tx distribution withdraw-rewards [validator-addr] [flags]
```

Example:
```bash
  $ undcli tx distribution withdraw-rewards undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey
  $ undcli tx distribution withdraw-rewards undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey --commission
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--commission`||also withdraw validator's commission|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx distribution set-withdraw-addr

Set the withdraw address for rewards associated with a delegator address.

Usage:
```bash
  undcli tx distribution set-withdraw-addr [withdraw-addr] [flags]
```

Example:
```bash
$ undcli tx distribution set-withdraw-addr und1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx distribution withdraw-all-rewards

Withdraw all rewards for a single delegator.

Usage:
```bash
  undcli tx distribution withdraw-all-rewards [flags]
```

Example:
```bash
  undcli tx distribution withdraw-all-rewards --from mykey
```

Flags:

| Flag | Type | Description |
|------|------|-------------|
|`--max-msgs`|`int`|Limit the number of messages per tx (0 for unlimited) (default 5)|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx distribution fund-community-pool

Funds the community pool with the specified amount

Usage:
```bash
  undcli tx distribution fund-community-pool [amount] [flags]
```

Example:
```bash
  undcli tx distribution fund-community-pool 1000000000nund --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx enterprise

Enterprise FUND transaction subcommands

Usage:
```bash
  undcli tx enterprise [flags]
  undcli tx enterprise [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[purchase](#undcli-tx-enterprise-purchase)|Raise a new Enterprise FUND purchase order|
|[process](#undcli-tx-enterprise-process)|Process an Enterprise FUND purchase order|
|[whitelist](#undcli-tx-enterprise-whitelist)|Add/Remove an address from the enterprise purchase order whitelist|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for enterprise|

## undcli tx enterprise purchase

Raise a new Enterprise FUND purchase order.

Usage:
```bash
  undcli tx enterprise purchase [amount] [flags]
```

Example:
```bash
  $ undcli tx enterprise purchase 1000000000000nund --from wrktest
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx enterprise process

Process an Enterprise FUND purchase order.

Only authorised addresses may process purchase orders

`[decision]` must be `accept` or `reject`

Usage:
```bash
  undcli tx enterprise process [purchase_order_id] [decision] [flags]
```

Example:
```bash
  $ undcli tx enterprise process 24 accept --from ent
  $ undcli tx enterprise process 24 reject --from ent
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx enterprise whitelist

Add/Remove an address from the enterprise purchase order whitelist.

Only authorised addesses may edit the whitelist.

`[action]` must be `add` or `remove`

Usage:
```bash
  undcli tx enterprise whitelist [action] [address] [flags]
```

Example:
```bash
  $ undcli tx enterprise whitelist add und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy --from ent
  $ undcli tx enterprise whitelist remove und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy --from ent
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|


## undcli tx gov

Governance transactions subcommands

Usage:
```bash
  undcli tx gov [flags]
  undcli tx gov [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[deposit](#undcli-tx-gov-deposit)|Deposit tokens for an active proposal|
|[vote](#undcli-tx-gov-vote)|Vote for an active proposal, options: yes/no/no_with_veto/abstain|
|[submit-proposal](#undcli-tx-gov-submit-proposal)|Submit a proposal along with an initial deposit|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for gov|

## undcli tx gov deposit

Submit a deposit for an active proposal. You can find the `[proposal-id]` by running "`undcli query gov proposals`".

Usage:
```bash
  undcli tx gov deposit [proposal-id] [deposit] [flags]
```

Example:
```bash
  undcli tx gov deposit 1 10stake --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx gov vote

Submit a vote for an active proposal. You can find the `[proposal-id]` by running "`undcli query gov proposals`".

Usage:
```bash
  undcli tx gov vote [proposal-id] [option] [flags]
```

Example:
```bash
  undcli tx gov vote 1 yes --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx gov submit-proposal

Submit a proposal along with an initial deposit.

Proposal title, description, type and deposit can be given directly or through a proposal JSON file.

Usage:
```bash
  undcli tx gov submit-proposal [flags]
  undcli tx gov submit-proposal [command]
```

Example:
```bash
  undcli tx gov submit-proposal --proposal="path/to/proposal.json" --from mykey
```

Where `proposal.json` contains:

```json
{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "10test"
}
```

Which is equivalent to:

```bash
  undcli tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10test" --from mykey
```

Available Commands:
| Command | Description |
|---------|-------------|
|[param-change](#undcli-tx-gov-submit-proposal-param-change)|Submit a parameter change proposal|
|[community-pool-spend](#undcli-tx-gov-submit-proposal-community-pool-spend)|Submit a community pool spend proposal|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--deposit`|`string`|deposit of proposal, e.g. 100000000000nund|
|`--description`|`string`|description of proposal|
|`--proposal`|`string`|proposal file path (if this path is given, other proposal flags are ignored)|
|`--title`|`string`|title of proposal|
|`--type`|`string`|Type of proposal, types: `text`/`parameter_change`/`software_upgrade`|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx gov submit-proposal param-change

Submit a parameter proposal along with an initial deposit.

The proposal details must be supplied via a JSON file. For values that contains
objects, only non-empty fields will be updated.

::: warning IMPORTANT
Currently parameter changes are evaluated but not validated, so it is
very important that any "value" change is valid (ie. correct type and within bounds) for its respective parameter, eg. "MaxValidators" should be an integer and not a decimal.

Proper vetting of a parameter change proposal should prevent this from happening
(no deposits should occur during the governance process), but it should be noted
regardless.
:::

Usage:
```bash
  undcli tx gov submit-proposal param-change [proposal-file] [flags]
```

Example:
```bash
$ undcli tx gov submit-proposal param-change <path/to/proposal.json> --from=<key_or_address>
```

Where `proposal.json` contains:
```json
{
  "title": "Slashing parameters",
  "description": "change the signed blocks window to 10,000, and minimum signed requirement to 5%",
  "changes": [
    {
      "subspace": "slashing",
      "key": "SignedBlocksWindow",
      "value": "10000"
    },
    {
      "subspace": "slashing",
      "key": "MinSignedPerWindow",
      "value": "0.050000000000000000"
    }
  ],
  "deposit": [
    {
      "denom": "nund",
      "amount": "1000000000000"
    }
  ]
}
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx gov submit-proposal community-pool-spend

Submit a community pool spend proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Usage:
```bash
  undcli tx gov submit-proposal community-pool-spend [proposal-file] [flags]
```

Example:
```bash
  undcli tx gov submit-proposal community-pool-spend <path/to/proposal.json> --from=<key_or_address>
```

Where `proposal.json` contains:

```json
{
  "title": "Community Pool Spend",
  "description": "Send some community pool FUND to this address",
  "recipient": "und17jv7rerc2e3undqumpf32a3xs9jc0kjk4z2car",
  "amount": [
    {
      "denom": "nund",
      "amount": "10000"
    }
  ],
  "deposit": [
    {
      "denom": "nund",
      "amount": "10000"
    }
  ]
}
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx slashing

Slashing transactions subcommands

Usage:
```bash
  undcli tx slashing [flags]
  undcli tx slashing [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[unjail](#undcli-tx-slashing-unjail)|unjail validator previously jailed for downtime|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for slashing|

## undcli tx slashing unjail

unjail a jailed validator:

Usage:
```bash
  undcli tx slashing unjail [flags]
```

Example:
```bash
  undcli tx slashing unjail --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx staking

Staking transaction subcommands

Usage:
```bash
  undcli tx staking [flags]
  undcli tx staking [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[create-validator](#undcli-tx-staking-create-validator)|create new validator initialised with a self-delegation to it|
|[edit-validator](#undcli-tx-staking-edit-validator)|edit an existing validator account|
|[delegate](#undcli-tx-staking-delegate)|Delegate UND to a validator|
|[redelegate](#undcli-tx-staking-redelegate)|Redelegate illiquid tokens from one validator to another|
|[unbond](#undcli-tx-staking-unbond)|Unbond shares from a validator|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for staking|

## undcli tx staking create-validator

create new validator initialised with a self-delegation to it

Usage:
```bash
  undcli tx staking create-validator [flags]
```

Example:
```bash
undcli tx staking create-validator \
  --amount=1000000000000nund \
  --pubkey=undvalconspub1zcjduepq6yq7drzefkavsrxhxk69cy63tj3r... \
  --moniker="MyAwesomeNode" \
  --website="https://my-node-site.com" \
  --details="My node is awesome" \
  --security-contact="security@my-node-site.com" \
  --commission-rate="0.05" \
  --commission-max-rate="0.10" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=my_new_wallet
```

::: warning
The values for `--commission-max-change-rate` and `--commission-max-rate` flags cannot be changed after the create-validator command has been run.
:::

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--amount`|`string`|Amount of coins to bond|
|`--commission-max-change-rate`|`string`|The maximum commission change rate percentage (per day)|
|`--commission-max-rate`|`string`|The maximum commission rate percentage|
|`--commission-rate`|`string`|The initial commission rate percentage|
|`--details`|`string`|The validator's (optional) details|
|`--identity`|`string`|The optional identity signature (ex. UPort or Keybase)|
|`--ip`|`string`|The node's public IP. It takes effect only when used in combination with `--generate-only`|
|`--min-self-delegation`|`string`|The minimum self delegation required on the validator|
|`--moniker`|`string`|The validator's name|
|`--node-id`|`string`|The node's ID|
|`--pubkey`|`string`|The Bech32 encoded PubKey of the validator|
|`--security-contact`|`string`|The validator's (optional) security contact email|
|`--website`|`string`|The validator's (optional) website|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx staking edit-validator

edit an existing validator account

Usage:
```bash
  undcli tx staking edit-validator [flags]
```

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--commission-rate`|`string`|The new commission rate percentage|
|`--details`|`string`|The validator's (optional) details (default "`[do-not-modify]`")|
|`--identity`|`string`|The (optional) identity signature (ex. UPort or Keybase) (default "`[do-not-modify]`")|
|`--min-self-delegation`|`string`|The minimum self delegation required on the validator|
|`--moniker`|`string`|The validator's name (default "`[do-not-modify]`")|
|`--security-contact`|`string`|The validator's (optional) security contact email (default "`[do-not-modify]`")|
|`--website`|`string`|The validator's (optional) website (default "`[do-not-modify]`")|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx staking delegate

Delegate an amount of FUND (in `nund`) to a validator from your wallet.

Usage:
```bash
  undcli tx staking delegate [validator-addr] [amount] [flags]
```

Example:
```bash
$ undcli tx staking delegate undvaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 1000000000nund --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx staking redelegate

Redelegate an amount of liquid staking tokens from one validator to another.

Usage:
```bash
  undcli tx staking redelegate [src-validator-addr] [dst-validator-addr] [amount] [flags]
```

Example:
```bash
$ undcli tx staking redelegate undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj undvaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 100nund --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx staking unbond

Unbond an amount of bonded shares from a validator.

Usage:
```bash
  undcli tx staking unbond [validator-addr] [amount] [flags]
```

Example:
```bash
  undcli tx staking unbond undvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100nund --from mykey
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx wrkchain

WRKChain transaction subcommands

Usage:
```bash
  undcli tx wrkchain [flags]
  undcli tx wrkchain [command]
```

Available Commands:
| Command | Description |
|---------|-------------|
|[register](#undcli-tx-wrkchain-register)|register a new WRKChain|
|[record](#undcli-tx-wrkchain-record)|record a WRKChain's block hashes|

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`-h`, `--help`||help for wrkchain|

## undcli tx wrkchain register

Register a new WRKChain, to enable WRKChain hash submissions

Usage:
```bash
  undcli tx wrkchain register [flags]
```

Example:
```bash
  undcli tx wrkchain register --moniker="MyWrkChain" --genesis="d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa" --name="My WRKChain" --base="geth" --from mykey
```

::: warning Note
The `--moniker` and `--base` flags are the minimum requirements for registering a WRKChain, and are mandatory flags.
:::

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--base`|`string`|(optional) WRKChain's chain type - `geth`/`tendermint`, etc.|
|`--genesis`|`string`|(optional) WRKChain's Genesis hash|
|`--moniker`|`string`|WRKChain's moniker|
|`--name`|`string`|(optional) WRKChain's name|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

## undcli tx wrkchain record

Record a new WRKChain block's hash(es)

Usage:
```bash
  undcli tx wrkchain record [wrkchain id] [flags]
```

Example:
```bash
  $ undcli tx wrkchain record 1 --wc_height=24 --block_hash="d04b98f48e8" --parent_hash="f8bcc15c6ae" --hash1="5ac050801cd6" --hash2="dcfd428fb5f9e" --hash3="65c4e16e7807340fa" --from mykey
  $ undcli tx wrkchain record 1 --wc_height=25 --block_hash="d04b98f48e8" --from mykey
  $ undcli tx wrkchain record 1 --wc_height=26 --block_hash="d04b98f48e8" --parent_hash="f8bcc15c6ae" --from mykey
```

::: warning Note
The `--wc_height` and `--block_hash` are the minimum requirements for submitting WRKChain block header hashes, and are mandatory flags.
:::

Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--block_hash`|`string`|WRKChain block's header (main) hash|
|`--hash1`|`string`|(optional) Additional WRKChain hash - e.g. State Merkle Root|
|`--hash2`|`string`|(optional) Additional WRKChain hash - e.g. Tx Merkle Root|
|`--hash3`|`string`|(optional) Additional WRKChain hash|
|`--parent_hash`|`string`|(optional) WRKChain block's parent hash|
|`--wc_height`|`uint`|WRKChain block's height/block number|
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

Global Flags:
| Flag | Type | Description |
|------|------|-------------|
|`--chain-id`|`string`|Chain ID of und Mainchain node|
|`-e`, `--encoding`|`string`|Binary encoding (`hex`\|`b64`\|`btc`) (default "`hex`")|
|`--keyring-backend`|`string`|Select keyring's backend (`os`\|`file`\|`test`) (default "`os`")|
|`-h`, `--help`||help for undcli|
|`--home`|`string`|directory for config and data (default "`$HOME/.und_cli`")|
|`-o`, `--output`|`string`|Output format (`text`\|`json`) (default "`text`")|
|`--trace`||print out full stack trace on errors|

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
