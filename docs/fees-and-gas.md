# Fees and Gas

Transactions consume gas, and the sender must pay a fee in order for the transaction be processed by the validator nodes. The fee is calculated from the amount of gas a Tx will consume multiplied by the gas price.

::: tip NOTE
The gas price for a transaction is set by the sender of the Tx, but each validator will have set their own `minimum-gas-prices` value, and will not process transactions that do not meet this minimum requirement.
:::

Fees are paid in `nund`, and may be either set or calculated depending on which flags are passed to the `undcli` command.

::: tip NOTE
only `--fees` or `--gas-prices` may be used - not both at the same time.
:::

::: tip
`--gas-prices` can be used along with the `--gas=auto` and `--gas-adjustment` flags to estimate the gas requirement and automatically calculate the Tx fees.
:::

## Example 1: setting --fees

In this example, we're manually setting the fee for a standard `send` transaction. The validator has a `minimum-gas-prices` of `0.25nund`. We'll set the `--fee` to 20000nund. For the purposes of simpler calculations, we'll assume the amount of **gas** consumed for this `send` transaction, including a small `--memo` will be around 65000.

::: tip NOTE
gas is defined on the chain as a flat cost per byte for a Tx, e.g. 10 gas per byte. The total size of our Tx will be around 6500 bytes, and therefore the gas consumed by the Tx will be 6500 * 10 = 65000.
:::

```
undcli tx send [from] [to] 123456nund --memo="some und from me to you" --fees=20000nund
```

In this instance, the `gas-price` is implied as approximately 0.31nund (fee / gas: 20000 / 65000), so the Validator will accept the Tx and include it in the block, since 0.31nund > 0.25nund.

If we had set the `--fees` to 10000, it would not have been processed by the Validator (10000 / 65000 = 0.15nund).

::: tip NOTE
the Tx with lower fees may remain the Tx pool until a validator with lower `minimum-gas-prices` picks it up and proposes the block.
:::

## Example 2: setting --gas and --gas-prices

In this example, we'll set our own `--gas-prices`, and ask `undcli` to estimate the amount of gas the Tx will consume based on the Tx input by passing the `--gas=auto` flag. We can also use the `--gas-adjustment` flag to increase/decrease this gas estimate. We'll assume again that the calculated estimate will be around 65000 gas:

```
undcli tx send from to 123456nund --memo="some und from me to you" --gas=auto --gas-prices=0.25nund
```

In this example, the Tx **fee** will be calculated and included in the transaction for us. The fee will be around 16250nund (`gas * gas-prices`: 65000 * 0.25). Since we have set `gas-prices` already to 0.25 (and assuming the gas estimate is also correct), this Tx will be processed by the validator.

::: tip NOTE
Adding the `--gas-adjustment` flag, for example `--gas-adjustment=1.5`, will increase the gas estimate and therefore the fee, but will increase the chances of the Tx being processed.
:::

Validators will prioritise Txs with higher `gas-prices`, so it is worth setting higher prices to ensure your Tx is processed.
