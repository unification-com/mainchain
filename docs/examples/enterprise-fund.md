# Enterprise FUND Example

#### Contents

[[toc]]

Raise purchase order:

::: tip
Enterprise Purchase Orders are raised using `nund`
:::

```
undcli tx enterprise purchase 1002000000000nund --from wrktest --gas=auto --gas-adjustment=1.15
```

List purchase orders:
```
undcli query enterprise get-all-pos
```

get specific purchase order:
```
undcli query enterprise get 1
```

Query total locked enterprise FUND
```
undcli query enterprise total-locked
```

Query locked enterprise FUND for an account
```
undcli query enterprise locked [address]
```
