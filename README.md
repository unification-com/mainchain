# Unification Mainchain

## Interacting with the Docker Enviroment

### Compositions

Pure upstream:
```
docker-compose -f Docker/docker-compose.upstream.yml up --build
```

Composition with local changes:
```
docker-compose -f Docker/docker-compose.local.yml up --build
```


### Importing docker composition keys
```
undcli keys import node0 Docker/assets/keys/node0.key
undcli keys import node1 Docker/assets/keys/node1.key
undcli keys import node2 Docker/assets/keys/node2.key
```


### Useful Defaults
```
undcli config chain-id UND-Mainchain-DevNet
undcli config node tcp://localhost:26661
```


### Query accounts
```
undcli query account cosmos1cxxsr89u77hu7ksz5nw2cu27pfg88g3v92u7dd
undcli query account cosmos1cvrv3atsm26t4qhssfzj4cs8u7rvsuv9gzwkn6
undcli query account cosmos1ss63vffqmpz68ext374cuxa0v3upavghwzw53p
```


### Interacting with WRKChain module

Register:
undcli tx wrkchain register --moniker="wrkchain1" --genesis="genesishashkjwnedjknwed" --name="Wrkchain 1" --type="geth" --from wrktest --gas=auto --gas-adjustment=1.15

then run:
```
undcli query tx [TX HASH]
```

this will return the generated WRKChain ID integer

Query metadata
```
undcli query wrkchain get 1
```

Filter WRKChain metadata

```
undcli query wrkchain search --moniker wrkchain1
undcli query wrkchain search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
undcli query wrkchain search --page=2 --limit=100
```

Record block hash(es)
```
undcli tx wrkchain record 1 1 d04b98f48e8 f8bcc15c6ae 5ac050801cd6 dcfd428fb5f9e 65c4e16e7807340fa --from wrktest --gas=auto --gas-adjustment=1.15
```

Query a block
```
undcli query wrkchain get-block 1 1
```

Query all blocks
```
undcli query wrkchain blocks 1
```

Filter Blocks

```
undcli query wrkchain search-blocks 1 --before 1574871069 --after 1573481124
undcli query wrkchain search-blocks 1 --min 123 --max 456
undcli query wrkchain search-blocks 1 --page=2 --limit=100
```

#### WRKChain REST

http://localhost:1317/wrkchain/wrkchains  
http://localhost:1317/wrkchain/wrkchains?moniker=wrkchain1  
http://localhost:1317/wrkchain/wrkchains?owner=[bech32address]  
http://localhost:1317/wrkchain/1  
http://localhost:1317/wrkchain/1/blocks  
http://localhost:1317/wrkchain/1/block/1  
http://localhost:1317/wrkchain/1/blocks?min=2&max=4

### Interacting with BEACON module

Register:
undcli tx beacon register beacon1 "Beacon 1" --from wrktest

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
undcli tx beacon record 1 d04b98f48e8 $(date +%s) --from wrktest --gas=auto --gas-adjustment=1.15
```

Query a Timestamp
```
undcli query beacon timestamp 1 1
```

Query all tmestamps
```
undcli query beacon timestamps 1
```


### Purchase Enterprise UND

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

Accept purchase order (must be sent from specified enterprise account)
```
undcli tx enterprise process 1 accept --from ent --gas=auto --gas-adjustment=1.15
```

Reject purchase order:
```
undcli tx enterprise process 1 reject --from ent --gas=auto --gas-adjustment=1.15
```

Query total locked enterprise UND
```
undcli query enterprise total-locked
```

Query locked enterprise UND for an account
```
undcli query enterprise locked [address]
```

## Invariance checking

Start a full node with the `--inv-check-period` flag. Value of 1 will
check every block for invariances:

```
und start --inv-check-period 1
```

Invariance Tx can be sent using:

```
undcli tx crisis invariant-broken enterprise module-account --from wrktest
```

## Simulation

### test full app simulation

```
go test -mod=readonly ./simapp \
    -run=TestFullAppSimulation \
    -Enabled=true \
    -NumBlocks=500 \
    -BlockSize=300 \
    -Commit=true \
    -Seed=24 \
    -Period=1 \
    -PrintAllInvariants=true \
    -ExportParamsPath=/path/to/.simapp/params.json \
    -ExportStatePath=/path/to/.simapp/state.json \
    -ExportStatsPath=/path/to/.simapp/statistics.json \
    -Verbose=true \
    -v \
    -timeout 24h
```

### benchmark test

```
go test -mod=readonly -benchmem -run=^$ github.com/unification-com/mainchain-cosmos/simapp -bench ^BenchmarkFullAppSimulation -Commit=true -cpuprofile /path/to/.simapp/cpu.out -v -timeout 24h
```

#### pprof tools

```
go tool pprof /path/to/.simapp/cpu.out
(pprof) top
(pprof) list [function]
(pprof) web
(pprof) quit
```
