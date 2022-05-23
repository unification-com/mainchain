# Enterprise FUND Example

#### Contents

[[toc]]

Raise purchase order:

::: tip
Enterprise Purchase Orders are raised using `nund`
:::

```
und tx enterprise purchase 1002000000000nund --from wrktest --gas=auto --gas-adjustment=1.15
```

List purchase orders:
```
und query enterprise get-all-pos
```

get specific purchase order:
```
und query enterprise get 1
```

Query total locked enterprise FUND
```
und query enterprise total-locked
```

Query locked enterprise FUND for an account
```
und query enterprise locked [address]
```
