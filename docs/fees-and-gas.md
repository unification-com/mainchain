# Fees and Gas

Transactions consume gas, and the sender must pay a fee in order for the transaction be processed by a validator. The fee is calculated from the amount of gas as Tx will consume multiplied by the gas price. The gas price can be set by the sender of the Tx, but each validator will have set their own `minimum-gas-prices` value, and will not process transactions that do not meet this minimum requirement.

Fees are paid in `nund`, and may be set or calculated depending on which flags are passed to the `undcli` command.

>**Note:** only `--fees` or `--gas-prices` may be used - not both at the same time.

>**Tip**: `--gas-prices` can be used along with the `--gas=auto` and `--gas-adjustment` flags to estimate the gas requirement and automatically calculate the Tx fees.

## Example 1

In this example, we're manually setting the fee for a standard `send` transaction. The validator has a `minimum-gas-prices` of `0.25nund`. We'll set the fee to 20000nund. We'll assume the amount of gas consumed for this `send` transaction, including a small `--memo` is around 65000.

```
undcli tx send from to ... --fees=20000nund
```
