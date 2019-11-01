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


### Interacting with WRKChain
./build/undcli tx wrkchain register wrkchain1 "genesishashkjwnedjknwed" "Wrkchain 1" --from wrktest --node=http://localhost:26661 --chain-id UND-Mainchain-DevNet
where 54372 is the WRKChainID

undcli query wrkchain get 54372
