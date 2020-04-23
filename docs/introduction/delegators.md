# Introduction to Delegators and Staking

## Overview

Many users may want to participate in the running of Unification Mainchain without running a full node and becoming a validator operator. Delegators can stake their UND to a validator in order to participate.

The current pool of 96 Validator nodes is selected based on not only the self-delegated UND staked by the Validator operator(s), but additionally the total UND delegated to it by other users.

Delegators play an important role in the staking process, and indeed the running of the network itself, by acting as a further safeguard against any potential misbehaviour of validators. If delegators feel that a validator is not behaving in the best interests of the network, they can simply move their staked UND away from them. If the validator's total stake falls below the top 96 validator stakes, they will be removed from the active validator pool.

Additionally, Delegators can (and should) take part in network governance by voting on proposals.

Delegating comes with both [risks](#risks) and [rewards](#rewards). Delegators share a percentage of the UND earned by their chosen validator(s) from processing transactions, and singing/producing blocks. The amount earned proportional the amount staked. Risks come from losing a small amount of staked UND should the validator misbehave - this includes prolonged periods of downtime, and more importantly, double-signing blocks.

## The Delegation process

### Selecting a Validator

Information regarding validators can be obtained from a number of sources, including block explorers, or the `undcli` query [validators](../software/undcli-commands.html#undcli-query-staking-validators) and query [validator](../software/undcli-commands.html#undcli-query-staking-validator) commands. There are several pieces of information a Validator should provide in order to help you make a decision:

- **Moniker** - a short identifier/name for the Validator
- **Description** - a brief description of the validator
- **Website** - a link to their Website
- **Security Contact** - an email address of the person(s) responsible for maintaining the node and its security
- **Commission Rate** - the percentage commission the validator charges delegators, and is deducted from the delegator's rewards
- **Maximum Commission** - the maximum percentage a validator can ever charge. This value is set by the validator when registering to become a validator and can never be changed. You may want to be wary, for example, of validators with very high maximum commission rates.
- **Maximum Rate Change** - the maximum daily percentage a validator can increase their commission rate.
- **Minimum self-delegation** - the minimum `nund` a validator can self-delegate. If their self-delegation drops below this amount (for example, they manually unbond, or through slashing due to bad behaviour), all delegations are automatically unbonded. This ensures validators behave in the best interest of the network. Higher numbers are better.

### Delegating your UND

The process is simple - a user sends a special delegation transaction, which tells the network to "bond" the chosen amount of UND to the selected validator. This can be done via the `undcli` CLI's [delegate](../software/undcli-commands.html#undcli-tx-staking-delegate) command, or via the [Web Wallet](https://chrome.google.com/webstore/detail/mkjjflkhdddfjhonakofipfojoepfndk) Chrome extension.

::: tip IMPORTANT
When delegating UND to a validator, ownership of the UND being staked is **NEVER** actually transferred anywhere. You will **ALWAYS** retain 100% ownership and full control of that UND, since it is simply "flagged" in your wallet as being bonded to a Validator. The validator has absolutely **ZERO** control over your UND.
:::

**Example using `undcli`**

You have done some research, and found a validator candidate to delegate 1000 UND to. You need a wallet with sufficient UND, and the Validator's Operator Address. This is different to a standard `und` address and begins with `undvaloper`, for example `undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps`.

Assuming you have [imported/added](../software/undcli-commands.html#undcli-keys-add) your wallet key into `undcli`'s keychain, you would run:

```bash
undcli tx staking delegate undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps 1000000000000nund \
  --gas=auto \
  --gas-adjustment=1.5 \
  --gas-prices=0.25nund \
  --node=http://[full-node-ip]:26656 \
  --chain-id=[chain_id] \
  --trust-node=false \
  --from=my_account
```

replacing `[full-node-ip]` and `[chain_id]` with the relevant IP and chain ID respectively.

## Delegator's Roles

Delegators may participate in several functions regarding the running of the network, including:

- **Exercising due diligence when selecting a validator to delegate to**: an important fist step, before delegating UND is to ensure that the chosen validator has a history of good behaviour. There are several tools available, from block explorers to the `undcli` command line tools for querying a validator.

- **Monitore the validator's behaviour after delegation**: this includes ensuring the validator maintains high uptimes (does not frequently miss blocks), does not double-sign blocks, and participates in governance.

- **Participate in network governance**: Delegators can and should participate in governance by voting on proposals. Similarly to rewards, "voting power" is proportional to the amount staked. By default, delegators who do not vote inherit their validator's voting decision, but can override this by voting themselves.

- **Unbonding from misbehaving validators to hold them accountable**: Delegators who feel that their selected validator is not behaving in the best interests of the network should remove their stake from that validator to reduce their chances of being included in the active validator pool. This can be done in one of two ways: unbonding, which simply removes the delegator's stake from the validator; or re-delegating - switching the delegated stake from one validator to another. Unbonding has a cool-down period of 3 weeks to process, but re-delegation is instant.

## Rewards

Every transaction sent to the network has a fee paid. Some Tx fees, such as submitting WRKChain hashes are fixed at 1 UND per Tx, and others are more flexible depending on how much as user is willing to pay. With each block, the fees are distributed among the active validators and their delegators as rewards. Rewards paid are proportional to the amount staked.

You can monitor and withdraw your rewards as often as you like with either the `undcli` [rewards](../software/undcli-commands.html#undcli-query-distribution-rewards) query and [withdraw-rewards](../software/undcli-commands.html#undcli-tx-distribution-withdraw-rewards) Tx, or via the [Web Wallet](https://chrome.google.com/webstore/detail/mkjjflkhdddfjhonakofipfojoepfndk) Chrome extension.

**Example using `undcli`**

Using the same Validator operator address as above, any outstanding rewards due can be queried using:

```bash
undcli query distribution rewards [my_delegator_address] undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps \
  --node=http://[full-node-ip]:26656 \
  --chain-id=[chain_id]
```
and withdraw any rewards to your account using:

```bash
undcli tx distribution withdraw-rewards undvaloper16twxa6lyj7uhp56tukrcfz2p6q93mrxgt60mps \
  --gas=auto \
  --gas-adjustment=1.5 \
  --gas-prices=0.25nund \
  --node=http://[full-node-ip]:26656 \
  --chain-id=[chain_id] \
  --trust-node=false \
  --from=my_account
```

once again replacing `[full-node-ip]` and `[chain_id]` with the relevant IP and chain ID respectively.

## Risks

Delegators that stake to a validator who continuously misbehave run the risk of having their stake slashed (as does the validator themselves) by a small percentage. Misbehaving includes double signing blocks, and prolonged periods of node downtime. This makes it all the more important to pay attention to a validator's behaviour history (e.g. if they have been slashed before) before selecting a validator to delegate to.

Slashing pays a role in further incentivising validators to perform well, and helps ensure their delegators hold them accountable.
