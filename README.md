# Unification Mainchain

## Interacting with the Docker Enviroment

### Compositions

Pure upstream:
docker-compose -f Docker/docker-compose.upstream.yml up --build

Composition with local changes:
docker-compose -f Docker/docker-compose.local.yml up --build


### Importing docker composition keys
undcli keys import node0 Docker/assets/keys/node0.key
undcli keys import node1 Docker/assets/keys/node1.key
undcli keys import node2 Docker/assets/keys/node2.key


### Useful Defaults
undcli config chain-id UND-Mainchain-DevNet
undcli config node tcp://localhost:26661


### Query accounts
undcli query account cosmos1cxxsr89u77hu7ksz5nw2cu27pfg88g3v92u7dd
undcli query account cosmos1cvrv3atsm26t4qhssfzj4cs8u7rvsuv9gzwkn6
undcli query account cosmos1ss63vffqmpz68ext374cuxa0v3upavghwzw53p


### Interacting with WRKChain module

Register:
undcli tx wrkchain register wrkchain1 "genesishashkjwnedjknwed" "Wrkchain 1" --from wrktest

then run:
undcli query tx [TX HASH]

this will return the generated WRKChain ID integer

Query metadata
undcli query wrkchain get 1

Record block hash(es)
undcli tx wrkchain record 1 1 d04b98f48e8 f8bcc15c6ae 5ac050801cd6 dcfd428fb5f9e 65c4e16e7807340fa --from wrktest

Query a block
undcli query wrkchain get-block 1 1

Query all blocks
undcli query wrkchain blocks 1

### Purchase Enterprise UND

Raise purchase order:
undcli tx enterprise purchase 1002000000000nund --from wrktest

List purchase orders:
undcli query enterprise get-all-pos

get specific purchase order:
undcli query enterprise get 1

Accept purchase order (must be sent from specified enterprise account)
undcli tx enterprise process 1 accept --from ent

Reject purchase order:
undcli tx enterprise process 1 reject --from ent

Query total locked enterprise UND
undcli query enterprise total-locked

Query locked enterprise UND for an account
undcli query enterprise locked [address]

## Invariance checking

Start a full node with the `--inv-check-period` flag. Value of 1 will
check every block for invariances:

und start --inv-check-period 1

Invariance Tx can be sent using:

undcli tx crisis invariant-broken enterprise module-account --from wrktest

## Simulation

go test -mod=readonly ./simapp -run=TestFullAppSimulation -Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5  -v -timeout 24h
