# WRKChain Example

This document will guide you though registering a new WRKChain and submitting hashes
via a manual process. Although hash submission is usually automated with the WRKOracle
software, this guide will help you understand the process, and allow for testing.

**Note**: It is _HIGHLY_ recommended that you only undertake this guide on
either [DevNet](local-devnet.md) or [TestNet](join-testnet.md). WRKChain hash 
submission on MainNet should preferably be automated using the WRKOracle
software.

## Registering your WRKChain

Registration is required so that the WRKChain has an identifier on Mainchain.
Registration incurs a one-time fee.

The following `undcli` command can be used to register your WRKChain, and assumes
you have a local full node running:

```bash
undcli tx wrkchain register --moniker="[YOUR_MONIKER]" --genesis="[GENESIS_BLOCK_HASH]" --name="[WRKCHAIN NAME]" --base="geth" --from [from_account] --chain-id [chain_id] --gas=auto --gas-adjustment=1.25
```

- `[YOUR_MONIKER]` - any alphanumeric identifier you want to give your WRKChain, e.g. wrkchain1
- `[GENESIS_BLOCK_HASH]` - the hash value for your genesis block. The `--genesis` flag is optional
- `[WRKCHAIN NAME]` - a name for your WRKChain, e.g. My First WRKChain. Optional
- `[from_account]` - your local account identifier (see [Accounts and Wallets](accounts-wallets.md))
- `[chain_id]` - the ID of the chain to run the transaction on

**Note**: The `--base` flag is used to define the base chain type your WRKChain has been
built with - for example "geth" (for a `go-ethereum` based WRKChain), "hyperledger" etc.)

For example, we have a local account and key set up called "testwrk", and want
to register a new WRKChain, with the moniker "wrkchain1" called "WRKChain Example 1":

```bash
undcli tx wrkchain register --moniker="wrkchain1" --genesis="78521D6EFBEDF6D7EE9C73EDD3443B8021DADBE06ECE81F639B6EC57D8E3F3EA" --name="WRKChain Example 1" --base="geth" --from testwrk --chain-id UND-Mainchain-TestNet --gas=auto --gas-adjustment=1.25
```

You will be prompted to accept, and enter testwrk's account password. Once
broadcast, you will receive confirmation with the TX Hash, which can be used
to query the Tx.

Your WRKCHain's on-chain ID will be embedded in the Tx query result. Alternatively,
you can run a search query to get your WRKChain's details, for example:

```bash
undcli query wrkchain search --moniker wrkchain1 --chain-id UND-Mainchain-TestNet
```

will return a result similar to:

```json
[
  {
    "wrkchain_id": "101",
    "moniker": "wrkchain1",
    "name": "WRKChain Example 1",
    "genesis": "78521D6EFBEDF6D7EE9C73EDD3443B8021DADBE06ECE81F639B6EC57D8E3F3EA",
    "type": "geth",
    "lastblock": "1",
    "reg_time": "1576858904",
    "owner": "und1n0d2qre7hrshud600rdtk4427428rjvnewnqfc"
  }
]
```

The `wrkchain_id` value is what is required to submit hashes, and find your WRKChain's submitted block hashes.

The `lastblock` tells us which block number was last submitted for the
WRKChain.

## Recording Hashes

Once successfully registered, you will be able to submit block hashes however
frequently suits your needs. To simulate how thw WRKOracle works, we can run the following
command to submit hashes:

```bash
undcli tx wrkchain record 1 --wc_height=[BLOCK_HEIGHT] --block_hash=[BLOCK_HASH] --parent_hash=[PARENT_HASH] --hash1=[HASH1] --hash2=[HASH2] --hash3=[HASH3] --from [account_name] --chain-id [chain_id] --gas=auto --gas-adjustment=1.5
```

- `[BLOCK_HEIGHT]` - the height/number of the WRKChain block you are submitting hashes for
- `[BLOCK_HASH]` - the main block hash
- `[PARENT_HASH]` - the main block's parent block hash (optional)
- `[HASH1]` - an optional, arbitrary hash. This can be, for example, the Tx Merkle root hash
- `[HASH2]` - an optional, arbitrary hash. This can be, for example, the Tx Merkle root hash
- `[HASH3]` - an optional, arbitrary hash. This can be, for example, the Tx Merkle root hash
- `[from_account]` - your local account identifier (see [Accounts and Wallets](accounts-wallets.md))
- `[chain_id]` - the ID of the chain to run the transaction on

For example, if we just want to submit the block hash for block number 123, we can run:

```bash
undcli tx wrkchain record 1 --wc_height=123 --block_hash=1BB457C575E72D7401C809B66290FAC56347223912F2484BA7E881D42495CD0F --from testwrk --chain-id UND-Mainchain-TestNet --gas=auto --gas-adjustment=1.5
```

## Querying a WRKChain on Mainchain

To retrieve a particular hash submitted for a WRKChain, we can run:

```bash
undcli query wrkchain block [WRKCHAIN_ID] [HEIGHT]
```

- `[WRKCHAIN_ID]` - the numeric ID for the WRKChain
- `[HEIGHT]` - the block number we wish to retrieve

If `[HEIGHT]` has been submitted for `[WRKCHAIN_ID]`, the data will be
returned in a JSON object, If not, the returned object will contain empty 
values, meaning the WRKChain has not submitted a value for this block
height.
