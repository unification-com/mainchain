# Becoming a TestNet Validator

If you intend to become a [MainNet Validator](become-mainnet-validator.md), it will be good practice to fist run a validator node on TestNet to become familiar with the process.

::: warning IMPORTANT
Whenever you use `undcli` to send Txs or query the chain ensure you pass the correct data to the `--chain-id` and if necessary `--node=` flags so that you connect to the correct network!
:::

#### Contents

[[toc]]

## Prerequisites

Before continuing, ensure you have gone through the following docs:

1. [Installing the software](installation.md)
2. [Join the Public TestNet](join-testnet.md)
3. [Accounts and Wallets](accounts-wallets.md)
4. **TestNet** Chain ID - if you haven't already, you can get the current chain ID by running:

```
jq --raw-output '.chain_id' $HOME/.und_mainchain/config/genesis.json
```

The above command assumes you have downloaded the TestNet genesis to the default `$HOME/.und_mainchain` directory.

::: warning IMPORTANT
you will need an account with sufficient UND to self-delegate to your validator node.
:::

::: tip
if you intend to fully participate in the running of **TestNet**, your node will need to be permanently available and online. In which case, you will need to investigate running `und` as a [background service](run-und-as-service.md)
:::

## Creating a validator

The first thing you will need is your node's Tendermint validator public key. This will be used to register your node as a Validator on the network. To get the key, open a terminal and run:

```bash
und tendermint show-validator
```
This will output your node's public key. Make a note of it, as it will be required soon.

::: warning IMPORTANT
Before continuing, ensure your full node has fully synced with TestNet and downloaded all the blocks (this may take a while, so go and make a brew). You can check the current network block height in the [TestNet block explorer](https://explorer-testnet.unification.io). For each block your full node syncs, you will see:

`I[2020-01-15|11:45:07.782] Executed block module=state height=37285 validTxs=0 invalidTxs=0`

When the height is the same as the current network block number, your full node has completed syncing.
:::

To create your Validator, you will need to generate, sign and broadcast a special transaction to the network which will register your Tendermint validator public key and stake the amount of UND specified (via self-delegation). Run the following command, modifying as required:

```bash
undcli tx staking create-validator \
  --amount=STAKE_IN_NUND \
  --pubkey=NODE_TENDERMINT_PUBLIC_KEY \
  --moniker="YOUR_EV_MONIKER" \
  --website="YOUR_WEBSITE_URL" \
  --identity=16_DIGIT_KEYBASE_IO_ID \
  --details="NODE_DESCRIPTION" \
  --security-contact="SECURITY_CONTACT_EMAIL" \
  --chain-id=CHAIN_ID \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-prices="0.25nund" \
  --gas-adjustment=1.5 \
  --from=SELF_DELEGATOR_ACCOUNT
```

**Mandatory fields**

`STAKE_IN_NUND`: the amount in nund you want to delegate to yourself. For example, if you want to stake 1000 UND, enter 1000000000000nund.

::: tip
You can use the built in `undcli` conversion tool to calculate this:

```
undcli convert 1000 und nund.
```
:::

::: warning IMPORTANT
do not attempt to stake more than you have in your account, and ensure you have enough UND to pay for the transaction fees, and enough left over for future transactions!
:::

`NODE_TENDERMINT_PUBLIC_KEY`: Your node's tendermint public key, obtained earlier via the `und tendermint show-validator` command.

`CHAIN_ID`: the chain you are creating a validator for. This was obtained previously via the `jq` command.

`SELF_DELEGATOR_ACCOUNT`: the name of the account being used to stake self-delegated UND and sign the transaction — for example, the identifier you entered when running the `undcli keys add` command earlier to create/import an account.

`YOUR_EV_MONIKER`: a moniker which will publicly identify your EV on the network.

**Optional fields**

::: tip
Ensure you create your validator with as much of the following additional information as you can. It will be publicly visible, and help potential stakers connect with you
:::

`YOUR_WEBSITE_URL`: the URL for the site promoting your validation node

`16_DIGIT_KEYBASE_IO_ID`: Your 16 digit public [keybase.io](https://keybase.io) PGP public key ID if you have one and want to associate your ID to your validator node.

`NODE_DESCRIPTION`: a brief description of your validator node

`SECURITY_CONTACT_EMAIL`: Email address for the security contact for your validator node

**Commission Rates**

Your commission rates can be set using the `--commission-rate` , `--commission-max-change-rate` and `--commission-max-rate` flags.

`--commission-rate`: The % commission you will earn from delegators’ rewards. Keeping this low can attract more delegators to your node.

`--commission-max-rate`: The maximum you will ever increase your commission rate to — you cannot raise commission above this value. Again, keeping this low can attract more delegators.

`--commission-max-change-rate`: The maximum you can change the commission-rate by in any one change request. For example, if your maximum change rate is 0.01, you can only make changes in 0.01 increments, so from 0.10 (10%) to 0.09 (9%).

::: warning IMPORTANT
The values for `--commission-max-change-rate` and `--commission-max-rate` flags cannot be changed after the create-validator command has been run.
:::

Finally, the `--min-self-delegation` flag is the minimum amount of `nund` you are required to keep self-delegated to your validator, meaning you must always have _at least_ this amount self-delegated to your node.

**Example**

```bash
undcli tx staking create-validator \
--amount=1000000000000nund \
--pubkey=undvalconspub1zcjduepq6yq7drzefkavsrxhxk69cy63tj3r... \
--moniker="MyAwesomeNode" \
--website="https://my-node-site.com" \
--details="My node is awesome" \
--security-contact="security@my-node-site.com" \
--chain-id=UND-Mainchain-TestNet-v3 \
--commission-rate="0.05" \
--commission-max-rate="0.10" \
--commission-max-change-rate="0.01" \
--min-self-delegation="1000000000" \
--gas="auto" \
--gas-prices="0.25nund" \
--gas-adjustment=1.5 \
--from=my_new_wallet
```

The command will return a Tx hash, which you can use to query whether or not the transaction was successful:

```bash
undcli query tx TX_HASH --chain-id UND-Mainchain-TestNet-v3
```

::: tip
you can set the `--broadcast-mode` flag in the command to `block`. This will tell `undcli` to wait for the transaction to be processed in a block before returning the result. This will take up to 5-6 seconds to complete, but the Tx result will be included in the output.
:::

### Verify

You can verify your node is registered as a validator by running:

```bash
undcli query staking validator \
$(undcli keys show SELF_DELEGATOR_ACCOUNT --bech=val -a) \
--chain-id=CHAIN_ID
```

replacing `SELF_DELEGATOR_ACCOUNT` and `CHAIN_ID` accordingly.
