# Fees and Gas

Transactions consume gas, and the sender must pay a fee in order for the transaction be processed by a validator. The fee is calculated from the amount of gas as Tx will consume multiplied by the gas price. The gas price can be set by the sender of the Tx, but each validator will have set their own `minimum-gas-prices` value, and will not process transactions that do not meet this minimum requirement.

Fees are paid in `nund`, and may be set or calculated depending on which flags are passed to the `undcli` command.

>**Note:** only `--fees` or `--gas-prices` may be used - not both at the same time.

>**Tip**: `--gas-prices` can be used along with the `--gas=auto` and `--gas-adjustment` flags to estimate the gas requirement and automatically calculate the Tx fees.

## Example 1: setting --fees

In this example, we're manually setting the fee for a standard `send` transaction. The validator has a `minimum-gas-prices` of `0.25nund`. We'll set the fee to 20000nund. We'll assume the amount of gas consumed for this `send` transaction, including a small `--memo` is around 65000.

```
undcli tx send from to 123456nund --fees=20000nund
```

In this instance, the `gas-price` is approximately 0.31nund (20000 / 65000), so the Validator will accept the Tx and include it in the block, since 0.31nund > 0.25nund.

If we had set the `--fees` to 10000, it would not have been processed by the Validator (10000 / 65000 = 0.15nund)

## Example 2: setting --gas and --gas-prices

In this example, we'll set our own `--gas-prices`, and have `undcli` estimate the amount of gas the Tx will consume. We can also use the `--gas-adjustment` to increase/decrease the estimate. We'll assume again that the calculated estimate will be around 65000 gas:

```
undcli tx send from to 123456nund --gas=auto --gas-prices=0.25nund
```

In this example, the Tx fee will be calculated and included in the transaction. The fee will be around 16250nund (`gas * gas-prices`: 65000 * 0.25). Since we have set `gas-prices` already to 0.25, this Tx will be processed by the validator.

>**Note**: Adding the `--gas-adjustment` flag, for example `--gas-adjustment=1.5`, will increase the gas estimate and therefore the fee, but will increase the chances of the Tx being processed.

Validators will prioritise Txs with higher `gas-prices`, so it is worth setting higher prices to ensure your Tx is processed.
