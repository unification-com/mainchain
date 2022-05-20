# BEACON Example

#### Contents

[[toc]]

Register:
und tx beacon register --moniker=beacon1 --name="Beacon 1" --from wrktest

then run:
```
und query tx [TX HASH]
```

this will return the generated Beacon ID integer

Query metadata
```
und query beacon beacon 1
```

Record Timestamp hash
```
und tx beacon record 1 --hash=d04b98f48e8 --subtime=$(date +%s) --from wrktest --gas=auto --gas-adjustment=1.15
```

Query a Timestamp
```
und query beacon timestamp 1 1
```
