# Frequently Asked Questions

#### Contents

[[toc]]

## 1. What is Unification Mainchain?

Unification is a scalable master blockchain for Enterprise.

1. WRKChains can be considered similar to side-chains which allow scaling processing power and cost metrics for enterprises who need an immutable blockchain without being directly on a public blockchain
2. Mainchain  is the master 100% public blockchain which WRKChains consume FUND to submit their block header hashes to Mainchain for public immutability.

### 1.1. What is the official project website?

[https://unification.com](https://unification.com)

### 1.2. What are the Github repo addresses?

- [https://github.com/unification-com](https://github.com/unification-com)  
- Mainchain repo - [https://github.com/unification-com/mainchain](https://github.com/unification-com/mainchain)

## 2. Where can I find the main API documents?

- `und` (server) command reference: [https://docs.unification.io/software/und-commands.html](https://docs.unification.io/software/und-commands.html)
- `undcli` (client) command reference: [https://docs.unification.io/software/undcli-commands.html](https://docs.unification.io/software/undcli-commands.html)
- Public REST API: [https://rest-testnet.unification.io/swagger-ui/](https://rest-testnet.unification.io/swagger-ui/) (Note: link for TestNet)
- Public RPC Interface: [https://rpc1-testnet.unification.io:26657](https://rpc1-testnet.unification.io:26657) (Nore: link for TestNet)

Public REST/RPC APIs for MainNet will be available via [https://rest.unification.io/swagger-ui](https://rest.unification.io/swagger-ui) and [https://rpc1.unification.io:26657/](https://rpc1.unification.io:26657/) respectively on launch.

## 3. Where is the block explorer?

- TestNet - [https://explorer-testnet.unification.io](https://explorer-testnet.unification.io)
- MainNet will be available on [https://explorer.unification.io](https://explorer.unification.io) when launched

## 4. Where can I get the pre-compiled binaries?

Latest compiled binaries available from [https://github.com/unification-com/mainchain/releases](https://github.com/unification-com/mainchain/releases)  

- `und` (server/full node) software has been compiled for Linux x86_64. Tested on CentOS 7, and Ubuntu 16.04  
- `undcli` (client) has been compiled for Linux x86_64, Windows and OSX. Tested on CentOS 7, Ubuntu 16.04, Windows 10.

## 5. How do I compile the code from source?

Build instructions are available here: [https://docs.unification.io/software/installation.html#building-from-source](https://docs.unification.io/software/installation.html#building-from-source)

This will build and install both `und` and `undcli` binaries into `$GOPATH/bin`

## 6. Key operation instructions

### 6.1. How to get the block height?

Several methods available:

```bash
undcli status --chain-id=CHAIN_ID --node=tcp://NODE:PORT
```

TestNet Example:

```bash
undcli status --chain-id=FUND-Mainchain-TestNet-v7 --node=tcp://rpc1-testnet.unification.io:26657
```

A JSON or text object is returned, and latest height available from `sync_info.latest_block_height`

RPC equivalent on TestNet: [https://rpc1-testnet.unification.io:26657/status](https://rpc1-testnet.unification.io:26657/status) and MainNet: [https://rpc1.unification.io:26657/status](https://rpc1.unification.io:26657/status)

`undcli query block` can also be used:

```bash
undcli query block --chain-id=CHAIN_ID --node=tcp://NODE:PORT
```

Will return the latest block info if no height is passed to the query.

TestNet example, using the public RPC node:

```bash
undcli query block --chain-id=FUND-Mainchain-TestNet-v7 --node=tcp://rpc1-testnet.unification.io:26657 --trust-node=false
```

### 6.2. How do I create new wallet address?

```bash
undcli keys add ACC_NAME
```

`ACC_NAME` is whatever ASCII identifier you want to give the account/wallet/address, and is used to reference the account when creating/signing Txs.

The command will output pertinent information - name (as passed in the command), wallet address, public key and recovery mnemonic in either JSON or text format.

Example:

```bash
undcli keys add some_new_account
```

Run `undcli keys add --help` or see [https://docs.unification.io/software/undcli-commands.html#undcli-keys-add](https://docs.unification.io/software/undcli-commands.html#undcli-keys-add) for details on flags/command options  etc.

### 6.3. How to transfer FUND?

```bash
undcli tx send [from_key_or_address] [to_address] [amount] --chain-id=CHAIN_ID --node=tcp://NODE_IP:PORT
```

Amount is `nund` - "Nano Unification Denomination", such that **1,000,000,000 nund == 1 FUND**. See [https://docs.unification.io/introduction/denomination.html](https://docs.unification.io/introduction/denomination.html).

Example to send **10 FUND** from `my_account` account (see Q7.2 about account names) on TestNet, using the public RPC node:

```bash
undcli tx send my_account und1nkhnc5e8pvph4phv93k0lkscc7yf5eh9kas5f6 10000000000nund --chain-id=FUND-Mainchain-TestNet-v7 --node=tcp://rpc1-testnet.unification.io:26657 --gas=auto --gas-adjustment=1.5 --gas-prices=0.25nund --trust-node=false
```

See [https://docs.unification.io/software/undcli-commands.html#undcli-tx-send](https://docs.unification.io/software/undcli-commands.html#undcli-tx-send) and [https://docs.unification.io/introduction/fees-and-gas.html](https://docs.unification.io/introduction/fees-and-gas.html) for more in-depth information.

### 6.4. How do I get all transactions related to one wallet/account?

`undcli query txs` can be used to query all transactions. Passing the `--events` flag will allow you to filter indexed events by a particular account. Data is returned paginated.

TestNet example to get Txs sent by `und17jv7rerc2e3undqumpf32a3xs9jc0kjk4z2car`, using the public RPC node:

```bash
undcli query txs --events 'message.sender=und17jv7rerc2e3undqumpf32a3xs9jc0kjk4z2car' --chain-id=FUND-Mainchain-TestNet-v7 --node=tcp://rpc1-testnet.unification.io:26657 --page 1 --limit 30
```

The `--events` flag can contain any `{eventType}.{eventAttribute}={value}` type query. For example `--events 'transfer.recipient=und17jv7rerc2e3undqumpf32a3xs9jc0kjk4z2car'` will return queries relating to transfers into the account. See [https://docs.unification.io/software/undcli-commands.html#undcli-query-txs](https://docs.unification.io/software/undcli-commands.html#undcli-query-txs) for further information.

### 6.5. How do I get the FUND balance for one wallet/account?

```bash
undcli query account [address] [flags]
```

Example on TestNet, using the public RPC node:

```bash
undcli query account und1eyn7s6qz2gcnfld0uskwxedyunpgjhlcjhvul9 --chain-id=FUND-Mainchain-TestNet-v7 --node=tcp://rpc1-testnet.unification.io:26657
```

Will return a JSON or text object (depending on options passed). `account.value.coins` in the returned result shows the amount of `nund`. The above example (currently) shows the account has **10000000000 nund (10 FUND)** on TestNet.

### 6.6. How do I export（dump/backup）a wallet?

```bash
undcli keys export some_new_account
```

will export an account private key in ASCII-armored encrypted format.

### 6.7. How do I import a wallet?

There are a couple of methods, depending on the import format. If the bip39 mnemonic is available, then:

```bash
undcli keys add some_new_account --recover
```

Will prompt you for the bip39 mnemonic. See [https://docs.unification.io/software/undcli-commands.html#undcli-keys-add](https://docs.unification.io/software/undcli-commands.html#undcli-keys-add)

If the private key has been exported (e.g. via `undcli keys export`), then the `undcli keys import` command can be used:

```bash
undcli keys import ACC_NAME KEYFILE
```

See [https://docs.unification.io/software/undcli-commands.html#undcli-keys-import](https://docs.unification.io/software/undcli-commands.html#undcli-keys-import)
