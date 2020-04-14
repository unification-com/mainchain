# Enterprise UND Example

[[toc]]

Raise purchase order:
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

Query total locked enterprise UND
```
undcli query enterprise total-locked
```

Query locked enterprise UND for an account
```
undcli query enterprise locked [address]
```
