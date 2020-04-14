# BEACON Example

[[toc]]

Register:
undcli tx beacon register --moniker=beacon1 --name="Beacon 1" --from wrktest

then run:
```
undcli query tx [TX HASH]
```

this will return the generated Beacon ID integer

Query metadata
```
undcli query beacon beacon 1
```

Record Timestamp hash
```
undcli tx beacon record 1 --hash=d04b98f48e8 --subtime=$(date +%s) --from wrktest --gas=auto --gas-adjustment=1.15
```

Query a Timestamp
```
undcli query beacon timestamp 1 1
```
