# WRKChain Example

This document will guide you though registering a new WRKChain and submitting hashes via a manual process. Although hash
submission is usually automated with the WRKOracle software, this guide will help you understand the process, and allow
for testing.

::: warning IMPORTANT
Whenever you use `und` to send Txs or query the chain ensure you pass the correct data to the `--chain-id` and if
necessary `--node=` flags so that you connect to the correct network!
:::

::: tip Note
It is _HIGHLY_ recommended that you only undertake this guide on
either [DevNet](local-devnet.md) or TestNet first.
:::

#### Contents

[[toc]]

## Registering your WRKChain

Registration is required so that the WRKChain has an identifier on Mainchain.
Registration incurs a one-time fee.

The following `und` command can be used to register your WRKChain, and assumes you have a local full node running,
connected to either DevNet or TestNet:

```bash
und tx wrkchain register --moniker="[YOUR_MONIKER]" --genesis="[GENESIS_BLOCK_HASH]" --name="[WRKCHAIN NAME]" --base="[WRKCHAIN_TYPE]" --from [from_account] --chain-id [chain_id] --gas=auto --gas-adjustment=1.5
```

- `[YOUR_MONIKER]` - any alphanumeric identifier you want to give your WRKChain, e.g. wrkchain1
- `[GENESIS_BLOCK_HASH]` - the hash value for your genesis block. The `--genesis` flag is optional. The hash can be
  obtained by running your genesis through a sha256 generator, for example.
- `[WRKCHAIN NAME]` - a name for your WRKChain, e.g. My First WRKChain. Optional
- `[from_account]` - your local account identifier (see [Accounts and Wallets](accounts-wallets.md)). This will be used
  as the WRKChain Owner. **Only the owner will be able to submit block hashes, so it is important to keep this account
  safe!**
- `[chain_id]` - the ID of the chain to run the transaction on
- `[WRKCHAIN_TYPE]` - the type of WRKChain. Currently supported by WRKOracle are WRKChains built using `cosmos`, `eos`
  , `geth`, `neo`, `stellar`, `tendermint`.

For example, we have a local account and key set up called "testwrk", and want
to register a new WRKChain, with the moniker "wrkchain1" called "WRKChain Example 1":

```bash
und tx wrkchain register --moniker="wrkchain1" --genesis="78521D6EFBEDF6D7EE9C73EDD3443B8021DADBE06ECE81F639B6EC57D8E3F3EA" --name="WRKChain Example 1" --base="tendermint" --from testwrk --chain-id FUND-DevNet-2 --gas=auto --gas-adjustment=1.25
```

Once broadcast, you will receive confirmation with the TX Hash, which can be used to query the Tx.

Your WRKCHain's on-chain ID will be embedded in the Tx query result. Alternatively, you can run a search query to get
your WRKChain's details, for example: 

```bash
und query wrkchain search --moniker wrkchain1 --chain-id FUND-DevNet-2
```

will return a result similar to:

```json
[
  {
    "wrkchain_id": "1",
    "moniker": "wrkchain1",
    "name": "WRKChain Example 1",
    "genesis": "78521D6EFBEDF6D7EE9C73EDD3443B8021DADBE06ECE81F639B6EC57D8E3F3EA",
    "type": "tendermint",
    "lastblock": "0",
    "num_blocks": "0",
    "reg_time": "1585752449",
    "owner": "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy"
  }
]

```

The `wrkchain_id` value is what is required to submit hashes, and find your WRKChain's submitted block hashes.

The `lastblock` tells us for which block number hashes were last submitted for the WRKChain, and `num_blocks` the number
of block hashes the WRKChain has submitted in total. Finally, `reg_time` is a UNIX timestamp for when the WRKChain was
registered on Mainchain.

> **Important**: Only the `owner` - i.e. the account used to register the WRKChain - will be able to submit block
> hashes.

## Filter WRKChain metadata

WRKChain metadata can be searched for on Mainchain:

```
und query wrkchain search --moniker wrkchain1
und query wrkchain search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
und query wrkchain search --page=2 --limit=100
```

## Recording Hashes

Once successfully registered, you will be able to submit block hashes however
frequently suits your needs. To simulate how the WRKOracle works, we can run the following command to submit hashes:

```bash
und tx wrkchain record 1 --wc_height=[BLOCK_HEIGHT] --block_hash=[BLOCK_HASH] --parent_hash=[PARENT_HASH] --hash1=[HASH1] --hash2=[HASH2] --hash3=[HASH3] --from [account_name] --chain-id [chain_id] --gas=auto --gas-adjustment=1.5
```

- `[BLOCK_HEIGHT]` - the height/number of the WRKChain block you are submitting hashes for
- `[BLOCK_HASH]` - the main block hash
- `[PARENT_HASH]` - the main block's parent block hash (optional)
- `[HASH1]` - an optional, arbitrary hash. This can be, for example, the Tx Merkle root hash
- `[HASH2]` - an optional, arbitrary hash. This can be, for example, the Tx Merkle root hash
- `[HASH3]` - an optional, arbitrary hash. This can be, for example, the Tx Merkle root hash
- `[from_account]` - your local account identifier (see [Accounts and Wallets](accounts-wallets.md)). This **must** be
  the same as the **owner** used to register the WRKChain.
- `[chain_id]` - the ID of the chain to run the transaction on

For example, if we just want to submit the block hash for block number 123, we can run:

```bash
und tx wrkchain record 1 --wc_height=123 --block_hash=1BB457C575E72D7401C809B66290FAC56347223912F2484BA7E881D42495CD0F --from testwrk --chain-id FUND-DevNet-2 --gas=auto --gas-adjustment=1.5
```

## Querying a WRKChain on Mainchain

To retrieve a particular hash submitted for a WRKChain, we can run:

```bash
und query wrkchain block [WRKCHAIN_ID] [HEIGHT]
```

- `[WRKCHAIN_ID]` - the numeric ID for the WRKChain
- `[HEIGHT]` - the block number we wish to retrieve

If `[HEIGHT]` has been submitted for `[WRKCHAIN_ID]`, the data will be
returned in a JSON object, If not, the returned object will contain empty
values, meaning the WRKChain has not submitted a value for this block
height.
